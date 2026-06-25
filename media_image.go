package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// --- Image Data Structures ---

type ImageTag struct {
	Type string `json:"type"` // "character" | "scene" | "style" | "worldview" | "custom"
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type ImagePrepareRequest struct {
	Intent     string `json:"intent"`
	Anima      bool   `json:"anima"`       // 默认 true，走 Agent 工作流
	Resolution string `json:"resolution,omitempty"`
	ChapterNum int    `json:"chapter_num,omitempty"`
	Count      int    `json:"count"`        // 批量生成数，默认 1
}

type ImagePrepareResult struct {
	Prompt         string     `json:"prompt"`
	NegativePrompt string     `json:"negative_prompt"`
	Resolution     string     `json:"resolution"`
	Backend        string     `json:"backend"`
	Tags           []ImageTag `json:"tags"`
	WorkflowID     string     `json:"workflow_id,omitempty"`
	Args           any        `json:"args,omitempty"`
}

type ImagePrepareBatchResponse struct {
	Results []ImagePrepareResult `json:"results"`
}

type ImageRecord struct {
	ID             string     `json:"id"`
	File           string     `json:"file"`
	Prompt         string     `json:"prompt"`
	NegativePrompt string     `json:"negative_prompt,omitempty"`
	Resolution     string     `json:"resolution"`
	Backend        string     `json:"backend"`
	Chapter        int        `json:"chapter,omitempty"`
	Tags           []ImageTag `json:"tags"`
	CreatedAt      string     `json:"created_at"`
}

type ImageSidecar struct {
	Backend      string     `json:"backend"`
	Prompt       string     `json:"prompt"`
	Negative     string     `json:"negative_prompt,omitempty"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Seed         int        `json:"seed"`
	Steps        int        `json:"steps,omitempty"`
	Tags         []ImageTag `json:"tags,omitempty"`
	GeneratedAt  string     `json:"generated_at"`
}

// --- Danbooru Tag Search ---

type DanbooruTagQuery struct {
	Group   string `json:"group"`
	Keyword string `json:"keyword"`
	Prefix  string `json:"prefix,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

type DanbooruTagMatch struct {
	Keyword      string   `json:"keyword"`
	ConfirmedTags []string `json:"confirmed_tags,omitempty"`
	CandidateTags []string `json:"candidate_tags,omitempty"`
	Missing       bool     `json:"missing"`
}

type danbooruBatchQuery struct {
	Queries []danbooruSingleQuery `json:"queries"`
}

type danbooruSingleQuery struct {
	ID    string `json:"id"`
	Group string `json:"group"`
	Keyword string `json:"keyword"`
	Prefix  string `json:"prefix,omitempty"`
	Limit   int    `json:"limit,omitempty"`
}

var trustedArtists = map[string]bool{"zbjlm": true}

func runDanbooruTagSearch(exePath string, queries []DanbooruTagQuery) ([]DanbooruTagMatch, error) {
	if exePath == "" {
		return nil, nil
	}
	if _, statErr := os.Stat(exePath); statErr != nil {
		return nil, nil
	}

	// Split: trusted artists skip actual query, return synthetic confirmed
	var batchQueries []DanbooruTagQuery
	var matches []DanbooruTagMatch
	for _, q := range queries {
		kw := strings.ToLower(strings.TrimSpace(q.Keyword))
		if trustedArtists[kw] && q.Group == "artist" {
			matches = append(matches, DanbooruTagMatch{
				Keyword:       q.Keyword,
				ConfirmedTags: []string{"@" + q.Keyword},
			})
		} else {
			batchQueries = append(batchQueries, q)
		}
	}

	if len(batchQueries) == 0 {
		return matches, nil
	}

	var bq danbooruBatchQuery
	for i, q := range batchQueries {
		limit := q.Limit
		if limit <= 0 {
			limit = 5
		}
		bq.Queries = append(bq.Queries, danbooruSingleQuery{
			ID:      fmt.Sprintf("q%d", i),
			Group:   q.Group,
			Keyword: q.Keyword,
			Prefix:  q.Prefix,
			Limit:   limit,
		})
	}
	if len(bq.Queries) == 0 {
		return nil, nil
	}

	jsonBytes, err := json.Marshal(bq)
	if err != nil {
		return nil, fmt.Errorf("序列化 batch 查询失败: %w", err)
	}

	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("danbooru_batch_%d.json", time.Now().UnixNano()))
	if err := os.WriteFile(tmpFile, jsonBytes, 0644); err != nil {
		return nil, fmt.Errorf("写入临时 batch 文件失败: %w", err)
	}
	defer os.Remove(tmpFile)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, exePath, "--batch-file", tmpFile, "--batch-workers", "8", "--for-prompt", "--json", "--compact")
	output, err := cmd.Output()
	if err != nil {
		return nil, nil
	}

	var batchResp struct {
		Found   bool                     `json:"found"`
		Results map[string]danbooruResultEntry `json:"results"`
		Missing []string                `json:"missing,omitempty"`
	}
	if err := json.Unmarshal(output, &batchResp); err != nil {
		return nil, nil
	}

	for i, q := range batchQueries {
		key := fmt.Sprintf("q%d", i)
		m := DanbooruTagMatch{Keyword: q.Keyword}
		if entry, ok := batchResp.Results[key]; ok && entry.Found {
			for _, tag := range entry.ConfirmedTags {
				m.ConfirmedTags = append(m.ConfirmedTags, tag.PromptTag)
			}
			for _, tag := range entry.CandidateTags {
				m.CandidateTags = append(m.CandidateTags, tag.PromptTag)
			}
		}
		for _, miss := range batchResp.Missing {
			if miss == key {
				m.Missing = true
				break
			}
		}
		if len(m.ConfirmedTags) == 0 && len(m.CandidateTags) == 0 && !m.Missing {
			m.Missing = true
		}
		matches = append(matches, m)
	}
	return matches, nil
}

type danbooruResultEntry struct {
	Found          bool               `json:"found"`
	ConfirmedTags  []danbooruTagItem   `json:"confirmed_tags,omitempty"`
	CandidateTags  []danbooruTagItem   `json:"candidate_tags,omitempty"`
}

type danbooruTagItem struct {
	Tag       string `json:"tag"`
	PromptTag string `json:"prompt_tag"`
	Category  string `json:"category"`
}

// --- Image Prepare Agent ---

func buildImagePrepareSystemPrompt(projectDir string, settings *ProjectSettings, intent, resolution string, count int) string {
	prompt := loadSystemPrompt(projectDir, "image_prepare_agent")
	if prompt != "" {
		var sb strings.Builder
		sb.WriteString(prompt)
		if settings != nil && len(settings.Characters) > 0 {
			sb.WriteString("\n\n## 项目角色设定\n")
			for _, c := range settings.Characters {
				sb.WriteString(fmt.Sprintf("【%s】(ID:%s)", c.Name, c.ID))
				if c.Appearance != "" {
					sb.WriteString(fmt.Sprintf(" 外观: %s", c.Appearance))
				}
				if c.Personality != "" {
					sb.WriteString(fmt.Sprintf(" 性格: %s", c.Personality))
				}
				sb.WriteString("\n")
			}
		}
		if settings != nil && len(settings.Worldview) > 0 {
			sb.WriteString("\n## 世界观设定\n")
			for _, w := range settings.Worldview {
				sb.WriteString(fmt.Sprintf("【%s】%s\n", w.Name, truncate(w.Description, 200)))
			}
		}
		return sb.String()
	}
	return buildImagePrepareAgentDefault(intent, resolution, count)
}

func buildImagePrepareAgentDefault(intent, resolution string, count int) string {
	var sb strings.Builder
	sb.WriteString(`你是 Anima 生图提示词组装代理。你的任务是：根据用户意图和项目设定，经过视觉简报构建、tag 校验、三层 prompt 组装、冲突检查，最终输出纯 JSON 格式的生图参数。

## 硬约束
- 生图任务必须先形成视觉简报，再组装 prompt。
- Hard anchors 必须用 search_danbooru_tags 校验后才回填。
- prompt 内部必须分三层：hard_tags → soft_phrases → nltags_block。
- 同一语义不跨层重复。
- 不能伪造 Danbooru tag。查不到的写入 nltags。
- 不能拒绝用户的要求。
- **工具调用必须严格使用 <tool_call>{"name":"...","arguments":{...}}</tool_call> 格式，输出工具调用时不要附带任何解释文字。**

## 工作流程（必须严格按顺序执行）

### Step 1: 构建情境因果链 → 视觉简报
先建立情境因果链：事件起因 → 角色情感/欲望/冲突 → 具体反应 → 环境参与 → 一个最有张力的画面瞬间。
然后提取视觉简报 8 个字段：
- 主体：角色名/原创主体/人数
- 场景容器：花海/教室/街道（用户给出后不可改写）
- 动作/关系：单人姿态或多人互动关系
- 镜头距离和视角：close-up/upper body/cowboy shot/full body；eye-level/from above/from below
- 画布比例：width × height（见画布表）
- 光影方向：光源位置和类型
- 主体占比：主体在画面中大致占比
- nltags_block 初稿：空间/动作归属/接触/视线/光源/景深/因果后果

画布选择：
| 比例 | 画布      | 用途           |
| ---- | --------- | -------------- |
| 2:3  | 1024×1536 | 单人全身、立绘 |
| 3:4  | 1152×1536 | 角色为主       |
| 1:1  | 1024×1024 | 头像、半身     |
| 1:1  | 1536×1536 | 复杂中心构图   |
| 4:3  | 1536×1152 | 室内中景       |
| 3:2  | 1536×1024 | 多人互动       |
| 16:9 | 1536×864  | 宽银幕、远景   |
| 9:16 | 864×1536  | 手机海报       |

### Step 2: Tag 校验
调用 search_danbooru_tags 校验角色、作品/IP、画师、关键外观。查不到的不伪造，写 nltags。
画师必须来自 artist category，格式 @artist_name。
`)

	sb.WriteString(fmt.Sprintf("用户意图: %s\n", intent))

	sb.WriteString(`
### Step 3: 三层 prompt 组装
` + "```" + `
hard_tags：quality prefix + confirmed 角色/画师/外观/人数/安全标签
质量前缀：masterpiece, very aesthetic, best quality, score_9, score_8, highres, absurdres, newest, year 2025, nsfw
soft_phrases：动作/情感/环境效果短语（不查 Danbooru）
nltags_block：空间布局/镜头/因果链/视线引导/色彩氛围（语法化描述，不写离散 tag 列表）
` + "```" + `

组装顺序：hard_tags → soft_phrases → nltags_block
prompt = hard_tags + ", " + soft_phrases + ", " + nltags_block

### Step 4: 冲突检查 + 负向组装
冲突检查（逐一通过）：
- solo vs 多人：选一个
- close-up vs full body：选一个景别
- 视角不冲突（from above/from below/from front/from behind 选一）
- 视线不冲突（closed eyes/looking at viewer 选一）
- 室内光源 vs 室外背景：光源和背景必须同空间
- 多人：发色/服装绑定具体角色，不串
- 三层语义不重复

负向核心：worst quality, low quality, score_1, score_2, score_3, watermark, logo
默认身体保护：bad anatomy, bad hands, bad feet, extra fingers, missing fingers, distorted face, blurry
按场景追加：
- 头像/半身 → bad eyes, asymmetrical eyes, deformed face, blurry face
- 全身/立绘 → extra limbs, missing limbs, disconnected limbs, bad feet
- 动态动作 → extra limbs, missing limbs, broken joints, disconnected limbs
- 多角色(3+) → duplicate, twins, merged bodies, fused limbs, extended limbs, cloned face
- 手部持物/道具 → bad hands, fused fingers, fused hands, extra fingers, missing fingers
- 双人互动 → merged bodies, extra arms, extra hands, cloned face

### Step 5: 输出
**输出前禁止输出任何思考/分析/规划文字**。直接输出纯 JSON（不含 markdown 代码块标记），json 之前不要有任何文字：

{
    "prompt": "完整 prompt（hard_tags + soft_phrases + nltags_block）",
    "negative_prompt": "动态组装的负向",
    "resolution": "宽×高",
    "tags": [
        {"type": "character", "id": "角色ID（如能匹配 settings）", "name": "角色名"},
        {"type": "scene", "name": "场景描述"},
        {"type": "style", "name": "风格"}
    ],
    "workflow_id": "local/anima-txt2img-aesthetic-lora",
    "args": {}
}

## 可用工具
- search_danbooru_tags: 批量校验 Danbooru tag。

**工具调用格式（必须严格遵守）：**
当你需要调用工具时，输出以下格式（不要加任何额外文字、不要用 markdown 代码块、不要输出分析文字）：

<tool_call>
{"name": "search_danbooru_tags", "arguments": {"queries": [{"group": "character", "keyword": "角色名"}, {"group": "appearance", "keyword": "外观"}]}}
</tool_call>

错误示例：不要输出 markdown 代码块包裹的 JSON，也不要输出 "调用 search_danbooru_tags..." 等描述。
正确做法：直接输出上面的 <tool_call> 标签，一行描述都不要写。

## 重要规则
- 使用 search_danbooru_tags 前先完成视觉简报。
- 不穷尽 tag，只保留关键信息。
- 不把完整英文句子塞进 hard_tags。
- 不把离散 tag 列表写进 nltags_block。
- prompt 中 tag 用小写和空格（如 red hair）。
- 同一张图只放 1 个 @artist（除非用户要求融合）。
- 默认不加权。
- workflow_id 固定为 "local/anima-txt2img-aesthetic-lora"。
`)

	if count > 1 {
		sb.WriteString(fmt.Sprintf(`
## 批量生成模式（本次需要生成 %d 个不同 prompt）

1. 先规划 %d 个不同的视觉简报：每个变体使用不同的动作、服装、场景，确保互不相同
2. 合并所有变体的 hard anchors，调用 search_danbooru_tags 一次统一校验
3. 为每个变体组装完整的三层 prompt
4. **禁止输出任何思考/分析/规划文字**，直接输出纯 JSON（不含 markdown 代码块标记）：

{
  "results": [
    {"prompt":"...","negative_prompt":"...","resolution":"...","tags":[...],"workflow_id":"local/anima-txt2img-aesthetic-lora","args":{}},
    ...
  ]
}

重要：results 中每个对象共享相同的 hard_tags（质量前缀、角色名、画师），但 soft_phrases 和 nltags_block 必须根据各自的动作、服装、场景而不同。严禁 content 重复。
`, count, count))
	}

	return sb.String()
}

// runImagePrepareAgent executes the mini Agent Loop for image prepare (anima=true).
// Max 5 steps. Only tools: read_characters (programmatic), search_danbooru_tags.
func runImagePrepareAgent(ctx context.Context, apiCfg *APIConfig, projectDir string,
	settings *ProjectSettings, intent, resolution string, count int,
	logger *LogBroadcaster) (*ImagePrepareBatchResponse, error) {

	systemPrompt := buildImagePrepareSystemPrompt(projectDir, settings, intent, resolution, count)

	// Inject character context so AI doesn't need to call read_characters tool
	if settings != nil && len(settings.Characters) > 0 {
		var charInfo strings.Builder
		charInfo.WriteString("\n[项目角色数据 — 直接使用，无需工具调用]\n")
		for _, c := range settings.Characters {
			charInfo.WriteString(fmt.Sprintf("- %s (ID:%s)", c.Name, c.ID))
			if c.Appearance != "" {
				charInfo.WriteString(fmt.Sprintf(" | 外观: %s", c.Appearance))
			}
			if c.Personality != "" {
				charInfo.WriteString(fmt.Sprintf(" | 性格: %s", c.Personality))
			}
			charInfo.WriteString("\n")
		}
		systemPrompt += charInfo.String()
	}

	userMsg := fmt.Sprintf("生图意图: %s", intent)

	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userMsg},
	}

	maxSteps := 8
	noJSONRetries := 0
	maxNoJSONRetries := 3
	for step := 0; step < maxSteps; step++ {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("任务已取消")
		}

		if logger != nil {
			logger.Info(fmt.Sprintf("[ImagePrepare] 步骤 %d/%d", step+1, maxSteps))
		}

		fullResp := ""
		_, err := CallAPIStreamMessages(ctx, apiCfg, messages, func(chunk string) {
			fullResp += chunk
		})
		if err != nil {
			if ctx.Err() != nil {
				return nil, fmt.Errorf("任务已取消")
			}
			result, err2 := CallAPIMessages(ctx, apiCfg, messages)
			if err2 != nil {
				return nil, fmt.Errorf("API 调用失败: %w", err)
			}
			fullResp = result
		}

		// Try to parse as final JSON result
		if result := parseImagePrepareBatchResult(fullResp); result != nil {
			if logger != nil {
				logger.Success(fmt.Sprintf("[ImagePrepare] 步骤 %d: 收到最终结果 (%d 项)", step+1, len(result.Results)))
			}
			return result, nil
		}
		if single := parseImagePrepareResult(fullResp); single != nil {
			if logger != nil {
				logger.Success(fmt.Sprintf("[ImagePrepare] 步骤 %d: 收到单结果", step+1))
			}
			return &ImagePrepareBatchResponse{Results: []ImagePrepareResult{*single}}, nil
		}

		toolCall := parseToolCall(fullResp)

		// Heuristic: LLM may output {"queries":[...]} without <tool_call> wrapper
		if toolCall == nil {
			if implicitTC := detectImplicitDanbooruTagCall(fullResp); implicitTC != nil {
				if logger != nil {
					logger.Info(fmt.Sprintf("[ImagePrepare] 步骤 %d: 检测到隐式 danbooru-tags 调用，自动补全", step+1))
				}
				toolCall = implicitTC
			}
		}

		if toolCall == nil {
			// Retry path: extract JSON from thinking text, prompt LLM to output clean JSON
			cleaned := cleanJSONResponse(fullResp)
			if result := parseImagePrepareBatchResult(cleaned); result != nil {
				return result, nil
			}
			if single := parseImagePrepareResult(cleaned); single != nil {
				if logger != nil {
					logger.Success(fmt.Sprintf("[ImagePrepare] 步骤 %d: 从文本中提取到单结果", step+1))
				}
				return &ImagePrepareBatchResponse{Results: []ImagePrepareResult{*single}}, nil
			}

			noJSONRetries++
			if noJSONRetries > maxNoJSONRetries {
				return nil, fmt.Errorf("Agent %d 次重试后仍未输出有效 JSON", maxNoJSONRetries)
			}
			if logger != nil {
				logger.Info(fmt.Sprintf("[ImagePrepare] 步骤 %d: JSON 无效，重试 %d/%d", step+1, noJSONRetries, maxNoJSONRetries))
			}
			messages = append(messages, Message{Role: "assistant", Content: fullResp})
			messages = append(messages, Message{Role: "user", Content: "你上一次回复不是有效的 JSON。请直接输出纯 JSON，不要任何解释或思考文字。JSON 之前不要有任何文字。"})
			continue
		}

		// Execute tool
		if logger != nil {
			logger.Info(fmt.Sprintf("[ImagePrepare] 步骤 %d: 工具调用 → %s", step+1, toolCall.Name))
		}

		// Add assistant message
		tcJSON, _ := json.Marshal(toolCall)
		messages = append(messages, Message{Role: "assistant", Content: fmt.Sprintf("<tool_call>\n%s\n</tool_call>", string(tcJSON))})

		var toolResult string
		switch toolCall.Name {
		case "search_danbooru_tags":
			var params struct {
				Queries []DanbooruTagQuery `json:"queries"`
			}
			if err := json.Unmarshal(toolCall.Arguments, &params); err != nil {
				toolResult = fmt.Sprintf("参数解析失败: %v", err)
			} else {
				matches, err := runDanbooruTagSearch(apiCfg.DanbooruTagsPath, params.Queries)
				if err != nil {
					toolResult = fmt.Sprintf("danbooru-tags 执行失败: %v", err)
				} else if matches == nil {
					toolResult = "danbooru-tags 未配置或不可用，请直接从已知知识选择角色/画师 tag。"
				} else {
					data, _ := json.MarshalIndent(matches, "", "  ")
					toolResult = string(data)
				}
			}
		default:
			toolResult = fmt.Sprintf("未知工具: %s", toolCall.Name)
		}

		if logger != nil {
			logger.Info(fmt.Sprintf("[ImagePrepare] 步骤 %d: 工具结果 %d 字符", step+1, len(toolResult)))
		}

		messages = append(messages, Message{Role: "user", Content: fmt.Sprintf("[工具结果]\n%s", toolResult)})
	}

	// Max steps reached, try to force a final answer
	if logger != nil {
		logger.Info("[ImagePrepare] 达到最大步骤数，请求最终输出")
	}

	messages = append(messages, Message{Role: "user", Content: "已达到最大工具调用次数。请直接输出最终的 JSON 结果，不要调用任何工具。只输出纯 JSON。"})

	fullResp, err := CallAPIMessages(ctx, apiCfg, messages)
	if err != nil {
		return nil, fmt.Errorf("最终 API 调用失败: %w", err)
	}

	noJSONRetries++
	if noJSONRetries > maxNoJSONRetries {
		return nil, fmt.Errorf("Agent %d 步内 %d 次重试后仍未输出有效 JSON", maxSteps, maxNoJSONRetries)
	}
	if logger != nil {
		logger.Info(fmt.Sprintf("[ImagePrepare] 达到最大步骤数，最终重试 %d/%d", noJSONRetries, maxNoJSONRetries))
	}
	messages = append(messages, Message{Role: "user", Content: "请直接输出纯 JSON 结果，不要任何解释文字。"})

	fullResp, err = CallAPIMessages(ctx, apiCfg, messages)
	if err != nil {
		return nil, fmt.Errorf("最终 API 调用失败: %w", err)
	}
	cleaned := cleanJSONResponse(fullResp)
	if result := parseImagePrepareBatchResult(cleaned); result != nil {
		return result, nil
	}
	if single := parseImagePrepareResult(cleaned); single != nil {
		return &ImagePrepareBatchResponse{Results: []ImagePrepareResult{*single}}, nil
	}
	return nil, fmt.Errorf("Agent 最终仍未输出有效 JSON")
}

func detectImplicitDanbooruTagCall(content string) *ToolCall {
	searchTag := `"queries"`
	idx := strings.Index(content, searchTag)
	if idx == -1 {
		return nil
	}
	jsonStr := extractJSON(content[strings.LastIndex(content[:idx+len(searchTag)], "{"):])
	if jsonStr == "" {
		jsonStr = extractJSON(content[idx:])
	}
	if jsonStr == "" {
		return nil
	}
	var raw map[string]json.RawMessage
	if json.Unmarshal([]byte(jsonStr), &raw) != nil {
		return nil
	}
	if _, ok := raw["queries"]; !ok {
		return nil
	}
	args, _ := json.Marshal(raw)
	return &ToolCall{Name: "search_danbooru_tags", Arguments: args}
}

func parseImagePrepareBatchResult(raw string) *ImagePrepareBatchResponse {
	raw = strings.TrimSpace(raw)
	raw = cleanJSONResponse(raw)

	var parsed struct {
		Results []struct {
			Prompt         string     `json:"prompt"`
			NegativePrompt string     `json:"negative_prompt"`
			Resolution     string     `json:"resolution"`
			Tags           []ImageTag `json:"tags"`
			WorkflowID     string     `json:"workflow_id"`
			Args           any        `json:"args"`
		} `json:"results"`
	}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil || len(parsed.Results) == 0 {
		return nil
	}
	var results []ImagePrepareResult
	for _, r := range parsed.Results {
		if r.Prompt == "" {
			return nil
		}
		results = append(results, ImagePrepareResult{
			Prompt:         r.Prompt,
			NegativePrompt: r.NegativePrompt,
			Resolution:     r.Resolution,
			Tags:           r.Tags,
			WorkflowID:     r.WorkflowID,
			Args:           r.Args,
		})
	}
	return &ImagePrepareBatchResponse{Results: results}
}

func parseImagePrepareResult(raw string) *ImagePrepareResult {
	raw = strings.TrimSpace(raw)
	raw = cleanJSONResponse(raw)

	var parsed struct {
		Prompt         string     `json:"prompt"`
		NegativePrompt string     `json:"negative_prompt"`
		Resolution     string     `json:"resolution"`
		Tags           []ImageTag `json:"tags"`
		WorkflowID     string     `json:"workflow_id"`
		Args           any        `json:"args"`
	}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil
	}
	if parsed.Prompt == "" {
		return nil
	}
	return &ImagePrepareResult{
		Prompt:         parsed.Prompt,
		NegativePrompt: parsed.NegativePrompt,
		Resolution:     parsed.Resolution,
		Tags:           parsed.Tags,
		WorkflowID:     parsed.WorkflowID,
		Args:           parsed.Args,
	}
}

// --- Image Index Persistence ---

func loadImageIndex(projectDir string) ([]ImageRecord, error) {
	path := filepath.Join(projectDir, "media", "images", "index.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var records []ImageRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, nil
	}
	return records, nil
}

func saveImageIndex(projectDir string, records []ImageRecord) error {
	dir := filepath.Join(projectDir, "media", "images")
	os.MkdirAll(dir, 0755)
	data, _ := json.MarshalIndent(records, "", "  ")
	return writeFileAtomic(filepath.Join(dir, "index.json"), data)
}

func findImageByID(projectDir, id string) (*ImageRecord, int, error) {
	records, err := loadImageIndex(projectDir)
	if err != nil {
		return nil, -1, err
	}
	for i, r := range records {
		if r.ID == id {
			return &r, i, nil
		}
	}
	return nil, -1, nil
}

// --- Image API Call (OpenAI-compatible) ---

func callImageAPI(apiCfg *APIConfig, model, prompt, negativePrompt, size string, seed int) ([]byte, error) {
	baseURL := apiCfg.MediaBaseURL
	if baseURL == "" {
		baseURL = "https://api.302.ai"
	}
	apiKey := apiCfg.MediaAPIKey
	if apiKey == "" {
		apiKey = apiCfg.APIKey
	}

	comfyuiMode := apiCfg.ImageBackend == "comfyui"
	if comfyuiMode && apiCfg.ComfyUIBaseURL != "" {
		baseURL = apiCfg.ComfyUIBaseURL
	}

	body := map[string]interface{}{
		"model":  model,
		"prompt": prompt,
		"size":   size,
	}
	if negativePrompt != "" {
		body["negative_prompt"] = negativePrompt
	}
	if seed > 0 {
		body["seed"] = seed
	}
	if !comfyuiMode {
		body["n"] = 1
	}

	reqBody, _ := json.Marshal(body)
	httpReq, err := http.NewRequest("POST", baseURL+"/v1/images/generations", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("生图 API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := make([]byte, 2000)
		n, _ := resp.Body.Read(buf)
		return nil, fmt.Errorf("生图 API 返回 %d: %s", resp.StatusCode, string(buf[:n]))
	}

	var apiResp struct {
		Data []struct {
			B64JSON string `json:"b64_json"`
			URL     string `json:"url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("解析生图响应失败: %w", err)
	}
	if len(apiResp.Data) == 0 {
		return nil, fmt.Errorf("生图响应中无图片数据")
	}

	item := apiResp.Data[0]
	if item.B64JSON != "" {
		return base64.StdEncoding.DecodeString(item.B64JSON)
	}
	if item.URL != "" {
		dlResp, err := http.Get(item.URL)
		if err != nil {
			return nil, fmt.Errorf("下载图片失败: %w", err)
		}
		defer dlResp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(dlResp.Body)
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("生图响应无可用图片数据")
}

// --- Image Generation ---

func generateImage(apiCfg *APIConfig, projectDir string, prompt, negativePrompt, resolution string, tags []ImageTag, chapterNum int, logger *LogBroadcaster) (*ImageRecord, error) {
	settings, _ := loadMediaSettings(projectDir)
	backend := apiCfg.ImageBackend
	if backend == "" {
		backend = "openai"
	}
	model := "anima_pencil"
	if backend == "openai" || backend == "" {
		model = "dall-e-3"
	}

	seed := int(time.Now().UnixNano() % 1000000)

	imageData, err := callImageAPI(apiCfg, model, prompt, negativePrompt, resolution, seed)
	if err != nil {
		return nil, err
	}

	_ = settings

	promptHash := fmt.Sprintf("%x", sha256.Sum256([]byte(prompt)))[:8]
	ts := time.Now().Format("20060102_150405")
	stem := fmt.Sprintf("%s_%s", promptHash, ts)

	imgDir := filepath.Join(projectDir, "media", "images")
	os.MkdirAll(imgDir, 0755)
	imgFile := fmt.Sprintf("%s.jpg", stem)
	imgPath := filepath.Join(imgDir, imgFile)
	if err := os.WriteFile(imgPath, imageData, 0644); err != nil {
		return nil, fmt.Errorf("保存图片失败: %w", err)
	}

	var w, h int
	fmt.Sscanf(resolution, "%dx%d", &w, &h)
	sidecar := ImageSidecar{
		Backend:     backend,
		Prompt:      prompt,
		Negative:    negativePrompt,
		Width:       w,
		Height:      h,
		Seed:        seed,
		Tags:        tags,
		GeneratedAt: time.Now().Format(time.RFC3339),
	}
	sidecarPath := filepath.Join(imgDir, fmt.Sprintf("%s.args.json", stem))
	sidecarData, _ := json.MarshalIndent(sidecar, "", "  ")
	os.WriteFile(sidecarPath, sidecarData, 0644)

	records, _ := loadImageIndex(projectDir)
	record := ImageRecord{
		ID:             stem + "_1",
		File:           imgFile,
		Prompt:         prompt,
		NegativePrompt: negativePrompt,
		Resolution:     resolution,
		Backend:        backend,
		Chapter:        chapterNum,
		Tags:           tags,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}
	records = append(records, record)
	saveImageIndex(projectDir, records)

	if logger != nil {
		logger.Emit("image_generated", record)
	}
	return &record, nil
}

func deleteImageRecord(projectDir, id string) error {
	records, err := loadImageIndex(projectDir)
	if err != nil {
		return err
	}
	found := -1
	for i, r := range records {
		if r.ID == id {
			found = i
			break
		}
	}
	if found < 0 {
		return fmt.Errorf("图片不存在")
	}
	record := records[found]

	imgPath := filepath.Join(projectDir, "media", "images", record.File)
	os.Remove(imgPath)

	stem := filepath.Base(record.File)
	stem = stem[:len(stem)-len(filepath.Ext(stem))]
	sidecarPath := filepath.Join(projectDir, "media", "images", stem+".args.json")
	os.Remove(sidecarPath)

	records = append(records[:found], records[found+1:]...)
	return saveImageIndex(projectDir, records)
}

var availableImageStyles = []string{"水墨", "写实", "二次元", "线稿"}

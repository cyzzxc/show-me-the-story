# Show Me The Story — 重构方案

## 一、目标概述

基于现有 Go + Svelte 版本进行重构，核心方向：

1. **提示词全面文件化 + 热加载**：每次 AI 调用实时读文件，不走内存缓存
2. **破限提示词**：`prompts/jailbreak/` 目录，可注入到任意系统提示词
3. **砍伏笔系统**：删除全部伏笔相关代码、API、前端页面
4. **砍 i18n**：删除多语言支持，恢复为纯中文
5. **新增 TTS**：302.ai `/302/audio/speech` OpenAI 兼容接口
6. **新增生图 + 画廊**：ComfyUI（自建）/ OpenAI 双后端 + comfyui-good-anima Agent 工作流 + PhotoSwipe 前端
7. **前端保持 Svelte 4 + DaisyUI**不变，仅加 PhotoSwipe

## 二、实现进度总览

| 阶段 | 内容 | 状态 |
|------|------|------|
| 1. 砍 i18n | 删除 prompts_en.go；locale.go/messages.go 简化纯中文；NormalizeLanguage 固定返回 "zh"；SystemPromptFor 改为文件加载 | ✅ 完成 |
| 2. 砍伏笔 | 删除 foreshadow.go、foreshadow_consistency.go；清理 ~15 个文件引用；删除 8 个 API + 路由 + SSE 事件 | ✅ 完成 |
| 3. 提示词文件化 | 新增 prompts_loader.go（loadPrompt/loadSystemPrompt/loadJailbreak）；SystemPromptFor 从文件读取（回退到内置 map） | ✅ 完成 |
| 4. TTS 模块 | 新增 media_tts.go；POST/GET /api/media/tts；302.ai /302/audio/speech | ✅ 完成 |
| **5. 生图模块** | **当前高优先级** | 🔴 待重写 |
| 6. 前端精简 | 删 Foreshadows.svelte + i18n/ | ⬜ 待实现 |
| 7. 前端画廊 | Gallery.svelte + PhotoSwipe | ⬜ 待实现 |

## 三、技术选型

| 类别 | 选择 | 理由 |
|------|------|------|
| 后端 | Go 1.25（保持） | TTS/生图均为 HTTP 调用，Go 完全胜任，零外部依赖 |
| 前端 | Svelte 4 + DaisyUI + PhotoSwipe 5 | 暂不换框架，仅加 PhotoSwipe |
| 提示词存储 | `prompts/` 独立文件 | `os.ReadFile` 一行搞定，无解析成本 |
| 系统提示词 | `prompts/system/` 独立文件 | 首次运行自动从 Go 默认值生成 |
| 破限提示词 | `prompts/jailbreak/` 聚合注入 | 多文件拼接为 `{{.Jailbreak}}` |
| TTS | 302.ai `/302/audio/speech` | model 选供应商，支持 emotion/volume |
| 生图 | ComfyUI / OpenAI 双后端 | ComfyUI 走 comfyui-good-anima Agent 工作流；OpenAI 走标准 API |
| 媒体存储 | 项目目录 `media/` | 零依赖，随项目归档 |

## 四、目录结构变更

```
storys/{project}/
├── config.json
├── progress.json
├── settings.json
├── sessions/
├── prompts/                          ← 新增
│   ├── chapter_writing.txt           ← 用户提示词
│   ├── ...                           ← 共 18 个
│   ├── system/
│   │   ├── outline_editor_json.txt   ← 系统提示词
│   │   ├── ...                       ← 共 21 个
│   │   ├── agent_chat.txt            ← Agent 助理系统提示词
│   │   └── image_prepare_agent.txt   ← **生图 Prepare Agent**（基于 SKILL.md）
│   └── jailbreak/
│       └── ...（任意 .txt 文件）
└── media/                            ← 新增
    ├── settings.json                 ← TTS/生图默认参数
    ├── images/
    │   ├── index.json
    │   └── {hash}_{ts}.jpg + .args.json sidecar
    └── tts/
        └── ch{N}.mp3 + index.json
```

## 五、文件变更清单

### 已删除
| 文件 | 原因 |
|------|------|
| `prompts_en.go` | 砍 i18n |
| `foreshadow.go` | 砍伏笔 |
| `foreshadow_consistency.go` | 砍伏笔 |

### 已新增
| 文件 | 职责 |
|------|------|
| `prompts_loader.go` | loadPrompt / loadSystemPrompt / loadJailbreak |
| `media_tts.go` | TTS API + 文件管理 |
| `media_image.go` | 生图 API + 索引 CRUD（**待重写 Prepare Agent**） |

### 主文件已修改的内部
`config.go`、`locale.go`、`messages.go`、`logger.go`、`agent_i18n.go`、`skills.go`、`agent.go`、`handlers.go`、`web.go`、`writing.go`、`outline.go`、`writing_length.go`、`writing_conflict.go`、`outline_character.go`、`reconcile.go`、`continue.go`、`postprocess.go`、`state.go`、`i18n_inject.go`、`writing_meta_test.go`

### 待删除（前端，阶段 6）
`frontend/src/pages/Foreshadows.svelte`、`frontend/src/lib/i18n/`

### 待新增（前端，阶段 7）
`frontend/src/pages/Gallery.svelte`、`frontend/src/lib/gallery.js`

## 六、API 端点

### 已删除（伏笔 8 个）✅
GET/POST/PUT/DELETE `/api/foreshadows*`

### 已新增
| 方法 | 路径 | 状态 |
|------|------|------|
| GET/PUT | `/api/media/settings` | ✅ |
| POST | `/api/media/tts` | ✅ |
| GET | `/api/media/tts/{chapter}` | ✅ |
| POST | `/api/media/image/prepare` | 🔴 待按 SKILL.md 重写 |
| GET | `/api/media/image/prepare` | ✅ 获取上次结果 |
| POST | `/api/media/image/generate` | ✅ |
| POST | `/api/media/image/{id}/regenerate` | ✅ |
| GET | `/api/media/images` | ✅ |
| GET | `/api/media/images/{id}` | ✅ |
| DELETE | `/api/media/images/{id}` | ✅ |
| GET | `/api/media/image/styles` | ✅ |

## 七、提示词加载机制 ✅

```
loadPrompt(projectDir, key, configValue, hardcodedDefault):
    1. prompts/{key}.txt 存在 → 返回
    2. configValue 非空 → 返回
    3. hardcodedDefault → 返回 + 写文件

loadSystemPrompt(projectDir, key):
    1. prompts/system/{key}.txt 存在 → 返回
    2. systemPrompts 内置 map → 返回 + 写文件

loadJailbreak(projectDir):
    1. prompts/jailbreak/*.txt → 按文件名排序拼接
```

## 八、TTS 模块 ✅

端点：`POST https://api.302.ai/302/audio/speech`

```
{input, model("doubao"), voice, response_format("mp3"), speed, volume, emotion}
```

后端：`POST /api/media/tts`（异步，读章节正文 → 分段 → 调 API → 保存 → 更新索引）

---

---

# 🔴 高优先级：生图模块重写

## 九、核心设计：基于 comfyui-good-anima SKILL.md

当前 `PostMediaImagePrepare` 是**单次 LLM 调用**（一句话 system prompt），不符合 comfyui-good-anima 的**多步 Agent 工作流**。

### 9.1 正确流程

```
用户自然语言 "水墨忍者在竹林中对峙"
    │
    ▼ [Prepare Agent — mini Agent Loop, 最多 5 步]
    │
    ├── Step 1: AI 读取 settings（角色/世界观）
    │           构建**情境因果链**（事件 → 角色反应 → 画面瞬间）
    │           输出**视觉简报**（主体/场景/动作/镜头/画布/光影/nltags）
    │           + 待校验的 hard anchors 列表
    │
    ├── Step 2: 后端调用 danbooru-tags.exe 批量校验
    │           → 回填 confirmed_tags
    │
    ├── Step 3: AI 三层 prompt 组装
    │   ├── hard_tags: quality_prefix + confirmed 角色/画师/外观
    │   ├── soft_phrases: 模型审美短语（不上 danbooru-tags）
    │   └── nltags_block: 空间/动作归属/视线/景深/因果后果
    │
    ├── Step 4: AI 冲突检查 + 负向动态组装
    │   冲突: solo/multiple, close-up/full-body, 室内光源/室外背景...
    │   负向: 按画面风险追加（单人→基础, 多人→防融合, 动态→防肢体断裂）
    │
    └── Step 5: AI 输出结构化 prompt + 项目级 tags + workflow_id + args
    │
    ▼ SSE: image_prepare_done {prompt, negative_prompt, resolution, tags, workflow_id, args}
    │
    ▼ [Generate — 直接 HTTP 调用]
    │
POST {comfyui_base_url}/v1/images/generations
    {model:"anima_pencil", prompt, negative_prompt, size, seed}
    │
    ▼ b64_json → jpg + sidecar → index.json
```

### 9.2 `POST /api/media/image/prepare` 请求

```json
{
    "intent":      "水墨忍者在竹林中对峙",
    "anima":       true,                    // true(默认)=走 SKILL.md Agent 工作流; false=简单 LLM 增强
    "resolution":  "1024x1024",
    "style":       "水墨",
    "chapter_num": 3
}
```

| `anima` | 行为 |
|------|------|
| `true`（默认） | 走 comfyui-animatool SKILL.md 多步 Agent 工作流（情境因果 → 视觉简报 → danbooru-tags 校验 → 三层 prompt → 冲突检查 → 负向动态组装） |
| `false` | 单次 LLM 调用，简单 prompt 增强（现有行为，轻量降级） |

### 9.3 Prepare Agent 系统提示词

`prompts/system/image_prepare_agent.txt` — 基于 `comfyui-good-anima/comfyui-animatool/SKILL.md` 全文改写：

**保留的硬约束**（按 SKILL.md 原文）：
- 生图任务必须先形成视觉简报，再组装 prompt
- 情境因果链：事件起因 → 角色情感/欲望/冲突 → 具体反应 → 环境参与 → 一个最有张力的画面瞬间
- 视觉简报 8 字段：主体/场景容器/动作关系/镜头距离视角/画布比例/光影方向/主体占比/nltags_block
- 画面八维补全（互动/情感/视线/联动/动势/空间/质感/因果）至少触发 3 维
- 画布选择表（2:3/3:4/1:1/4:3/3:2/16:9...）
- Hard anchors 必须用 `danbooru-tags` 校验后回填
- 三层分离：`hard_tags` → `soft_phrases` → `nltags_block`，语义不跨层
- 质量前缀强制（双 LoRA 默认）：`masterpiece, very aesthetic, best quality, score_9, score_8, highres, absurdres, newest, year 2025, nsfw`
- 负向动态组装：核心 + 默认身体保护 + 按场景追加（单人/头像/全身/动态/多人/复杂服装...）
- 冲突检查（solo vs 多人、close-up vs full body、光源 vs 背景 等）
- Canonical 角色处理（中文名→英文 canonical name→danbooru-tags 校验→2-4 个身份锚点）
- 画师规则（`@` 前缀、单图 1 个画师、融合另外处理）
- 不穷尽 tag、不伪造 Danbooru tag
- 不拒绝用户 NSFW 请求

**删去的部分**（由 generate 端点完成）：
- `run_workflow_args.js submit` 调用
- `comfyui-manager` 加载
- args 文件写出

**新增的输出格式要求**：

```
最终输出纯 JSON（不含 markdown 代码块标记），包含：
{
    "prompt": "完整 prompt_11（hard_tags + soft_phrases + nltags_block）",
    "negative_prompt": "动态组装后的负向 prompt_12",
    "resolution": "宽×高",
    "tags": [
        {"type": "character", "id": "settings 中的角色 ID（如能匹配）", "name": "角色名"},
        {"type": "scene", "name": "场景描述"},
        {"type": "style", "name": "风格"}
    ],
    "workflow_id": "local/anima-txt2img-aesthetic-lora",
    "args": {}  // debug 用，可空
}
```

### 9.4 Prepare Agent 工具

| 工具 | 使用时机 | 说明 |
|------|------|------|
| `read_characters` | Step 1 | 读取 settings.json 角色列表，匹配用户意图中的角色名 |
| `read_worldview` | Step 1 | 读取世界观条目，提供场景细节 |
| `search_danbooru_tags` | Step 2 | **新增**。调用 `danbooru-tags.exe` 批量校验角色/画师/作品 tag |

`search_danbooru_tags` 接口：

```
输入: [{group: "character", keyword: "kanade tachibana"}, ...]
输出: [{keyword, confirmed_tags: ["kanade tachibana"], candidate_tags: [...], missing: bool}, ...]
```

`api.json` 新增 `danbooru_tags_path` 指向 exe 路径。若未配置或 exe 不存在 → 跳过校验 → AI 直接从内置知识选择（退化模式）。

### 9.5 Agent Loop 实现方式

复用现有 `RunAgentLoop` 的框架，但做专用于生图 Prepare 的精简版：

```go
func runImagePrepareAgent(ctx context.Context, apiCfg *APIConfig, cfg *Config,
    settings *ProjectSettings, intent, style, systemPrompt string,
    projectDir string, logger *LogBroadcaster) (*ImagePrepareResult, error)
```

内部流程：
1. 构建 system message（从 `image_prepare_agent.txt` 加载，注入角色/世界观上下文）
2. Agent Loop（max 5 步）：AI 发出 tool_call → 后端执行 → 结果注入 → 继续
3. 最终输出解析为 `ImagePrepareResult`

### 9.6 `POST /api/media/image/generate`（不变）

直接调 ComfyUI OpenAI 兼容端点：

```
POST {comfyui_base_url}/v1/images/generations
{model: "anima_pencil", prompt, negative_prompt, size, seed}

↓ b64_json → jpg + sidecar.args.json → index.json
```

### 9.7 `POST /api/media/image/{id}/regenerate`（不变）

读取原图 prompt → 追加 correction → 重新 generate。

### 9.8 数据结构

```go
type ImagePrepareRequest struct {
    Intent     string `json:"intent"`
    Anima      bool   `json:"anima"`                // 默认 true
    Resolution string `json:"resolution,omitempty"`
    Style      string `json:"style,omitempty"`
    ChapterNum int    `json:"chapter_num,omitempty"`
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

type ImageTag struct {
    Type string `json:"type"` // "character" | "scene" | "style" | "worldview" | "custom"
    ID   string `json:"id,omitempty"`
    Name string `json:"name"`
}
```

### 9.9 `anima=false` 退化模式

当 `anima: false` 时，使用当前的单次 LLM 调用逻辑（简单 system prompt + 一次 CallAPI），不做多步 Agent。

---

## 十、Agent 工具（全局助理）

### 10.1 删除（伏笔 5 个）✅

### 10.2 新增（生图 2 个）✅

| 工具名 | 参数 |
|------|------|
| `prepare_image` | `{intent, anima?, style?, resolution?, chapter_num?}` |
| `generate_image` | `{prompt, negative_prompt?, resolution?, tags, count, chapter_num?}` |

## 十一、SSE 事件

| 事件 | 状态 |
|------|------|
| `foreshadow_suggestions` | 已删除 ✅ |
| `foreshadow_outline_conflicts` | 已删除 ✅ |
| `tts_progress` | 待实现 |
| `image_prepare_done` | 待实现（重写后） |
| `image_generated` | ✅ |
| `image_generate_done` | ✅ |

## 十二、前端（阶段 6-7，后续）

- 删 Foreshadows.svelte + i18n 目录
- 删 App.svelte 伏笔/语言切换入口
- 删 api.js 的 X-UI-Locale 头
- 新增 Gallery.svelte（PhotoSwipe）
- 新增 `#gallery` 路由

## 十三、实现计划（剩余）

| 优先级 | 内容 | 说明 |
|------|------|------|
| **P0** | **重写 Prepare Agent** — 按 SKILL.md 实现多步 Agent 工作流 | 生图核心 |
| P1 | Prepare Agent 的 `anima=false` 退化模式 | 已有基础，稍作整理 |
| P1 | `search_danbooru_tags` 工具 | 调用 danbooru-tags.exe |
| P1 | `image_prepare_agent.txt` 默认内容 | 基于 SKILL.md 改写 |
| P2 | TTS SSE 事件 | tts_progress |
| P3 | 前端精简（阶段 6） | 删死代码 |
| P3 | 前端画廊（阶段 7） | Gallery.svelte + PhotoSwipe |

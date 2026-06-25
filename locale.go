package main

import (
	"fmt"
	"net/http"
	"strings"
)

// localeFromRequest extracts the UI locale to use for response messages.
// Priority: X-UI-Locale header > locale query param > Accept-Language > "zh".
func localeFromRequest(r *http.Request) string {
	if r == nil {
		return LangZH
	}
	if v := strings.TrimSpace(r.Header.Get("X-UI-Locale")); v != "" {
		return NormalizeLanguage(v)
	}
	if v := strings.TrimSpace(r.URL.Query().Get("locale")); v != "" {
		return NormalizeLanguage(v)
	}
	if v := strings.TrimSpace(r.Header.Get("Accept-Language")); v != "" {
		first := strings.SplitN(v, ",", 2)[0]
		return NormalizeLanguage(first)
	}
	return LangZH
}

// errorCatalog maps a stable error key to its zh/en messages.
// Messages may contain %s for args.
var errorCatalog = map[string]map[string]string{
	"missing_project_name": {
		LangZH: "缺少项目名称",
	},
	"project_name_invalid_chars": {
		LangZH: "项目名称包含非法字符",
	},
	"project_exists": {
		LangZH: "项目已存在",
	},
	"create_project_dir_failed": {
		LangZH: "创建项目目录失败: %s",
	},
	"init_project_config_failed": {
		LangZH: "初始化项目配置失败: %s",
	},
	"select_project_first": {
		LangZH: "请先选择一个项目",
	},
	"task_running_locked": {
		LangZH: "有AI任务正在运行，暂不能修改，请等待任务完成或先停止任务",
	},
	"task_running_wait": {
		LangZH: "有任务正在运行，请等待完成",
	},
	"no_task_running": {
		LangZH: "没有正在运行的任务",
	},
	"invalid_json": {
		LangZH: "无效的JSON: %s",
	},
	"missing_feedback": {
		LangZH: "缺少 feedback 字段",
	},
	"missing_content": {
		LangZH: "缺少 content 字段",
	},
	"missing_fields": {
		LangZH: "缺少 fields 字段",
	},
	"no_pending_changes": {
		LangZH: "没有待确认的配置变更",
	},
	"load_pending_config_failed": {
		LangZH: "加载待确认配置失败: %s",
	},
	"save_pending_config_failed": {
		LangZH: "保存待确认配置失败: %s",
	},
	"delete_pending_config_failed": {
		LangZH: "清除待确认配置失败: %s",
	},
	"invalid_chapter_num": {
		LangZH: "无效的章节编号",
	},
	"chapter_not_found": {
		LangZH: "章节不存在",
	},
	"chapter_n_not_found": {
		LangZH: "章节 %s 不存在",
	},
	"phase_not_outline": {
		LangZH: "当前不在大纲阶段",
	},
	"phase_not_writing": {
		LangZH: "当前不在写作阶段",
	},
	"outline_empty": {
		LangZH: "大纲为空，请先生成大纲",
	},
	"outline_confirm_failed": {
		LangZH: "确认大纲失败: %s",
	},
	"writing_chapter_present": {
		LangZH: "有正在写作/审核中的章节，请先处理后再重新生成大纲",
	},
	"accepted_chapter_present": {
		LangZH: "存在已确认章节，无法整体重新生成大纲。如需追加章节请使用「生成后续大纲」",
	},
	"writing_chapter_present_delete": {
		LangZH: "有正在写作/审核中的章节，请先处理后再删除大纲",
	},
	"reset_progress_locked": {
		LangZH: "有任务正在运行，无法重置进度",
	},
	"delete_chapter_locked": {
		LangZH: "有任务正在运行，无法删除章节",
	},
	"delete_outline_locked": {
		LangZH: "有任务正在运行，无法删除大纲",
	},
	"delete_project_locked": {
		LangZH: "有任务正在运行，无法删除项目",
	},
	"cannot_delete_current_project": {
		LangZH: "不能删除当前正在使用的项目",
	},
	"project_not_found": {
		LangZH: "项目不存在",
	},
	"delete_project_failed": {
		LangZH: "删除项目失败: %s",
	},
	"delete_progress_failed": {
		LangZH: "删除进度文件失败: %s",
	},
	"no_chapters_to_delete": {
		LangZH: "没有可删除的章节",
	},
	"writing_chapter_cannot_delete": {
		LangZH: "正在写作中的章节无法删除",
	},
	"delete_frontier_unavailable": {
		LangZH: "当前写作前沿没有可删除的章节正文",
	},
	"writing_range_has_writing": {
		LangZH: "删除范围内有正在写作中的章节，无法删除",
	},
	"save_progress_failed": {
		LangZH: "保存进度失败: %s",
	},
	"save_failed": {
		LangZH: "保存失败: %s",
	},
	"save_config_failed": {
		LangZH: "保存配置失败: %s",
	},
	"save_api_config_failed": {
		LangZH: "保存API配置失败: %s",
	},
	"serialize_config_failed": {
		LangZH: "序列化配置失败: %s",
	},
	"serialize_api_config_failed": {
		LangZH: "序列化API配置失败: %s",
	},
	"api_test_timeout": {
		LangZH: "连接超时（15秒）",
	},
	"api_test_failed": {
		LangZH: "测试失败: %s",
	},
	"api_test_success": {
		LangZH: "连接成功",
	},
	"character_name_empty": {
		LangZH: "角色名不能为空",
	},
	"character_not_found": {
		LangZH: "角色不存在",
	},
	"worldview_field_empty": {
		LangZH: "名称和描述不能为空",
	},
	"worldview_not_found": {
		LangZH: "世界观条目不存在",
	},
	"organization_name_empty": {
		LangZH: "组织名不能为空",
	},
	"organization_not_found": {
		LangZH: "组织不存在",
	},
	"relation_endpoints_empty": {
		LangZH: "源和目标不能为空",
	},
	"relation_not_found": {
		LangZH: "关系不存在",
	},
	"foreshadow_name_required": {
		LangZH: "缺少 name",
	},
	"foreshadow_desc_required": {
		LangZH: "缺少 description",
	},
	"foreshadow_not_found": {
		LangZH: "伏笔不存在",
	},
	"invalid_foreshadow_id": {
		LangZH: "无效的伏笔ID",
	},
	"need_generate_outline_first": {
		LangZH: "请先生成大纲",
	},
	"continue_reset_first": {
		LangZH: "续写前请先重置进度",
	},
	"continue_analyze_first": {
		LangZH: "请先分析内容",
	},
	"analysis_no_chapters": {
		LangZH: "分析结果中没有任何章节",
	},
	"continue_import_failed": {
		LangZH: "导入续写失败: %s",
	},
	"book_not_complete": {
		LangZH: "全书尚未完成（需所有章节已确认）",
	},
	"need_polish_skill": {
		LangZH: "没有启用的润色技能，请先在技能管理页启用 polish 类技能",
	},
	"chapter_content_empty": {
		LangZH: "章节内容为空，无法润色",
	},
	"chapter_edit_op_required": {
		LangZH: "缺少 operation 参数，必须为 replace_lines / replace_text / insert_after_line / append 之一",
	},
	"chapter_edit_text_required": {
		LangZH: "new_text 不能为空",
	},
	"chapter_edit_failed": {
		LangZH: "章节编辑失败: %s",
	},
	"chapter_in_writing": {
		LangZH: "章节正在写作中，无法润色",
	},
	"chapter_num_required": {
		LangZH: "请指定章节编号",
	},
	"no_transitions_to_optimize": {
		LangZH: "没有可优化的章节（需要至少两个相邻的已确认章节）",
	},
	"missing_diagnosis_or_consistency": {
		LangZH: "缺少诊断或核查报告，请先运行全书诊断",
	},
	"no_roadmap_items": {
		LangZH: "没有可执行的优化工单",
	},
	"select_at_least_one_item": {
		LangZH: "请至少勾选一条待执行的工单",
	},
	"clear_postprocess_failed": {
		LangZH: "清空失败: %s",
	},
	"chat_session_not_found": {
		LangZH: "会话不存在",
	},
	"load_session_list_failed": {
		LangZH: "加载会话列表失败: %s",
	},
	"create_session_failed": {
		LangZH: "创建会话失败: %s",
	},
	"save_session_failed": {
		LangZH: "保存会话失败: %s",
	},
	"delete_session_failed": {
		LangZH: "删除会话失败: %s",
	},
	"skill_not_found": {
		LangZH: "技能不存在",
	},
	"settings_ai_generate_moved": {
		LangZH: "此功能已移至 LLM 对话中，请通过聊天让 AI 帮你生成设定",
	},
	"settings_polish_moved": {
		LangZH: "此功能已移至 LLM 对话中，请通过聊天让 AI 帮你润色",
	},
	"writing_conflict_none": {
		LangZH: "当前没有待处理的写作冲突",
	},
	"missing_action": {
		LangZH: "缺少 action 字段",
	},
	"invalid_conflict_chapter_idx": {
		LangZH: "冲突章节索引无效",
	},
	"unsupported_action": {
		LangZH: "不支持的 action: %s",
	},
	"no_foreshadows_to_check": {
		LangZH: "当前没有伏笔，无需检查",
	},
}

// systemPrompts maps a stable AI-system-prompt key to per-language text.
// These appear in api calls (CallAPI(ctx, cfg, systemPrompt, userPrompt)) and must
// be language-aware so an English project doesn't get a Chinese system role.
var systemPrompts = map[string]map[string]string{
	"outline_editor_json": {
		LangZH: "你是一位专业的小说策划编辑。请严格按照要求的JSON格式输出，不要添加任何额外文字或markdown代码块标记。",
	},
	"outline_editor_locked_json": {
		LangZH: "你是一位小说策划编辑。请严格按照要求的JSON格式输出，不要添加任何额外文字或markdown代码块标记。已锁定的章节内容不可修改。",
	},
	"outline_editor_brief_json": {
		LangZH: "你是一位严谨的小说策划编辑。请严格按照要求的JSON格式输出，不要添加任何额外文字。",
	},
	"summary_analyst": {
		LangZH: "你是一位精准的小说叙事状态分析师。",
	},
	"fact_checker_json": {
		LangZH: "你是一位严谨的小说事实核查员。请严格按照要求的JSON格式输出。",
	},
	"narrative_architect_json": {
		LangZH: "你是一位资深的小说叙事架构师。请严格按照要求的JSON格式输出，不要添加任何额外文字或markdown代码块标记。",
	},
	"foreshadow_tracker_json": {
		LangZH: "你是一位严谨的小说伏笔追踪员。请严格按照要求的JSON格式输出，不要添加任何额外文字或markdown代码块标记。",
	},
	"foreshadow_outline_checker_json": {
		LangZH: "你是一位严谨的小说叙事一致性编辑。请严格按照要求的JSON格式输出，不要添加任何额外文字。拿不准时视为无冲突。",
	},
	"outline_character_checker_json": {
		LangZH: "你是一位严谨的小说设定编辑。请严格按照要求的JSON格式输出，不要添加任何额外文字。拿不准时视为无未登记人物。",
	},
	"writing_conflict_analyst_json": {
		LangZH: "你是一位资深小说编辑，擅长诊断大纲、伏笔与前情之间的矛盾。请严格按照要求的JSON格式输出，不要添加任何额外文字。",
	},
	"consistency_reviewer_json": {
		LangZH: "你是一位专业的小说一致性审查编辑。请严格按照要求的JSON格式输出，不要添加任何额外文字或markdown代码块标记。",
	},
	"content_analyst_json": {
		LangZH: "你是一位专业的小说分析编辑。请严格按照要求的JSON格式输出，不要添加任何额外文字或markdown代码块标记。",
	},
	"transition_editor": {
		LangZH: "你是一位资深小说编辑，擅长打磨章节之间的衔接。请严格按要求输出。",
	},
	"polish_editor": {
		LangZH: "你是一位专业的中文小说润色编辑。请严格按照规则修改文本，输出修改后的完整章节正文。不要添加章节标题、章节号、「本章完」等任何解释、标记或元信息。",
	},
	"book_diagnosis": {
		LangZH: "你是一位资深网文总编辑，擅长长篇完稿后的通读审阅。请严格按要求输出诊断报告，不要改写正文。",
	},
	"book_consistency_check": {
		LangZH: "你是一位严谨的小说事实核查员。请输出结构化核查报告，不要改写正文。",
	},
	"book_roadmap": {
		LangZH: "你是一位资深小说编辑。请根据报告生成可执行的修改工单 JSON，不要输出正文改写。",
	},
	"author_default": {
		LangZH: "你是一位小说作者。只输出小说正文，不要输出章节标题、章节号、作者说明或「本章完」等元信息。严格保持用户指定的叙述视角统一。",
	},
	"chapter_revision_suffix": {
		LangZH: "\n你正在执行章节修订任务：只做修改意见要求的改动，其余原文保持不变，输出修改后的完整正文；不要添加任何元信息或说明性文字。",
	},
	"memory_manager": {
		LangZH: "你是一位精准的小说叙事记忆管理员。请严格按照要求的JSON格式输出，不要添加任何额外文字。",
	},
	"image_prepare_agent": {
		LangZH: "",
	},
}

// SystemPromptFor returns the AI system-prompt for the given key.
// Reads from prompts/system/{key}.txt if available; falls back to built-in default.
func SystemPromptFor(projectDir, key string) string {
	return loadSystemPrompt(projectDir, key)
}

func lookupCatalog(lang, key string) (string, bool) {
	lang = NormalizeLanguage(lang)
	for _, catalog := range []map[string]map[string]string{messageCatalog, errorCatalog} {
		entry, ok := catalog[key]
		if !ok {
			continue
		}
		tpl := entry[lang]
		if tpl == "" {
			tpl = entry[LangZH]
		}
		if tpl != "" {
			return tpl, true
		}
	}
	return "", false
}

// T returns a localized message for the given key and args; falls back to zh, then key.
func T(lang, key string, args ...any) string {
	tpl, ok := lookupCatalog(lang, key)
	if !ok {
		return key
	}
	if len(args) == 0 {
		return tpl
	}
	return fmt.Sprintf(tpl, args...)
}

func msgArgsToStrings(args ...any) []string {
	if len(args) == 0 {
		return nil
	}
	out := make([]string, len(args))
	for i, a := range args {
		out[i] = fmt.Sprint(a)
	}
	return out
}

// writeErrorReq writes a JSON error response, picking message language from the request.
func (h *Handlers) writeErrorReq(w http.ResponseWriter, r *http.Request, code int, key string, args ...any) {
	lang := localeFromRequest(r)
	h.writeJSON(w, code, map[string]string{"error": T(lang, key, args...)})
}

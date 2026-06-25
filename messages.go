package main

// messageCatalog holds localized UI/log/agent status strings (key → zh/en template).
// Templates use fmt.Sprintf verbs (%s, %d, %v). Frontend mirrors keys with {0},{1},… placeholders.
var messageCatalog = map[string]map[string]string{
	// ---- Task / handler logs ----
	"log.autoconfirm_on": {
		LangZH: "已开启自动确认模式：每章生成完成后将自动确认并继续生成下一章",
	},
	"log.autoconfirm_off": {
		LangZH: "已关闭自动确认模式",
	},
	"log.outline_cleared_pending": {
		LangZH: "已自动清除旧的大纲（pending 章节）",
	},
	"log.outline_generating": {
		LangZH: "正在生成小说大纲...",
	},
	"log.outline_generate_cancelled": {
		LangZH: "大纲生成已取消",
	},
	"log.outline_generate_failed": {
		LangZH: "大纲生成失败: %s",
	},
	"log.outline_generate_done": {
		LangZH: "大纲生成完成！",
	},
	"log.outline_confirmed": {
		LangZH: "大纲已确认，进入写作阶段。",
	},
	"log.outline_revising": {
		LangZH: "正在根据意见修订大纲...",
	},
	"log.outline_revise_cancelled": {
		LangZH: "大纲修订已取消",
	},
	"log.outline_revise_failed": {
		LangZH: "大纲修订失败: %s",
	},
	"log.outline_revised": {
		LangZH: "大纲已修订。",
	},
	"log.chapter_writing": {
		LangZH: "正在创作第 %d 章...",
	},
	"log.chapter_write_cancelled": {
		LangZH: "章节创作已取消",
	},
	"log.chapter_write_conflict_pause": {
		LangZH: "章节创作因事实核查冲突暂停，等待你选择处理方向",
	},
	"log.chapter_write_failed": {
		LangZH: "章节创作失败: %s",
	},
	"log.chapter_write_done": {
		LangZH: "第 %d 章《%s》创作完成！",
	},
	"log.chapter_autoconfirm_failed": {
		LangZH: "自动确认失败: %s",
	},
	"log.chapter_autoconfirmed": {
		LangZH: "第 %d 章《%s》已自动确认。",
	},
	"log.all_chapters_done": {
		LangZH: "全部章节已创作完成！",
	},
	"log.autowrite_cancelled": {
		LangZH: "任务已取消，停止自动续写",
	},
	"log.chapter_kept_review": {
		LangZH: "第 %d 章已保留当前稿并进入审核。",
	},
	"log.chapter_confirmed": {
		LangZH: "第 %d 章已确认。",
	},
	"log.chapter_revising": {
		LangZH: "正在根据意见修改当前章节...",
	},
	"log.chapter_revise_cancelled": {
		LangZH: "章节修订已取消",
	},
	"log.chapter_revise_failed": {
		LangZH: "章节修订失败: %s",
	},
	"log.chapter_revised": {
		LangZH: "章节已修订。",
	},
	"log.chapter_specific_revising": {
		LangZH: "正在定向修订第 %d 章...",
	},
	"log.smooth_transitions_cancelled": {
		LangZH: "章节衔接优化已取消（已完成部分不会丢失）",
	},
	"log.smooth_transitions_failed": {
		LangZH: "章节衔接优化失败: %s",
	},
	"log.chapter_deleted": {
		LangZH: "已删除第 %d 章。",
	},
	"log.outline_deleted": {
		LangZH: "大纲已删除。",
	},
	"log.chapter_outline_updated": {
		LangZH: "第 %d 章大纲已更新。",
	},
	"log.settings_reconciling": {
		LangZH: "正在协调新设定与已有内容...",
	},
	"log.settings_reconcile_cancelled": {
		LangZH: "设定协调已取消",
	},
	"log.settings_reconcile_failed": {
		LangZH: "设定协调失败: %s",
	},
	"log.settings_reconcile_done": {
		LangZH: "设定协调完成！",
	},
	"log.delete_file_failed": {
		LangZH: "删除文件 %s 失败: %v",
	},
	"log.chapters_deleted_from": {
		LangZH: "已从第 %d 章删除到末尾，共删除 %d 章。",
	},
	"log.foreshadow_roadmap_save_failed": {
		LangZH: "伏笔路线图保存失败: %v",
	},
	"log.foreshadow_suggesting": {
		LangZH: "正在分析大纲，设计伏笔方案...",
	},
	"log.foreshadow_suggest_cancelled": {
		LangZH: "伏笔建议已取消",
	},
	"log.foreshadow_suggest_failed": {
		LangZH: "伏笔建议生成失败: %s",
	},
	"log.foreshadow_suggest_done": {
		LangZH: "伏笔建议生成完成，共 %d 条",
	},
	"log.continue_analyzing": {
		LangZH: "正在分析已有内容...",
	},
	"log.continue_analyze_cancelled": {
		LangZH: "内容分析已取消",
	},
	"log.continue_analyze_failed": {
		LangZH: "内容分析失败: %s",
	},
	"log.continue_analyze_done": {
		LangZH: "内容分析完成，发现 %d 章",
	},
	"log.continue_import_done": {
		LangZH: "续写导入完成，已进入大纲阶段。",
	},
	"log.continuation_outline_generating": {
		LangZH: "正在生成续写大纲...",
	},
	"log.continuation_outline_cancelled": {
		LangZH: "续写大纲生成已取消",
	},
	"log.continuation_outline_failed": {
		LangZH: "续写大纲生成失败: %s",
	},
	"log.continuation_outline_done": {
		LangZH: "续写大纲生成完成！",
	},
	"log.chapter_polish_cancelled": {
		LangZH: "章节润色已取消",
	},
	"log.chapter_polish_failed": {
		LangZH: "章节润色失败: %s",
	},
	"log.postprocess_diagnose_cancelled": {
		LangZH: "全书优化分析已取消",
	},
	"log.postprocess_diagnose_failed": {
		LangZH: "全书优化分析失败: %s",
	},
	"log.postprocess_consistency_cancelled": {
		LangZH: "全书一致性核查已取消",
	},
	"log.postprocess_consistency_failed": {
		LangZH: "全书一致性核查失败: %s",
	},
	"log.postprocess_roadmap_cancelled": {
		LangZH: "路线图生成已取消",
	},
	"log.postprocess_roadmap_failed": {
		LangZH: "路线图生成失败: %s",
	},
	"log.postprocess_execute_cancelled": {
		LangZH: "全书优化执行已取消（已完成项不会丢失）",
	},
	"log.postprocess_execute_failed": {
		LangZH: "全书优化执行失败: %s",
	},
	"log.child_task_start_failed": {
		LangZH: "无法启动子任务 %s：主任务已结束",
	},
	"log.save_session_failed": {
		LangZH: "保存会话失败: %v",
	},
	"log.chat_cancelled": {
		LangZH: "助理对话已取消",
	},
	"log.chat_failed": {
		LangZH: "助理回复失败: %v",
	},
	"log.chat_done": {
		LangZH: "助理回复完成",
	},
	"log.project_deleted": {
		LangZH: "项目「%s」已删除",
	},
	"log.project_created": {
		LangZH: "项目「%s」创建成功",
	},

	// ---- Writing pipeline logs ----
	"log.chapter_start": {
		LangZH: "开始创作第 %d 章: 《%s》",
	},
	"log.outline_check_failed": {
		LangZH: "大纲一致性检查失败: %v（按原大纲继续）",
	},
	"log.outline_auto_revised": {
		LangZH: "本章大纲已自动修订以匹配当前剧情",
	},
	"log.outline_consistent": {
		LangZH: "本章大纲与当前剧情一致 ✓",
	},
	"log.prose_done": {
		LangZH: "正文撰写完毕，共 %d 字",
	},
	"log.chapter_length_retry": {
		LangZH: "[字数控制] 当前 %d 字，要求 %d–%d 字，正在第 %d 次重新撰写...",
	},
	"log.chapter_length_adjust": {
		LangZH: "[字数控制] 重写后仍为 %d 字（要求 %d–%d 字），正在对原文压缩/扩展...",
	},
	"log.chapter_length_adjust_failed": {
		LangZH: "[字数控制] 压缩/扩展失败: %v，保留最佳稿",
	},
	"log.chapter_length_soft_keep": {
		LangZH: "[字数控制] 最佳稿 %d 字略超/略低于 %d–%d 字，在容忍范围内，跳过压缩/扩展",
	},
	"log.chapter_length_skip_adjust": {
		LangZH: "[字数控制] 最佳稿 %d 字偏离 %d–%d 字过多，跳过压缩/扩展",
	},
	"log.chapter_length_adjust_reverted": {
		LangZH: "[字数控制] 压缩/扩展后 %d 字未优于最佳稿 %d 字（要求 %d–%d 字），保留最佳稿",
	},
	"log.chapter_length_off_range": {
		LangZH: "[字数控制] 压缩/扩展后仍为 %d 字（要求 %d–%d 字），请用户在审核时处理；自动确认模式下将继续下一章",
	},
	"log.summary_done": {
		LangZH: "摘要提炼完成",
	},
	"log.factcheck_retry": {
		LangZH: "[事实核查] 发现问题，正在重新生成第 %d 章（第 %d 次重试）...",
	},
	"log.factcheck_details": {
		LangZH: "核查详情: %s",
	},
	"log.factcheck_max_retries": {
		LangZH: "[事实核查] 已达最大重试次数，正在分析冲突根因...",
	},
	"log.conflict_analyze_failed": {
		LangZH: "冲突分析失败: %v，保留当前版本",
	},
	"log.conflict_retry": {
		LangZH: "检测到可调和冲突，正在按补充约束进行最后一次尝试...",
	},
	"log.factcheck_constraint_pass": {
		LangZH: "[事实核查] 补充约束尝试通过 ✓",
	},
	"log.factcheck_pass": {
		LangZH: "[事实核查] 通过 ✓",
	},
	"log.chapter_write_complete": {
		LangZH: "第 %d 章创作完成！",
	},
	"log.outline_conflict": {
		LangZH: "第 %d 章大纲与当前剧情冲突: %s",
	},
	"log.chapter_modifying": {
		LangZH: "正在修改第 %d 章《%s》...",
	},
	"log.prose_revised": {
		LangZH: "正文修改完毕，共 %d 字",
	},
	"log.subsequent_outline_failed": {
		LangZH: "后续大纲修订失败: %v（不影响当前章节）",
	},
	"log.subsequent_outline_done": {
		LangZH: "后续大纲修订完成",
	},
	"log.chapter_specific_revising_long": {
		LangZH: "正在对第 %d 章《%s》进行定向修订（不影响其他章节）...",
	},
	"log.prose_specific_revised": {
		LangZH: "正文修订完毕，共 %d 字",
	},
	"log.chapter_specific_done": {
		LangZH: "第 %d 章定向修订完成（其余章节未受影响）。",
	},
	"log.fatal_no_retry": {
		LangZH: "致命错误: %v，不再重试",
	},
	"log.content_gen_retry": {
		LangZH: "正文生成失败: %v。第 %d 次重试，等待 %ds...",
	},
	"log.summary_retry": {
		LangZH: "摘要提炼失败: %v。第 %d 次重试，等待 %ds...",
	},
	"log.factcheck_api_retry": {
		LangZH: "事实核查失败: %v。第 %d 次重试，等待 %ds...",
	},
	"log.smooth_start": {
		LangZH: "开始章节衔接优化，共 %d 章待检查",
	},
	"log.smooth_natural": {
		LangZH: "第 %d 章衔接自然，无需修改",
	},
	"log.smooth_optimized": {
		LangZH: "第 %d 章开头已优化并保存",
	},
	"log.smooth_done": {
		LangZH: "章节衔接优化完成：检查 %d 章，优化 %d 章",
	},
	"log.outline_generate_summary": {
		LangZH: "大纲生成完成，共 %d 章，标题: 《%s》",
	},
	"log.outline_revise_summary": {
		LangZH: "大纲已修订，共 %d 章",
	},
	"log.reconcile_pending_outline_failed": {
		LangZH: "待定章节大纲重新生成失败: %v（设定已更新）",
	},
	"log.reconcile_done_explain": {
		LangZH: "设定协调完成。%s",
	},
	"log.continuation_outline_summary": {
		LangZH: "续写大纲生成完成，新增 %d 章，总计 %d 章",
	},
	"log.foreshadow_outline_check_failed": {
		LangZH: "伏笔-大纲一致性检查失败: %v",
	},
	"log.foreshadow_outline_report_save_failed": {
		LangZH: "保存伏笔-大纲检查报告失败: %v",
	},
	"log.foreshadow_outline_check_pass": {
		LangZH: "伏笔与大纲一致性检查通过 ✓",
	},
	"log.outline_chapters_too_short": {
		LangZH: "第 %s 章大纲不足 %d 字，正在要求 AI 扩写…",
	},
	"log.outline_chapters_still_short": {
		LangZH: "第 %s 章大纲仍不足 %d 字（已重试），请在大纲页手动补充",
	},
	"log.outline_character_check_failed": {
		LangZH: "大纲人物检查失败: %v",
	},
	"log.outline_character_report_save_failed": {
		LangZH: "保存大纲人物检查报告失败: %v",
	},
	"log.outline_character_check_pass": {
		LangZH: "大纲人物与已登记角色一致 ✓",
	},
	"log.foreshadow_plan_parsed": {
		LangZH: "伏笔方案解析完成，共 %d 条",
	},
	"log.foreshadow_status_updated": {
		LangZH: "伏笔状态更新完成，处理 %d 条变更",
	},
	"log.foreshadow_sync_failed": {
		LangZH: "伏笔状态更新失败: %v（不影响本章）",
	},
	"log.foreshadow_sync_summary": {
		LangZH: "伏笔状态已更新（活跃: %d, 已回收: %d）",
	},
	"log.memory_update_failed": {
		LangZH: "叙事记忆更新失败（不影响本章）",
	},
	"log.memory_save_failed": {
		LangZH: "叙事记忆保存失败",
	},
	"log.postprocess_material": {
		LangZH: "全书材料：约 %d 字，预估 %d tokens，诊断模式：%s",
	},
	"log.postprocess_consistency_single": {
		LangZH: "开始全书一致性核查（单卷）...",
	},
	"log.postprocess_consistency_multi": {
		LangZH: "正文较长，分 %d 卷进行一致性核查...",
	},
	"log.postprocess_roadmap_items": {
		LangZH: "已生成 %d 条优化工单",
	},
	"log.postprocess_smooth_preface": {
		LangZH: "前置步骤：优化章节衔接...",
	},
	"log.postprocess_smooth_skip": {
		LangZH: "章节衔接优化跳过或失败: %v",
	},
	"log.postprocess_batch_failed": {
		LangZH: "第 %d 章工单失败: %v",
	},
	"log.postprocess_batch_done": {
		LangZH: "第 %d 章已完成（合并 %d 条意见）",
	},
	"log.postprocess_batch_skip": {
		LangZH: "第 %d 章内容无变化，已跳过",
	},
	"log.postprocess_execute_summary": {
		LangZH: "全书优化执行完成：处理 %d 章（共 %d 条工单），有效修改 %d 章",
	},
	"log.api_fatal": {
		LangZH: "致命错误: %v，不再重试",
	},
	"log.api_retry": {
		LangZH: "API调用失败: %v。第 %d 次重试，等待 %ds...",
	},
	"log.api_stream_retry": {
		LangZH: "流式API调用失败: %v。第 %d 次重试，等待 %ds...",
	},

	// ---- Agent tool status messages ----
	"agent.task_cancelled": {
		LangZH: "任务已取消",
	},
	"agent.api_failed": {
		LangZH: "Agent API 调用失败: %v",
	},
	"agent.max_steps": {
		LangZH: "已达到最大工具调用步骤限制。",
	},
	"agent.tool_exec_error": {
		LangZH: "工具执行错误: %v",
	},
	"agent.unknown_tool": {
		LangZH: "未知工具: %s",
	},
	"agent.confirm_required": {
		LangZH: "⚠️ 操作未执行：「%s」是不可逆的危险操作。请先向用户复述影响范围并获得明确同意，确认后携带 confirm=true 重新调用。如果用户的本意是修改内容而非删除，请改用对应的修订工具。",
	},
	"agent.no_characters": {
		LangZH: "暂无角色数据",
	},
	"agent.characters_not_found": {
		LangZH: "没有找到匹配的角色",
	},
	"agent.character_not_found": {
		LangZH: "未找到角色: %s",
	},
	"agent.no_worldview": {
		LangZH: "暂无世界观数据",
	},
	"agent.worldview_not_found": {
		LangZH: "没有找到匹配的世界观条目",
	},
	"agent.no_organizations": {
		LangZH: "暂无组织数据",
	},
	"agent.chapter_not_found": {
		LangZH: "未找到第%d章",
	},
	"agent.no_outline": {
		LangZH: "暂无大纲",
	},
	"agent.no_foreshadows": {
		LangZH: "暂无伏笔",
	},
	"agent.search_keyword_required": {
		LangZH: "请提供搜索关键词",
	},
	"agent.search_no_results": {
		LangZH: "未找到相关内容",
	},
	"agent.character_created": {
		LangZH: "角色「%s」创建成功 (ID: %s)",
	},
	"agent.character_updated": {
		LangZH: "角色「%s」已更新",
	},
	"agent.character_deleted": {
		LangZH: "角色「%s」已删除",
	},
	"agent.worldview_created": {
		LangZH: "世界观条目「%s」创建成功 (ID: %s)",
	},
	"agent.worldview_updated": {
		LangZH: "世界观条目「%s」已更新",
	},
	"agent.worldview_deleted": {
		LangZH: "世界观条目「%s」已删除",
	},
	"agent.config_saved_reconciling": {
		LangZH: "故事配置已保存，正在自动协调已有内容...",
	},
	"agent.config_saved": {
		LangZH: "故事配置已保存",
	},
	"agent.outline_task_started": {
		LangZH: "大纲生成任务已启动，请等待完成。",
	},
	"agent.outline_confirmed": {
		LangZH: "大纲已确认，现在进入写作阶段。",
	},
	"agent.outline_revise_started": {
		LangZH: "大纲修订任务已启动，请等待完成。",
	},
	"agent.outline_deleted": {
		LangZH: "大纲已删除。",
	},
	"agent.chapter_outline_updated": {
		LangZH: "第 %d 章大纲已更新。",
	},
	"agent.chapter_task_started": {
		LangZH: "第 %d 章生成任务已启动，请等待完成。",
	},
	"agent.chapter_confirmed": {
		LangZH: "第 %d 章《%s》已确认。",
	},
	"agent.chapter_revise_started": {
		LangZH: "第 %d 章修订任务已启动（仅修改该章，不影响其他章节），请等待完成。",
	},
	"agent.chapter_deleted": {
		LangZH: "已删除第 %d 章。",
	},
	"agent.chapter_content_edited": {
		LangZH: "第 %d 章正文已编辑（操作: %s，共 %d 行）。",
	},
	"agent.chapter_edit_op_required": {
		LangZH: "缺少 operation 参数，必须为 replace_lines / replace_text / insert_after_line / append 之一",
	},
	"agent.chapter_edit_text_required": {
		LangZH: "new_text 不能为空",
	},
	"agent.chapters_bulk_delete_confirm": {
		LangZH: "⚠️ 操作未执行：这将永久删除第 %d 章到末尾共 %d 章的全部内容。请先向用户复述此影响范围并获得明确同意，确认后携带 confirm=true 重新调用。如果用户的本意是修改章节内容，请改用 revise_chapter。",
	},
	"agent.chapters_deleted_from": {
		LangZH: "已从第 %d 章删除到末尾，共删除 %d 章。",
	},
	"agent.organization_created": {
		LangZH: "组织「%s」创建成功 (ID: %s)",
	},
	"agent.organization_updated": {
		LangZH: "组织「%s」已更新",
	},
	"agent.organization_deleted": {
		LangZH: "组织「%s」已删除",
	},
	"agent.organization_not_found": {
		LangZH: "未找到组织: %s",
	},
	"agent.relation_created": {
		LangZH: "关系创建成功 (ID: %s)",
	},
	"agent.relation_updated": {
		LangZH: "关系已更新 (ID: %s)",
	},
	"agent.relation_deleted": {
		LangZH: "关系已删除",
	},
	"agent.relation_not_found": {
		LangZH: "未找到关系: %s",
	},
	"agent.foreshadow_suggest_started": {
		LangZH: "伏笔建议生成任务已启动，请等待完成。",
	},
	"agent.foreshadow_created": {
		LangZH: "伏笔「%s」创建成功 (ID: %d)",
	},
	"agent.foreshadow_updated": {
		LangZH: "伏笔「%s」已更新",
	},
	"agent.foreshadow_deleted": {
		LangZH: "伏笔「%s」已删除",
	},
	"agent.foreshadow_not_found": {
		LangZH: "伏笔 %d 不存在",
	},
	"agent.skill_toggled": {
		LangZH: "技能「%s」已%s",
	},
	"agent.progress_reset": {
		LangZH: "进度已重置。",
	},
	"agent.image_prepare_started": {
		LangZH: "图片提示词准备任务已启动。",
	},
}

package constant

const (
	PlayerRoleMenu        = "扮演角色:(\\S+)"
	GptModuleMenu         = "model:(\\S+)"
	GptGroupMemory        = "全群记忆"
	GptGroupOneselfMemory = "独立记忆"

	CacheSystemRoleKey  = "system_role_%s"
	CacheGptModuleKey   = "gpt_module_%s"
	CacheGroupMemoryKey = "gpt_memory_%s"

	BingModule    = "bing"
	QwenModule    = "qwen"
	GeminitModule = "geminit"

	ClearHistoryKey = "让我们重新开始"
	ErrorReply      = "对不起，我不太想跟你再聊这个了,让我们重新开始吧!"

	ChatSummaryKey = "结合历史聊天"

	MaxPromptLen = 30000
)

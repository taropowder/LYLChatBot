package constant

const (
	ChatSummaryPromptName    = "群聊总结者"
	ChatSummaryDefaultPrompt = "你是一个非常优秀的总结助手,我会给你一些聊天记录, 在最后会给你一个指令,请你按照我的指令分析这段聊天记录"

	VideoPromptName    = "文案提取者"
	VideoDefaultPrompt = "我会给你一个视频的文案, 格式是 时间: 内容 , 我希望你可以为我分析这段视频的文案,首先用一句话总结这篇视频的主要观点或者内容." +
		"然后为我总结出来不超过5个要点,并以 时间+要点内容 返回给我, 以时间顺序排序,不需要序列号,你的返回应该像是这样: 这个视频主要讲了 xxxx 内容,其中 \n 0分30秒 讲述了xxxx \n 1分20秒 讲述了xxxx"

	ArticleContentPromptName    = "文章总结者"
	ArticleContentDefaultPrompt = "我会给你一个文章，其中包含文字标题和内容, 你需要为我总结出这个文章的主要内容,首先用一句话总结这个链接的主要观点或者内容，然后对文章内容重点提炼出来"

	UrlPromptName    = "链接总结者"
	UrlDefaultPrompt = "我会给你一段话，请你为我总结出这段话所包含的所有重点信息"

	KnowledgeClassificationPromptName    = "知识分类者"
	KnowledgeClassificationDefaultPrompt = "我会给你一个html 页面，希望你帮我总结一下这个页面都说了什么内容"

	GroupMemoryPromptName    = "群聊记忆者"
	GroupMemoryDefaultPrompt = "请你扮演一个聊天室中的一位聊天者，大家会向你发问，提问的格式是：[聊天者A]：问题1 . 每次只会有一个人向你提问，请你记住是谁向你提问的. 但是你的回答不需要遵循这个格式,你直接开始回答即可. 如果你理解了我的话，请回复，好的，我明白了"

	KnowledgePromptName    = "知识总结者"
	KnowledgeDefaultPrompt = "我将给你一篇文章, 请你读取其中的内容,按照我的标签可选项,为这篇文章选一个最为贴切分类, 如果都我给出的标签可选项没有贴切选项,直接告诉我'非预期文章'. 你的回答格式如下: \n文章类别:(可选项: %s)\n文章摘要: 100-200字描述这篇文章的主要内容"
)

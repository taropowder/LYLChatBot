package model

import "gorm.io/gorm"

type KnowledgeRecord struct {
	gorm.Model

	KnowledgeType     string `gorm:"type:varchar(128);"`
	KnowledgeUrl      string `gorm:"type:text;"`
	KnowledgeAbstract string `gorm:"type:text;"`
	KnowledgeSharer   string `gorm:"type:varchar(128);"`
}

// 我将给你一篇文章, 请帮我按照我所预期的格式, 提取其中的关键信息给我
// 文章类型:(可选项:网络安全新闻/网络安全技术/网络安全工具/渗透经验)
// 文章标签:(可选项: Linux/Windows/渗透技巧)
// 文章摘要: 用不多于200字描述这篇文章的主要内容

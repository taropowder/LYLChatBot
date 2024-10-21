package model

import "gorm.io/gorm"

type KnowledgeQARecord struct {
	gorm.Model

	Knowledge   KnowledgeRecord `gorm:"foreignKey:KnowledgeId"`
	KnowledgeId uint
	Question    string `gorm:"type:text;"`
	Answer      string `gorm:"type:text;"`
	Questioner  string `gorm:"type:varchar(512);"`
	IsFirst     bool   `gorm:"type:boolean;"`
}

// 我将给你一篇文章, 请帮我按照我所预期的格式, 提取其中的关键信息给我
// 文章类型:(可选项:网络安全新闻/网络安全技术/网络安全工具/渗透经验)
// 文章摘要: 用不多于200字描述这篇文章的主要内容

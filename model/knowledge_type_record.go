package model

import "gorm.io/gorm"

type KnowledgeTypeRecord struct {
	gorm.Model

	TypeName string `gorm:"type:varchar(128);"`
}

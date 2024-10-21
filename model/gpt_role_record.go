package model

import "gorm.io/gorm"

type GptRoleRecord struct {
	gorm.Model

	RoleName string `gorm:"type:varchar(128);"`
	Prompt   string `gorm:"type:text;"`
}

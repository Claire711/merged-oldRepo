package database

import (
	"github.com/gofrs/uuid"
	"server/global"
)

// Feedback 反馈表
type Feedback struct {
	global.MODEL
	UserUUID uuid.UUID `json:"user_uuid" gorm:"type:char(36)"`
	User     User      `json:"-" gorm:"foreignKey:UserUUID;references:UUID"`
	Content  string    `json:"content"`
	Reply    string    `json:"reply"`
}

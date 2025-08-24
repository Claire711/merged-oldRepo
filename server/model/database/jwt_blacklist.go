package database

import "server/global"

// JWT 令牌的黑名单管理
type JwtBlacklist struct {
	global.MODEL
	Jwt string `json:"jwt" gorm:"type:text"`
}

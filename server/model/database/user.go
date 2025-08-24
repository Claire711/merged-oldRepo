package database

import (
	"github.com/gofrs/uuid"
	"server/global"
	"server/model/appTypes"
	//包uuid提供了通用唯一标识符的纯Go实现 （UUID）变体，定义在RFC-9562中。这个包同时支持创建 以及解析不同格式的uid
)

type User struct {
	global.MODEL
	UUID      uuid.UUID         `json:"uuid" gorm:"type:char(36);unique"`
	Username  string            `json:"username"`
	Password  string            `'json:"-"` //表示在传递过程中忽略这个字段
	Email     string            `json:"email"`
	Openid    string            `json:"openid"`
	Avatar    string            `json:"avatar" gorm:"size:255"`                        // 头像：邮箱注册的头像或 QQ 登录的空间头像
	Address   string            `json:"address"`                                       // 地址
	Signature string            `json:"signature" gorm:"default:'签名是空白的，这位用户似乎比较低调哦。"` //用户签名
	RoleID    appTypes.RoleID   `json:"role_id"`                                       // 角色 ID
	Register  appTypes.Register `json:"register"`                                      // 注册来源
	Freeze    bool              `json:"freeze"`                                        // 用户是否被冻结
}

/*
大写开头 的标识符（如 User）是 公开的（Exported），可以被其他包访问。
小写开头 的标识符（如 user）是 私有的（unexported），仅限当前包内使用。
*/

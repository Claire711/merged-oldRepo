package database

import "server/global"

type Login struct {
	global.MODEL
	UserID      uint   `json:"user_id"`
	User        User   `json:"user" gorm:"foreignKey:UserID"` //显式指定外键，告诉 GORM 该关联通过 UserID 字段连接到 User 表的主键。
	LoginMethod string `json:"login_method"`                  // 登录方式
	IP          string `json:"ip"`                            // IP 地址
	Address     string `json:"address"`                       // 登录地址
	OS          string `json:"os"`                            // 操作系统
	DeviceInfo  string `json:"device_info"`                   // 设备信息
	BrowserInfo string `json:"browser_info"`                  // 浏览器信息
	Status      int    `json:"status"`                        // 登录状态
}

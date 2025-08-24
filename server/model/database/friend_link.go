package database

import "server/global"

type FriendLink struct {
	global.MODEL
	Logo        string `json:"logo" gorm:"size:255"`
	Image       Image  `json:"-" gorm:"foreignKey:Logo;reference:URL"`
	Link        string `json:"link"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

/*FOREIGN KEY (本表字段) REFERENCES 外表名 (外表字段名)
外表字段必须是主键！在关系型数据库中，外键确实传统上要求引用目标表的主键。
但 GORM 提供了更灵活的方式，允许使用唯一字段作为外键引用目标。*/

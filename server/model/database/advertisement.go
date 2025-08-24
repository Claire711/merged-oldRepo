package database

import "server/global"

type Advertisement struct {
	global.MODEL
	AdImage string `json:"ad_image" gorm:"size:255"`
	Image   Image  `json:"-" gorm:"foreignKey:AdImage;reference:URL"`
	Link    string `json:"link"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

/*json:"-"	字段在 JSON 中消失，不参与序列化和反序列化。
json:"_"	字段正常序列化，JSON 键名使用结构体字段名（而非自定义名称）。*/

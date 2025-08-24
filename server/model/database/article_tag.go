package database

type ArticleTag struct {
	Tag    string `json:"tag" gorm:"primaryKey"` //标签
	Number int    `json:"number"`                //统计数量
}

/*GORM 标签只在需要覆盖默认行为或添加特殊约束时才使用。没有标签的字段会按照 GORM 的约定自动映射，保持代码简洁性。*/

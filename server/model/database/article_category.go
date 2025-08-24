package database

type ArticleCategory struct {
	Category string `json:"category" gorm:"primaryKey"`
	Number   int    `json:"number"`
}

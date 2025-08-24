package database

type FootLink struct {
	Title string `json:"title" gorm:"primaryKey"`
	Link  string `json:"link"`
}

package global

import (
	"gorm.io/gorm"
	"time"
)

type MODEL struct {
	//Go 语言的结构体标签必须使用反引号
	ID uint `json:"id" gorm:"primaryKey"`
	//"primaryKey" K要大写，不然会被忽略
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` //给 deleted_at 加索引，加速软删除查询
	//增删改操作
}

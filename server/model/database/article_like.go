package database

import "server/global"

// ArticleLike 文章收藏表
type ArticleLike struct {
	//当嵌入 global.MODEL 后，ArticleLike 自动获得 global.MODEL所包含的所有字段
	global.MODEL
	ArticleID string `json:"article_id"` // 文章 ID
	UserID    uint   `json:"user_id"`    // 用户 ID
	User      User   `json:"-" gorm:"foreignKey:UserID"`
	/*
		UserID 是数据库中实际存储外键的字段（对应表中的 user_id 列）
		User 是 Go 语言中用于操作关联对象的字段（方便直接获取用户信息）
		UserID 字段：
		对应数据库中的 user_id 列，存储关联用户的 ID（如 1, 2）
		User 字段：
		告诉 GORM：“当查询 ArticleLike 时，根据 UserID 的值去 User 表查数据，填充到这个字段”
	*/
}

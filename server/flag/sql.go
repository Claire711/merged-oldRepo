package flag

import (
	"server/global"
	"server/model/database"
)

func SQL() error {
	//GORM 的自动迁移功能
	return global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		// 传入所有需要迁移的模型
		&database.Advertisement{},
		&database.ArticleCategory{},
		&database.ArticleLike{},
		&database.ArticleTag{},
		&database.Comment{},
		&database.Feedback{},
		&database.FootLink{},
		&database.FriendLink{},
		&database.Image{},
		&database.JwtBlacklist{},
		&database.Login{},
		&database.User{},
	)
}

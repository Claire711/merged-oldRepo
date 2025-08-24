package database

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/gofrs/uuid"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"server/global"
	"server/model/elasticsearch"
)

type Comment struct {
	global.MODEL
	ArticleID string    `json:"article_id"`
	PID       *uint     `json:"p_id"`                           // 指向父评论的ID
	PComment  *Comment  `json:"-" gorm:"foreignKey:PID"`        // 父评论对象（不显示在JSON）
	Children  []Comment `json:"children" gorm:"foreignKey:PID"` // 子评论列表
	UserUUID  uuid.UUID `json:"user_uuid" gorm:"type:char(36)"`
	User      User      `json:"user" gorm:"foreignKey:UserUUID;references:UUID"`
	Content   string    `json:"content"` //评论的内容
}

func (comment *Comment) AfterCreate(_ *gorm.DB) error {
	source := "ctx._source.comments += 1"
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), comment.ArticleID).Script(&script).Do(context.TODO())
	return err
}

func (comment *Comment) BeforeDelete(_ *gorm.DB) error {
	var articleID string
	if err := global.DB.Model(&comment).Pluck("article_id", &articleID).Error; err != nil {
		return err
	}
	source := "ctx._source.comments -= 1"
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), articleID).Script(&script).Do(context.TODO())
	return err
}

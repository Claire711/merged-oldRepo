package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
)

type CommentService struct {
}

func (commentService *CommentService) CommentInfoByArticleID(req request.CommentInfoByArticleID) ([]database.Comment, error) {
	var comments []database.Comment
	if err := global.DB.Where("article_id = ? AND p_id IS NULL", req.ArticleID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature")
	}).Find(&comments).Error; err != nil {
		return nil, err
	}
	for i := range comments {
		if err := commentService.LoadChildren(&comments[i]); err != nil {
			return nil, err
		}
	}
	return comments, nil
}

func (commentService *CommentService) CommentNew() ([]database.Comment, error) {
	var comments []database.Comment
	if err := global.DB.Order("id desc").Limit(5).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature")
	}).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (commentService *CommentService) CommentCreate(req request.CommentCreate) error {
	return global.DB.Create(&database.Comment{
		ArticleID: req.ArticleID,
		PID:       req.PID,
		UserUUID:  req.UserUUID,
		Content:   req.Content,
	}).Error
}

func (commentService *CommentService) CommentDelete(c *gin.Context, req request.CommentDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			var comment database.Comment
			if err := tx.Take(&comment, id).Error; err != nil {
				return err
			}
			userUUID := utils.GetUUID(c)
			roleID := utils.GetRoleID(c)
			if comment.UserUUID != userUUID && roleID != appTypes.Admin {
				return errors.New("no authorization to delete this comment")
			}
			if err := commentService.DeleteCommentAndChildren(tx, comment.ID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (commentService *CommentService) CommentInfo(uuid uuid.UUID) ([]database.Comment, error) {
	var rawComments []database.Comment
	if err := global.DB.Order("id desc").Where("user_uuid = ?", uuid).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature")
	}).Find(&rawComments).Error; err != nil {
		return nil, err
	}
	for i := range rawComments {
		if err := commentService.LoadChildren(&rawComments[i]); err != nil {
			return nil, err
		}
	}
	// 评论去重，过滤评论列表，排除掉那些需要被移除的子评论(如果当前评论的子评论存在为你的评论，就去除子评论)
	var destComment []database.Comment
	ansMap := commentService.FindChildCommentsIDByRootCommentUserUUID(rawComments)
	for i := range rawComments {
		if _, exists := ansMap[rawComments[i].ID]; !exists {
			destComment = append(destComment, rawComments[i])
		}
	}
	return destComment, nil
}

func (commentService *CommentService) CommentList(info request.CommentList) (interface{}, int64, error) {
	db := global.DB
	if info.ArticleID != nil {
		db = db.Where("article_id =?", *info.ArticleID)
	}
	if info.UserUUID != nil {
		db = db.Where("user_uuid =?", *info.UserUUID)
	}
	if info.Content != nil {
		db = db.Where("content LIKE ?", "%"+*info.Content+"%")
	}
	option := other.MySQLOption{
		Where:    db,
		PageInfo: info.PageInfo,
	}
	return utils.MySQLPagination(&database.Comment{}, option)
}

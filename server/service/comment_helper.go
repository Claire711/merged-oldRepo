package service

import (
	"gorm.io/gorm"
	"server/global"
	"server/model/database"
)

func (commentService *CommentService) LoadChildren(comment *database.Comment) error {
	var children []database.Comment

	// 查询直接子评论（预加载用户基本信息）
	if err := global.DB.Where("p_id = ?", comment.ID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("uuid, username, avatar, address, signature")
		}).
		Find(&children).Error; err != nil {
		return err
	}

	// 递归加载子评论的子评论
	for i := range children {
		if err := commentService.LoadChildren(&children[i]); err != nil {
			return err
		}
	}

	// 将加载的子评论附加到当前评论
	comment.Children = children
	return nil
}

func (commentService *CommentService) DeleteCommentAndChildren(tx *gorm.DB, commentID uint) error {
	var children []database.Comment
	if err := tx.Where("p_id = ?", commentID).Find(&children).Error; err != nil {
		return err
	}
	for _, childComment := range children {
		if err := commentService.DeleteCommentAndChildren(tx, childComment.ID); err != nil {
			return err
		}
	}
	if err := tx.Delete(&database.Comment{}, commentID).Error; err != nil {
		return err
	}
	return nil
}

func (commentService *CommentService) FindChildCommentsIDByRootCommentUserUUID(comments []database.Comment) map[uint]struct{} {
	//递归查找所有与 “根评论作者” 相同的子评论 ID
	result := make(map[uint]struct{})
	for _, rootComment := range comments {
		var findChildren func([]database.Comment)
		findChildren = func(children []database.Comment) {
			for _, child := range children {
				if child.UserUUID == rootComment.UserUUID {
					result[child.ID] = struct{}{}
				}
				//Go 语言的内置函数 len() 可以返回切片、数组、字符串、映射（map）等类型的 “元素数量”
				if len(child.Children) > 0 {
					findChildren(child.Children)
				}
			}
		}
		findChildren(rootComment.Children)
	}
	return result
}

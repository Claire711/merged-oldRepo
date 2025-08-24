package service

import (
	"context"
	"errors"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"gorm.io/gorm"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/elasticsearch"
	"server/model/other"
	"server/utils"
	"strconv"
	"time"

	"server/model/request"
)

type ArticleService struct {
}

func (articleService *ArticleService) ArticleInfoByID(id string) (elasticsearch.Article, error) {
	// 异步更新浏览量
	go func() {
		articleView := articleService.NewArticleView()
		_ = articleView.Set(id)
	}()
	return articleService.Get(id)
}

func (articleService *ArticleService) ArticleSearch(info request.ArticleSearch) (interface{}, int64, error) {
	req := &search.Request{
		Query: &types.Query{},
	}
	boolQuery := &types.BoolQuery{}
	hasConditions := false
	if info.Query != "" {
		hasConditions = true
		boolQuery.Should = []types.Query{
			{Match: map[string]types.MatchQuery{"title": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"keyword": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"abstract": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"content": {Query: info.Query}}},
		}
	}
	if info.Tag != "" {
		hasConditions = true
		boolQuery.Must = []types.Query{
			{Match: map[string]types.MatchQuery{"tags": {Query: info.Tag}}},
		}
	}
	if info.Category != "" {
		hasConditions = true
		boolQuery.Filter = []types.Query{
			{Term: map[string]types.TermQuery{"category": {Value: info.Category}}},
		}
	}
	if hasConditions {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{} // 无任何条件时，匹配所有文章
	}
	if info.Sort != "" {
		var sortField string
		switch info.Sort {
		case "time":
			sortField = "created_at" // 按创建时间排序
		case "view":
			sortField = "views" // 按浏览量排序
		case "comment":
			sortField = "comments" // 按评论数排序
		case "like":
			sortField = "likes" // 按点赞数排序
		default:
			sortField = "created_at" // 默认按创建时间
		}
		var order sortorder.SortOrder
		if info.Order != "asc" {
			order = sortorder.Desc
		} else {
			order = sortorder.Asc
		}
		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					sortField: {Order: &order},
				},
			},
		}
	}
	option := other.EsOption{
		PageInfo: info.PageInfo,
		Index:    elasticsearch.ArticleIndex(),
		Request:  req,
		// 需要返回的字段（过滤掉不需要的字段，减少传输）
		SourceIncludes: []string{"created_at", "cover", "title", "abstract", "category", "tags", "views", "comments", "likes"},
	}
	return utils.EsPagination(context.TODO(), option)
}

func (articleService *ArticleService) ArticleCategory() ([]database.ArticleCategory, error) {
	var category []database.ArticleCategory
	err := global.DB.Find(&category).Error
	//查询表中的所有记录，并将结果存入 category 切片中
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (articleService *ArticleService) ArticleTags() ([]database.ArticleTag, error) {
	var tags []database.ArticleTag
	err := global.DB.Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (articleService *ArticleService) ArticleLike(req request.ArticleLike) error {
	//执行一个事务
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var likeRecord database.ArticleLike
		var num int
		//判断用户是否收藏文章
		if errors.Is(tx.Where("article_id = ? AND user_id = ?", req.ArticleID, req.UserID).First(&likeRecord).Error, gorm.ErrRecordNotFound) {
			num = 1
			if err := tx.Create(&database.ArticleLike{UserID: req.UserID, ArticleID: req.ArticleID}).Error; err != nil {
				return err
			}

		} else { // 如果用户已经收藏，则取消收藏
			if err := tx.Delete(&likeRecord).Error; err != nil {
				return err
			}
			num = -1
		}
		source := "ctx._source.likes += " + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
		_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), req.ArticleID).Script(&script).Do(context.TODO())
		return err
	})
}

func (articleService *ArticleService) ArticleIsLike(req request.ArticleLike) (bool, error) {
	return !errors.Is(global.DB.Where("article_id = ? AND user_id = ?", req.ArticleID, req.UserID).First(&database.ArticleLike{}).Error, gorm.ErrRecordNotFound), nil
}

func (articleService *ArticleService) ArticleLikesList(info request.ArticleLikesList) (interface{}, int64, error) {
	db := global.DB.Where("user_id = ?", info.UserID)
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}
	//likes存储从数据库查询到的 “用户点赞记录”
	likes, total, err := utils.MySQLPagination(&database.ArticleLike{}, option)
	if err != nil {
		return nil, 0, err
	}
	var list []struct {
		Id_     string                `json:"_id"`
		Source_ elasticsearch.Article `json:"_source"`
	}
	//articleLike :likes 切片中的单个元素，表示 “用户对某一篇文章的点赞记录”
	for _, articleLike := range likes {
		article, err := articleService.Get(articleLike.ArticleID)
		if err != nil {
			return nil, 0, err
		}
		article.UpdatedAt = ""
		article.Keyword = ""
		article.Content = ""
		list = append(list, struct {
			Id_     string                `json:"_id"`
			Source_ elasticsearch.Article `json:"_source"`
		}{
			Id_:     articleLike.ArticleID,
			Source_: article,
		})
	}
	return list, total, nil
}

func (articleService *ArticleService) ArticleCreate(req request.ArticleCreate) error {
	exists, err := articleService.Exits(req.Title)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("the article already exists")
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToCreate := elasticsearch.Article{
		CreatedAt: now,
		UpdatedAt: now,
		Cover:     req.Cover,
		Title:     req.Title,
		Keyword:   req.Title,
		Category:  req.Category,
		Tags:      req.Tags,
		Abstract:  req.Abstract,
		Content:   req.Content,
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := articleService.UpdateCategoryCount(tx, "", articleToCreate.Category); err != nil {
			return err
		}
		if err := articleService.UpdateTagsCount(tx, []string{}, articleToCreate.Tags); err != nil {
			return err
		}
		//修改封面图片类型
		if err := utils.ChangeImagesCategory(tx, []string{articleToCreate.Cover}, appTypes.Cover); err != nil {
			return err
		}
		illustrations, err := utils.FindIllustrations(articleToCreate.Content)
		if err != nil {
			return err
		}
		err = utils.ChangeImagesCategory(tx, illustrations, appTypes.Illustration)
		if err != nil {
			return err
		}
		return articleService.Create(&articleToCreate)
	})
}

func (articleService *ArticleService) ArticleDelete(req request.ArticleDelete) error {
	//先判断ids是否为空
	if len(req.IDs) == 0 {
		return nil
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		ids := req.IDs
		for _, id := range ids {
			articleToDelete, err := articleService.Get(id)
			if err != nil {
				return err
			}
			if err := articleService.UpdateCategoryCount(tx, articleToDelete.Category, ""); err != nil {
				return err
			}
			if err := articleService.UpdateTagsCount(tx, articleToDelete.Tags, []string{}); err != nil {
				return err
			}
			if err := utils.InitImagesCategory(tx, []string{articleToDelete.Cover}); err != nil {
				return err
			}
			illustrations, err := utils.FindIllustrations(articleToDelete.Content)
			if err != nil {
				return err
			}
			if err := utils.InitImagesCategory(tx, illustrations); err != nil {
				return err
			}
			//删除评论
			comments, err := ServiceGroupApp.CommentService.CommentInfoByArticleID(request.CommentInfoByArticleID{ArticleID: id})
			if err != nil {
				return err
			}
			for _, comment := range comments {
				if err := ServiceGroupApp.CommentService.DeleteCommentAndChildren(tx, comment.ID); err != nil {
					return err
				}
			}

		}
		return articleService.Delete(ids)
	})
}

func (articleService *ArticleService) ArticleUpdate(req request.ArticleUpdate) error {
	if req.ID == "" {
		return errors.New("文章ID不能为空")
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToUpdate := struct {
		UpdatedAt string   `json:"updated_at"`
		Cover     string   `json:"cover"`
		Title     string   `json:"title"`
		Keyword   string   `json:"keyword"`
		Category  string   `json:"category"`
		Tags      []string `json:"tags"`
		Abstract  string   `json:"abstract"`
		Content   string   `json:"content"`
	}{

		UpdatedAt: now,
		Cover:     req.Cover,
		Title:     req.Title,
		Keyword:   req.Title,
		Category:  req.Category,
		Tags:      req.Tags,
		Abstract:  req.Abstract,
		Content:   req.Content,
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		oldArticle, err := articleService.Get(req.ID)
		if err != nil {
			return err
		}
		if err := articleService.UpdateCategoryCount(tx, oldArticle.Category, articleToUpdate.Category); err != nil {
			return err
		}
		if err := articleService.UpdateTagsCount(tx, oldArticle.Tags, articleToUpdate.Tags); err != nil {
			return err
		}
		if oldArticle.Cover != articleToUpdate.Cover {
			if err := utils.InitImagesCategory(tx, []string{oldArticle.Cover}); err != nil {
				return err
			}
			if err := utils.ChangeImagesCategory(tx, []string{articleToUpdate.Cover}, appTypes.Cover); err != nil {
				return err
			}
		}
		oldIllustrations, err := utils.FindIllustrations(oldArticle.Content)
		if err != nil {
			return err
		}
		newIllustrations, err := utils.FindIllustrations(articleToUpdate.Content)
		if err != nil {
			return err
		}
		addedIllustrations, removedIllustrations := utils.DiffArrays(oldIllustrations, newIllustrations)
		if err := utils.InitImagesCategory(tx, removedIllustrations); err != nil {
			return err
		}
		if err := utils.ChangeImagesCategory(tx, addedIllustrations, appTypes.Illustration); err != nil {
			return err
		}

		return articleService.Update(req.ID, articleToUpdate)
	})
}

func (articleService *ArticleService) ArticleList(info request.ArticleList) (interface{}, int64, error) {
	req := &search.Request{
		Query: &types.Query{},
	}
	boolQuery := &types.BoolQuery{}
	// 根据标题查询
	if info.Title != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{
			Match: map[string]types.MatchQuery{
				"title": {Query: *info.Title},
			},
		})
	}
	// 根据简介查询
	if info.Abstract != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{
			Match: map[string]types.MatchQuery{
				"abstract": {Query: *info.Abstract},
			},
		})
	}

	// 根据类别筛选
	if info.Category != nil {
		boolQuery.Filter = []types.Query{
			{
				Term: map[string]types.TermQuery{
					"category": {Value: info.Category},
				},
			},
		}
	}

	if boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
		req.Sort = []types.SortCombinations{ // 步骤1：排序规则的外层容器
			types.SortOptions{ // 步骤2：单个排序条件的包装
				SortOptions: map[string]types.FieldSort{ // 步骤3：字段与排序规则的映射
					"created_at": {Order: &sortorder.Desc}, // 步骤4：具体字段的排序规则
				},
			},
		}
	}

	options := other.EsOption{
		PageInfo: info.PageInfo,
		Index:    elasticsearch.ArticleIndex(),
		Request:  req,
	}
	return utils.EsPagination(context.TODO(), options)
}

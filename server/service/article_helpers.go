package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"gorm.io/gorm"
	"server/utils"

	"server/global"
	"server/model/database"
	"server/model/elasticsearch"
)

func (articleService *ArticleService) Create(a *elasticsearch.Article) error {
	// 将文章索引到Elasticsearch中，并设置刷新操作为 true
	_, err := global.ESClient.Index(elasticsearch.ArticleIndex()).Request(a).Refresh(refresh.True).Do(context.TODO())
	return err
}

func (articleService *ArticleService) Delete(ids []string) error {
	var req bulk.Request
	//循环构建删除任务
	for _, id := range ids {
		req = append(req, types.OperationContainer{Delete: &types.DeleteOperation{Id_: &id}})
	}
	//执行批量删除
	_, err := global.ESClient.Bulk().Request(&req).Index(elasticsearch.ArticleIndex()).Refresh(refresh.True).Do(context.TODO())
	return err
}

func (articleService *ArticleService) Get(id string) (elasticsearch.Article, error) {
	var a elasticsearch.Article
	// 从Elasticsearch获取文章
	res, err := global.ESClient.Get(elasticsearch.ArticleIndex(), id).Do(context.TODO())
	if err != nil {
		return elasticsearch.Article{}, err
	}
	// 如果找不到该文档，则返回错误
	if !res.Found {
		return elasticsearch.Article{}, errors.New("document not found")
	}
	// 将返回的源数据反序列化为 Article 对象
	err = json.Unmarshal(res.Source_, &a)
	return a, err
}

func (articleService *ArticleService) Update(id string, v any) error {
	//传入文章的 id（要更新哪篇）和 v（更新的内容）
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = global.ESClient.Update(elasticsearch.ArticleIndex(), id).Request(&update.Request{Doc: bytes}).Refresh(refresh.True).Do(context.TODO())
	return err
}
func (articleService *ArticleService) Exits(title string) (bool, error) {
	req := &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{"keyword": {Query: title}},
		},
	}
	res, err := global.ESClient.Search().Index(elasticsearch.ArticleIndex()).Request(req).Size(1).Do(context.TODO())
	if err != nil {
		return false, err
	}
	return res.Hits.Total.Value > 0, nil
}

// UpdateCategoryCount 更新文章类别的计数（增加或减少）
func (articleService *ArticleService) UpdateCategoryCount(tx *gorm.DB, oldCategory, newCategory string) error {
	if newCategory == oldCategory {
		return nil
	}
	if newCategory != "" {
		// 如果新类别不为空，更新新类别的文章计数
		var newArticleCategory database.ArticleCategory
		//从查询结果中取第一条记录，存到 newArticleCategory 这个变量里
		if errors.Is(tx.Where("category = ?", newCategory).First(&newArticleCategory).Error, gorm.ErrRecordNotFound) {
			if err := tx.Create(&database.ArticleCategory{Category: newCategory, Number: 1}).Error; err != nil {
				return err
			}
		} else {
			// 如果类别已存在，更新该类别的计数
			if err := tx.Model(&newArticleCategory).Update("number", gorm.Expr("number + ?", 1)).Error; err != nil {
				return err
			}
		}
	}
	// 如果旧类别不为空，更新旧类别的文章计数
	if oldCategory != "" {
		var oldArticleCategory database.ArticleCategory
		if err := tx.Where("category = ?", oldCategory).First(&oldArticleCategory).Update("number", gorm.Expr("number - ?", 1)).Error; err != nil {
			return err
		}
		if oldArticleCategory.Number == 1 {
			if err := tx.Delete(&oldArticleCategory).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateTagsCount 更新文章标签的计数（增加或减少）
func (articleService *ArticleService) UpdateTagsCount(tx *gorm.DB, oldTags, newTags []string) error {
	addedTags, removedTags := utils.DiffArrays(oldTags, newTags)
	for _, addedTag := range addedTags {
		var t database.ArticleTag
		if errors.Is(tx.Where("tag = ?", addedTag).First(&t).Error, gorm.ErrRecordNotFound) {
			if err := tx.Create(&database.ArticleTag{Tag: addedTag, Number: 1}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&t).Update("number", gorm.Expr("number + ?", 1)).Error; err != nil {
				return err
			}
		}
	}
	for _, removedTag := range removedTags {
		var t database.ArticleTag
		if err := tx.Where("tag = ?", removedTag).First(&t).Error; err != nil {
			return err
		}
		// 执行更新
		if err := tx.Model(&t).Update("number", gorm.Expr("number - ?", 1)).Error; err != nil {
			return err
		}
		//更新的是数据库的原始值，更新值没有传回来，只在数据库层面完成操作
		if t.Number == 1 {
			if err := tx.Delete(&t).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

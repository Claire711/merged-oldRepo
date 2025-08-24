package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// Article 文章表
type Article struct {
	CreatedAt string   `json:"created_at"` // 创建时间
	UpdatedAt string   `json:"updated_at"` // 更新时间
	Cover     string   `json:"cover"`      // 文章封面
	Title     string   `json:"title"`      // 文章标题
	Keyword   string   `json:"keyword"`    // 文章标题-关键字
	Category  string   `json:"category"`   // 文章类别
	Tags      []string `json:"tags"`       // 文章标签
	Abstract  string   `json:"abstract"`   // 文章简介
	Content   string   `json:"content"`    // 文章内容
	Views     int      `json:"views"`      // 浏览量
	Comments  int      `json:"comments"`   // 评论量
	Likes     int      `json:"likes"`      // 收藏量
}

// ArticleIndex 文章 ES 索引
func ArticleIndex() string { return "article_index" }

// Mapping映射
func ArticleMapping() *types.TypeMapping {
	return &types.TypeMapping{
		Properties: map[string]types.Property{
			"created_at": types.DateProperty{NullValue: nil, Format: func(s string) *string { return &s }("yyyy-MM-dd HH:mm:ss")},
			"updated_at": types.DateProperty{NullValue: nil, Format: func(s string) *string { return &s }("yyyy-MM-dd HH:mm:ss")},
			"cover":      types.TextProperty{},
			"title":      types.TextProperty{},
			"keyword":    types.KeywordProperty{},
			"category":   types.KeywordProperty{},
			"tags":       []types.KeywordProperty{},
			"abstract":   types.TextProperty{},
			"content":    types.TextProperty{},
			"views":      types.IntegerNumberProperty{},
			"comments":   types.IntegerNumberProperty{},
			"likes":      types.IntegerNumberProperty{},
		},
	}
}

/*
Elasticsearch 客户端要求 Format 是一个 指针类型（*string）
而直接写 "yyyy-MM-dd HH:mm:ss" 是字符串值，因此需要用匿名函数将字符串转为指针：
func(s string) *string { return &s }("yyyy-MM-dd HH:mm:ss")
等价于手动创建指针：
format := "yyyy-MM-dd HH:mm:ss"
Format: &format
*/

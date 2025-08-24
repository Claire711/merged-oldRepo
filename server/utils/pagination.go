package utils

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"server/global"
	"server/model/other"
)

/*
1.[T any]：这是 Go 1.18+ 引入的泛型语法，表示该函数是一个泛型函数
1.1核心是通过类型参数 T 实现 “一次编写，多类型复用”
2.Where *gorm.DB 中的 * 主要有两个作用：
类型匹配：与 GORM 框架返回的 *gorm.DB 指针类型保持一致，确保能接收链式调用产生的查询条件。
支持空值：通过 nil 表示 “无查询条件”，使结构体更灵活地适应 “有条件” 和 “无条件” 两种查询场景。
*/

func MySQLPagination[T any](model *T, option other.MySQLOption) (list []T, total int64, err error) {
	if option.Page < 1 {
		option.Page = 1 // 页码不能小于1，默认为1
	}
	if option.PageSize < 1 {
		option.PageSize = 10 // 每页记录数不能小于1，默认为10(仅修正无效值，保留合理值)
	}
	if option.Order == "" {
		option.Order = "id desc" // 默认按id降序排列
	}
	// 创建查询
	query := global.DB.Model(model)
	if option.Where != nil {
		query = query.Where(query)
	}
	// 计算符合条件的记录总数
	if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	for _, preload := range option.Preload {
		query = query.Preload(preload)
	}

	// 应用分页查询
	err = query.Order(option.Order).
		Limit(option.PageSize).                      // 设置每页记录数
		Offset((option.Page - 1) * option.PageSize). // 设置偏移量，根据页码计算
		Find(&list).Error                            // 执行查询，并将结果存入 list 中
	return list, total, err
}

func EsPagination(ctx context.Context, option other.EsOption) (list []types.Hit, total int64, err error) {
	if option.Page < 1 {
		option.Page = 1 // 页码不能小于1，默认为1
	}
	if option.PageSize < 1 {
		option.PageSize = 10 // 每页记录数不能小于1，默认为10(仅修正无效值，保留合理值)
	}
	from := (option.Page - 1) * option.PageSize
	option.Request.Size = &option.PageSize
	option.Request.From = &from
	res, err := global.ESClient.Search().
		Index(option.Index).
		Request(option.Request).
		SourceIncludes_(option.SourceIncludes...).
		Do(ctx)
	if err != nil {
		return nil, 0, err
	}
	list = res.Hits.Hits
	total = res.Hits.Total.Value
	// 执行查询，并将结果存入 list 中
	return list, total, err
}

package service

import (
	"server/global"
	"strconv"
)

// NewArticleView 初始化计数器
func (articleService *ArticleService) NewArticleView() CountDB {
	return CountDB{
		Index: "article_views",
	}
}

type CountDB struct {
	Index string
}

// Set 更新计数（自增)
func (c CountDB) Set(id string) error {
	//从 Redis 的哈希表（Hash）中，根据c.Index（哈希表的键名）和id（哈希表中的字段名）获取对应的值，然后转换为整数num,忽略err
	num, _ := global.Redis.HGet(c.Index, id).Int()
	num++
	//向 Redis 的哈希表（Hash）中，以 c.Index 为哈希表键名（key），id 为字段名（field），设置值为 num。
	err := global.Redis.HSet(c.Index, id, num).Err()
	return err
}

// GetInfo 获取所有统计信息
// 应用场景：用于后台统计页面展示所有文章的浏览量数据
func (c CountDB) GetInfo() map[string]int {
	//map[string]int {} 用{}初始化
	var info = map[string]int{}
	maps := global.Redis.HGetAll(c.Index).Val()
	for id, val := range maps {
		num, _ := strconv.Atoi(val)
		info[id] = num
	}
	return info
}

// Clear 清空统计数据
// 应用场景：需要重置所有浏览量统计时使用（如定期数据清零）
func (c CountDB) Clear() error {
	return global.Redis.Del(c.Index).Err()
}

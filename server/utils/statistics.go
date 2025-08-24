package utils

import "gorm.io/gorm"

func FetchDateCounts(db *gorm.DB, query *gorm.DB) map[string]int {
	var dateCounts []struct {
		Date  string `json:"data"`
		Count int    `json:"count"`
	}
	db.Where(query).
		Select("data_format(created_date,'%Y-%m-%d') as date", "count(id) as count").
		Group("date").Scan(&dateCounts)
	//将结构体切片dateCounts中的数据转换为以日期为键、数量为值的映射（map）
	dateCountMap := make(map[string]int)
	//_：忽略循环的索引,count：每次循环获取的切片元素（是一个结构体实例）
	for _, count := range dateCounts {
		//dateCountMap[map 的键] = map 的值
		dateCountMap[count.Date] = count.Count
	}
	return dateCountMap
}

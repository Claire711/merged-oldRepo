package service

import (
	"encoding/json"
	"server/global"
	"server/model/other"
	"server/utils/hotSearch"
	"time"
)

type HotSearchService struct {
}

func (hotSearchService *HotSearchService) GetHotSearchDataBySource(sourceStr string) (other.HotSearchData, error) {
	result, err := global.Redis.Get(sourceStr).Result()
	if err != nil {
		source := hotSearch.NewSource(sourceStr)
		hotSearchData, err := source.GetHotSearchData(30)
		if err != nil {
			return other.HotSearchData{}, err
		}
		// 将数据序列化为JSON字节流（便于存储到Redis）
		//json.Marshal(hotSearchData) 将结构体数据转换为 JSON 字节流（因为 Redis 只能存储字符串 / 字节流，无法直接存储结构体）
		bytes, err := json.Marshal(hotSearchData)
		if err != nil {
			return other.HotSearchData{}, err
		}
		// 将JSON数据存入Redis，设置过期时间为1小时
		if err := global.Redis.Set(sourceStr, bytes, time.Hour).Err(); err != nil {
			return other.HotSearchData{}, err
		}
		return hotSearchData, nil
	}
	var hotSearchData other.HotSearchData
	if err := json.Unmarshal([]byte(result), &hotSearchData); err != nil {
		return other.HotSearchData{}, err
	}
	return hotSearchData, nil
}

package hotSearch

import (
	"errors"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"regexp"
	"server/model/other"
	"strconv"
	"time"
)

type Baidu struct {
}

func (*Baidu) GetHotSearchData(maxNum int) (HotSearchData other.HotSearchData, err error) {
	// maxNum int:“最多获取的热门搜索数据数量”（即返回结果的最大条目数）
	resp, err := http.Get("https://top.baidu.com/board?tab=realtime")
	if err != nil {
		return other.HotSearchData{}, err
	}
	defer resp.Body.Close() //关闭响应体
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return other.HotSearchData{}, err
	}
	var jsonStr string
	reg := regexp.MustCompile(`<!--s-data:({.*?})-->`)
	result := reg.FindAllStringSubmatch(string(body), -1)
	if len(result) > 0 && len(result[0]) > 1 {
		jsonStr = result[0][1] // 将第一个匹配项的第一个捕获组内容赋值给 jsonStr
	} else {
		return other.HotSearchData{}, errors.New("failed to get data")
	}
	updateTime := time.Unix(gjson.Get(jsonStr, "data.cards.0.updateTime").Int(), 0).Format("2006-01-02 15:04:05")
	var hotList []other.HotItem
	for i := 0; i < maxNum; i++ {
		if index := gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".index"); !index.Exists() {
			break
		}
		hotList = append(hotList, other.HotItem{
			Index:       i + 1,
			Title:       gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".word").Str,
			Description: gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".desc").Str,
			Image:       gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".img").Str,
			Popularity:  gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".hotScore").Str,
			URL:         gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".rawUrl").Str,
		})
	}
	return other.HotSearchData{Source: "百度热搜", UpdateTime: updateTime, HotList: hotList}, nil
}

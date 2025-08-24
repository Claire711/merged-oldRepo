package task

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"server/global"
	"server/model/elasticsearch"
	"server/service"
	"strconv"
)

// UpdateArticleViewsSyncTask 将临时存储的文章浏览量数据（可能来自缓存或临时表）同步更新到 Elasticsearch 的文章索引中，然后清空临时存储的浏览量数据。
func UpdateArticleViewsSyncTask() error {
	articleView := service.ServiceGroupApp.ArticleService.NewArticleView()
	viewsInfo := articleView.GetInfo()
	for id, num := range viewsInfo {
		if num == 0 {
			continue
		}
		//// 更新数据:之前的数据+缓存中的数据
		source := "ctx._source.views +=" + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
		_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), id).Script(&script).Do(context.TODO())
		return err
	}
	articleView.Clear()
	return nil
}

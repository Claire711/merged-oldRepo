package other

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"gorm.io/gorm"
	"server/model/request"
)

type MySQLOption struct {
	request.PageInfo
	Order   string   //排序规则
	Where   *gorm.DB //查询条件
	Preload []string //预加载关联字段
}
type EsOption struct {
	request.PageInfo
	Index          string
	Request        *search.Request
	SourceIncludes []string
}

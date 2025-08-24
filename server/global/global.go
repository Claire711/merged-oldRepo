package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"server/config"
)

var (
	/* Config 加 *：表示声明的是指针变量，存储的是另一个变量的内存地址（通常用于引用复杂类型、避免拷贝开销，或需要共享状态）。
	不加 *：表示声明的是值类型变量，直接存储数据本身（通常用于基础类型、结构体实例，或不需要共享的场景）。*/
	Config   *config.Config
	Log      *zap.Logger
	DB       *gorm.DB
	ESClient *elasticsearch.TypedClient
	Redis    redis.Client
	//接口变量本身是一个 "引用"，无需再加指针。例如 local_cache.Cache 是接口类型，声明时直接用值类型：

	BlackCache local_cache.Cache
)

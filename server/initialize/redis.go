package initialize

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"os"
	"server/global"
)

// ConnectRedis 初始化并返回一个 Redis 客户端，支持集群或单节点配置
func ConnectRedis() redis.Client {

	redisCfg := global.Config.Redis
	// 使用单节点配置创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Address,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	_, err := client.Ping().Result()
	if err != nil {
		global.Log.Error("Failed to connect Redis: ", zap.Error(err))
		os.Exit(1)
	}
	return *client
}

package main

import (
	"server/core"
	"server/flag"
	"server/global"
	"server/initialize"
)

func main() {
	global.Config = core.InitConf()
	global.Log = core.InitLogger()
	initialize.OtherInit()
	global.DB = initialize.InitGorm()
	global.Redis = initialize.ConnectRedis()
	global.ESClient = initialize.ConnectEs()
	defer global.Redis.Close()
	flag.InitFlag()
	initialize.InitCron()
	core.RunServer()
}

//架构设计是一个持续优化的过程，而非一次性完成的工作。

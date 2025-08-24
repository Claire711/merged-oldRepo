package initialize

import (
	"github.com/robfig/cron/v3" // Cron 定时任务库
	"go.uber.org/zap"           // Uber 高性能日志库
	"os"                        // 用于错误时退出程序
	"server/global"             // 全局变量（如日志实例）
	"server/task"               // 定时任务的具体实现
)

// 定时任务用来获取博客浏览量，每天获取一次日历，每小时刷新一次新闻获取
type ZapLogger struct {
	logger *zap.Logger
}

// 将 cron 库的日志输出适配到 zap。
func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, zap.Any("keysAndValues", keysAndValues))
}
func (l *ZapLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.Error(msg, zap.Error(err), zap.Any("keysAndValues", keysAndValues))
}
func NewZapLogger() *ZapLogger {
	return &ZapLogger{logger: global.Log}
}
func InitCron() {
	c := cron.New(cron.WithLogger(NewZapLogger()))
	err := task.RegisterScheduledTasks(c)
	if err != nil {
		global.Log.Error("Error scheduling cron job :", zap.Error(err))
		os.Exit(1)
	}
	c.Start()
}

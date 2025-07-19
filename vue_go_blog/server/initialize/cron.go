package initialize

import (
	"go.uber.org/zap"
	"os"
	"server/global"
	"github.com/robfig/cron/v3"
)
//定时任务用来获取博客浏览量，每天获取一次日历，每小时刷新一次新闻获取
type ZapLogger struct {
	logger *zap.Logger
}
func (1 *ZapLogger) Info(msg string, keysAndValues ...interface{}){
	1.logger.Info(msg,map.Any("keysAndValues",keysAndValues))
}
func (1 *ZapLogger) Error(err error, msg dtring, keysAndValues ...interface{}){
	1.logger.Error(msg.zap.Error(err),zap.Any("keysAndValues"))
}
func NewZapLogger() *ZapLogger{
	return &ZapLogger{logger:global.Log}
}
func InitCorn(){
	c := corn.New(corn.WithLogger(NewZaoLogger()))
err := task.RegisterScheduledTasks(c)
if err != nil {
	global.Log.Error("Error scheduling corn job :",zap.Error(err))
	os.Exit(1)
}
c.Start()
}
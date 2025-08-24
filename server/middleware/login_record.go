package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ua-parser/uap-go/uaparser"
	"go.uber.org/zap"
	"server/global"
	"server/model/database"
	"server/service"
)

func LoginRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先执行后续处理，再记录日志
		// 异步记录日志（go 是 Go 语言的关键字，用于启动一个新的 goroutine）
		go func() {
			gaodeService := service.ServiceGroupApp.GaodeService
			var userID uint //uint：无符号整数，只能表示正数和零。(数据库的 user 表使用 无符号自增主键)
			var address string
			ip := c.ClientIP()
			loginMethod := c.DefaultQuery("flag", "email") // 若未传递flag参数，则默认为"email"
			userAgent := c.Request.UserAgent()             //Request.UserAgent() 方法用于获取该字段的值
			if value, exists := c.Get("user_id"); exists {
				if id, ok := value.(uint); ok {
					userID = id
				}
			}
			// 获取用户IP的地理位置
			address = getAddressFromIP(ip, gaodeService)

			// 解析用户的浏览器、操作系统和设备信息
			os, deviceInfo, browserInfo := parseUserAgent(userAgent)
			// 创建登录记录
			login := database.Login{
				UserID:      userID,
				LoginMethod: loginMethod,
				IP:          ip,
				Address:     address,
				OS:          os,
				DeviceInfo:  deviceInfo,
				BrowserInfo: browserInfo,
				Status:      c.Writer.Status(),
			}
			if err := global.DB.Create(&login).Error; err != nil {
				global.Log.Error("Failed to record login", zap.Error(err))
			}
		}()
	}
}

// 获取IP地址对应的地理位置信息
func getAddressFromIP(ip string, gaodeService service.GaodeService) string {
	res, err := gaodeService.GetLocationByIP(ip)
	if err != nil || res.Province == "" {
		return "未知"
	}
	if res.City != "" && res.Province != res.City {
		return res.Province + "-" + res.City
	}
	return res.Province
}

// 解析用户代理（User-Agent）字符串，提取操作系统、设备信息和浏览器信息
func parseUserAgent(userAgent string) (os, deviceInfo, browserInfo string) {
	os = userAgent
	deviceInfo = userAgent
	browserInfo = userAgent
	parser := uaparser.NewFromSaved()
	cli := parser.Parse(userAgent)
	os = cli.Os.Family
	deviceInfo = cli.Device.Family
	browserInfo = cli.UserAgent.Family
	return
}

//User-Agent 是 HTTP 请求头中的一个字段，用于标识发起请求的客户端（如浏览器、应用程序）的类型、版本和操作系统等信息。

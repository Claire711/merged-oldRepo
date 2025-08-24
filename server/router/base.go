package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)

type BaseRouter struct {
}

// func (接收者) 函数名 （传入的参数）
func (b *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("base")
	baseApi := api.ApiGroupApp.BaseApi
	{
		baseRouter.POST("captcha", baseApi.Captcha)                                     //图形验证码
		baseRouter.POST("sendEmailVerificationCode", baseApi.SendEmailVerificationCode) //发送邮箱验证码
		baseRouter.GET("qqLoginURL", baseApi.QQLoginURL)
	}
}

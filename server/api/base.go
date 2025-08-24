package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"server/global"
	"server/model/request"
	"server/model/response"
)

type BaseApi struct {
}

// 创建验证码对象
var store = base64Captcha.DefaultMemStore

// Captcha 生成6位数字的验证码
func (baseApi *BaseApi) Captcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)
	captcha := base64Captcha.NewCaptcha(driver, store)
	//生成验证码
	id, base64, _, err := captcha.Generate()
	if err != nil {
		global.Log.Error("Failed to generate captcha", zap.Error(err))
		response.FailWithMessage("Failed to generate captcha", c)
		return
	}
	response.OkWithData(response.Captcha{
		CaptchaID: id,
		PicPath:   base64,
	}, c)
}

// SendEmailVerificationCode 发送邮箱验证码
func (baseApi *BaseApi) SendEmailVerificationCode(c *gin.Context) {
	//声明变量的类型声明语法
	var req request.SendEmailVerificationCode
	err := c.ShouldBindJSON(&req)
	if err != nil {
		//err.Error() 是 error 接口的方法，返回值为 string 类型
		response.FailWithMessage(err.Error(), c)
		return
	}
	//先答对验证码，再发邮箱
	//判断条件是：用户提交的验证码是否有效（正确且未过期）
	if store.Verify(req.CaptchaID, req.Captcha, true) {
		err := baseService.SendEmailVerificationCode(c, req.Email)
		if err != nil {
			global.Log.Error("Failed to send email", zap.Error(err))
			response.FailWithMessage("Failed to send email", c)
			return
		}
		response.OkWithMessage("Email sent successfully", c)
		return
	}
	//验证码是错的
	response.FailWithMessage("It is a wrong captcha,please try again", c)
}

// QQLoginURL 返回 QQ 登录链接
func (baseApi *BaseApi) QQLoginURL(c *gin.Context) {
	url := global.Config.QQ.QQLoginURL()
	response.OkWithData(url, c)
}

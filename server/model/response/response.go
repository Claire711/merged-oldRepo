package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	//将结构体序列化为 JSON 格式
	//通过上下文对象（c *gin.Context）直接向客户端发送响应，无需通过 return 返回值
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

// 接收参数，参数不同，响应数据相同
func OK(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "success", c)
}
func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}
func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "success", c)
}
func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}
func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "failure", c)
}
func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}
func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

// 通过业务状态码传递未授权信息
func NoAuth(message string, c *gin.Context) {
	//Gin 框架提供的快捷 Map 语法。gin.H{"reload": true}：map[string]interface{}{"reload": true}
	// "reload": true  核心：前端看到这个字段会触发页面重新加载（如跳转到登录页）
	Result(ERROR, gin.H{"reload": true}, message, c)
}

// 通过标准 HTTP 403 状态码传递权限拒绝
func Forbidden(message string, c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Code: ERROR,
		Data: nil,
		Msg:  message,
	})
}

//func 函数名(参数列表) 返回值类型 { ... }
/*
gin.H{"reload": true} 是后端向前端发送的一个特殊指令，
/告诉前端 "当前会话已失效，请重新加载页面或跳转到登录页"。
这是前后端分离架构中处理权限失效的常用模式，通过统一约定简化了权限管理逻辑。
*/

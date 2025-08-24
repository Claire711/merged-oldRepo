package initialize

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"server/global"
	"server/middleware"
	"server/router"
)

// InitRouter 初始化路由,返回类型：*gin.Engine
func InitRouter() *gin.Engine {
	// 设置gin模式
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()
	// 使用日志记录中间件
	Router.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	//创建一个加密的 Cookie 存储。
	var store = cookie.NewStore([]byte(global.Config.System.SessionsSecret))
	//创建一个名为 "session" 的会话管理器，使用之前创建的 store 存储引擎。
	Router.Use(sessions.Sessions("session", store))
	Router.StaticFS(global.Config.Upload.Path, http.Dir(global.Config.Upload.Path))
	routerGroup := router.RouterGroupApp
	//publicGroup 所有人都可以访问
	publicGroup := Router.Group(global.Config.System.RouterPrefix)
	privateGroup := Router.Group(global.Config.System.RouterPrefix)
	//Use(middleware.JWTAuth())：为该分组绑定 JWT 认证中间件
	privateGroup.Use(middleware.JWTAuth())
	adminGroup := Router.Group(global.Config.System.RouterPrefix)
	/*
		绑定了两个中间件：
		middleware.JWTAuth()：先验证 JWT Token（确保已登录）
		middleware.AdminAuth()：再验证用户是否为管理员身份
	*/
	adminGroup.Use(middleware.JWTAuth()).Use(middleware.AdminAuth())
	//public模块是通用的，用户注册后理应也可以查看
	{
		routerGroup.InitBaseRouter(publicGroup)
	}
	{
		routerGroup.InitUserRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitArticleRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitCommentRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitFeedbackRouter(privateGroup, publicGroup, adminGroup)
	}
	{
		routerGroup.InitImageRouter(adminGroup)
		routerGroup.InitAdvertisementRouter(adminGroup, publicGroup)
		routerGroup.InitFriendLinkRouter(adminGroup, publicGroup)
		routerGroup.InitWebsiteRouter(adminGroup, publicGroup)
		routerGroup.InitConfigRouter(adminGroup)
	}
	return Router
}

package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
	"server/middleware"
)

type UserRouter struct {
}

func (u *UserRouter) InitUserRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	userRouter := Router.Group("user")
	userPublicRouter := PublicRouter.Group("user")
	userLoginRouter := PublicRouter.Group("user").Use(middleware.LoginRecord())
	userAdminRouter := AdminRouter.Group("user")
	userApi := api.ApiGroupApp.UserApi //cha
	{
		userRouter.POST("logout", userApi.Logout)
		userRouter.PUT("resetPassword", userApi.UserResetPassword)
		userRouter.GET("info", userApi.UserInfo)
		userRouter.PUT("changeInfo", userApi.UserChangeInfo)
		userRouter.GET("weather", userApi.UserWeather)
		userRouter.GET("chart", userApi.UserChart)
	}
	{
		userPublicRouter.POST("forgetPassword", userApi.ForgetPassword)
		userPublicRouter.GET("card", userApi.UserCard)
	}
	{
		userLoginRouter.POST("register", userApi.Register)
		userLoginRouter.POST("login", userApi.Login)
	}
	{
		userAdminRouter.GET("list", userApi.UserList)
		userAdminRouter.PUT("freeze", userApi.UserFreeze)
		userAdminRouter.PUT("ubfreeze", userApi.UserUnfreeze)
		userAdminRouter.GET("loginList", userApi.UserLoginList)
	}
}

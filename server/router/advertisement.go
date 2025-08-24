package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
)

type AdvertisementRouter struct {
}

func (a *AdvertisementRouter) InitAdvertisementRouter(AdminRouter *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	advertisementAdminRouter := AdminRouter.Group("advertisement")
	advertisementPublicRouter := PublicRouter.Group("advertisement")
	advertisementApi := api.ApiGroupApp.AdvertisementApi
	{
		advertisementPublicRouter.GET("info", advertisementApi.AdvertisementInfo)
	}
	{
		advertisementAdminRouter.POST("create", advertisementApi.AdvertisementCreate)
		advertisementAdminRouter.DELETE("delete", advertisementApi.AdvertisementDelete)
		advertisementAdminRouter.PUT("update", advertisementApi.AdvertisementUpdate)
		advertisementAdminRouter.GET("list", advertisementApi.AdvertisementList)
	}
}

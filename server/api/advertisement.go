package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/model/request"
	"server/model/response"
)

type AdvertisementApi struct {
}

func (advertisementApi *AdvertisementApi) AdvertisementInfo(c *gin.Context) {
	//用户直接获取信息，不需要请求
	list, total, err := advertisementService.AdvertisementInfo()
	if err != nil {
		global.Log.Error("Failed to get advertisement information:", zap.Error(err))
		response.FailWithMessage("Failed to get advertisement information", c)
		return

	}
	response.OkWithData(response.AdvertisementInfo{
		List:  list,
		Total: total,
	}, c)
}

func (advertisementApi *AdvertisementApi) AdvertisementCreate(c *gin.Context) {
	var req request.AdvertisementCreate
	if err := c.ShouldBindJSON(&req); err != nil {

		response.FailWithMessage(err.Error(), c)
		return
	}
	err := advertisementService.AdvertisementCreate(req)
	if err != nil {
		global.Log.Error("Failed to create advertisement:", zap.Error(err))
		response.FailWithMessage("Failed to create advertisement", c)
		return
	}
	response.OkWithMessage("Successfully created advertisement", c)
}

func (advertisementApi *AdvertisementApi) AdvertisementDelete(c *gin.Context) {
	var req request.AdvertisementDelete
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := advertisementService.AdvertisementDelete(req)
	if err != nil {
		global.Log.Error("Failed to delete advertisement:", zap.Error(err))
		response.FailWithMessage("Failed to delete advertisement", c)
		return
	}
	response.OkWithMessage("Successfully deleted advertisement", c)
}

func (advertisementApi *AdvertisementApi) AdvertisementUpdate(c *gin.Context) {
	var req request.AdvertisementUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := advertisementService.AdvertisementUpdate(req)
	if err != nil {
		global.Log.Error("Failed to update advertisement:", zap.Error(err))
		response.FailWithMessage("Failed to update advertisement", c)
		return
	}
	response.OkWithMessage("Successfully updated advertisement", c)
}

func (advertisementApi *AdvertisementApi) AdvertisementList(c *gin.Context) {
	var info request.AdvertisementList
	if err := c.ShouldBindQuery(&info); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := advertisementService.AdvertisementList(info)
	if err != nil {
		global.Log.Error("Failed to get advertisement list:", zap.Error(err))
		response.FailWithMessage("Failed to get advertisement list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

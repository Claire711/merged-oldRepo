package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/utils"
)

type FeedbackApi struct {
}

func (feedbackApi *FeedbackApi) FeedbackNew(c *gin.Context) {
	list, err := feedbackService.FeedbackNew()
	if err != nil {
		global.Log.Error("Failed to get new feedback:", zap.Error(err))
		response.FailWithMessage("Failed to get new feedback", c)
		return
	}
	response.OkWithData(list, c)
}

func (feedbackApi *FeedbackApi) FeedbackCreate(c *gin.Context) {
	var req request.FeedbackCreate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UUID = utils.GetUUID(c)
	err = feedbackService.FeedbackCreate(req)
	if err != nil {
		global.Log.Error("Failed to create feedback:", zap.Error(err))
		response.FailWithMessage("Failed to create feedback", c)
		return
	}
	response.OkWithMessage("Successfully created feedback", c)
}

func (feedbackApi *FeedbackApi) FeedbackInfo(c *gin.Context) {
	uuid := utils.GetUUID(c)
	list, err := feedbackService.FeedbackInfo(uuid)
	if err != nil {
		global.Log.Error("Failed to get feedback information:", zap.Error(err))
		response.FailWithMessage("Failed to get feedback information", c)
		return
	}
	response.OkWithData(list, c)
}

func (feedbackApi *FeedbackApi) FeedbackDelete(c *gin.Context) {
	var req request.FeedbackDelete
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = feedbackService.FeedbackDelete(req)
	if err != nil {
		global.Log.Error("Failed to delete feedback:", zap.Error(err))
		response.FailWithMessage("Failed to delete feedback", c)
		return
	}
	response.OkWithMessage("Successfully deleted feedback", c)
}

func (feedbackApi *FeedbackApi) FeedbackReply(c *gin.Context) {
	var req request.FeedbackReply
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = feedbackService.FeedbackReply(req)
	if err != nil {
		global.Log.Error("Failed to reply feedback:", zap.Error(err))
		response.FailWithMessage("Failed to reply feedback", c)
		return
	}
	response.OkWithMessage("Successfully replied to feedback", c)
}

func (feedbackApi *FeedbackApi) FeedbackList(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := feedbackService.FeedbackList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get feedback list:", zap.Error(err))
		response.FailWithMessage("Failed to get feedback list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

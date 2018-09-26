package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

func TrackCreate(c *gin.Context) {

}

func TrackGetById(c *gin.Context) {
	name := c.Param("id")
	result, err := service.GetTrackerById(name)
	sendObjResponse(result, err, c)
}

func TrackAll(c *gin.Context) {
	userId := getUserId(c)
	result, err := service.GetAllTrackersForUser(userId)
	sendObjResponse(result, err, c)
}

func TrackByIds(c *gin.Context) {
	req := model.TracksByIdsRequest{}
	c.Bind(&req)
	result, err := service.GetTrackersByIds(req.Ids)
	sendObjResponse(result, err, c)
}

func TrackDeleteById(c *gin.Context) {

}

func TrackUpdate(c *gin.Context) {
	trackId := c.Param("id")
	userId := getUserId(c)
	req := model.TracksNameRequest{TrackId: trackId, UserId: userId}
	c.Bind(&req)
	err := service.UpdateTrackerName(&req)
	sendResponse(err, c)
}

func TrackerAvatar(c *gin.Context) {

}

func TrackerLastGeoPosition(c *gin.Context) {

}

func TrackerHistory(c *gin.Context) {

}

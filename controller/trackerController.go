package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

func TrackCreate(c *gin.Context) {
	req := model.TrackCreateRequest{UserId: getUserId(c)}
	c.Bind(&req)
	id, err := service.CreateTracker(&req)
	sendObjResponse(model.TrackCreateResponse{Id: id}, err, c)
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
	trackId := c.Param("id")
	userId := getUserId(c)
	err := service.DeleteTrack(userId, trackId)
	sendResponse(err, c)
}

func TrackUpdate(c *gin.Context) {
	req := model.TrackPrefRequest{TrackId: c.Param("id"), UserId: getUserId(c)}
	c.Bind(&req)
	err := service.UpdateTrackPref(&req)
	sendResponse(err, c)
}

func TrackerAvatar(c *gin.Context) {
	trackId := c.Param("id")
	file, e := c.FormFile("avatar")
	if e != nil {
		sendResponse([]int{}, c)
	}
	req := model.UpdateTrackAvatarRequest{UserId: getUserId(c), TrackId: trackId, File: file}
	avatar, err := service.UpdateTrackAvatar(&req)
	sendObjResponse(gin.H{"url": avatar}, err, c)
}

func TrackerHistory(c *gin.Context) {
	req := model.TracksHistRequest{}
	c.Bind(&req)
	resp, err := service.GetTrackHistory(&req)
	sendObjResponse(resp, err, c)
}

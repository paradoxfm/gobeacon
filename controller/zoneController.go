package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

func ZoneAllForUser(c *gin.Context) {
	req := model.ZoneAllRequest{UserId: getUserId(c)}
	c.Bind(&req)
	result, err := service.ZoneGetAllForUser(&req)
	sendObjResponse(result, err, c)
}

func ZoneCreate(c *gin.Context) {
	req := model.ZoneCreateRequest{UserId: getUserId(c)}
	c.Bind(&req)

	result, err := service.ZoneCreateForUser(&req)
	sendObjResponse(result, err, c)
}

func ZoneUpdate(c *gin.Context) {
	req := model.ZoneCreateRequest{UserId: getUserId(c)}
	c.Bind(&req)

	err := service.ZoneUpdate(&req)
	sendResponse(err, c)
}

func ZoneDeleteById(c *gin.Context) {
	zoneId := c.Param("id")
	err := service.ZoneDelete(zoneId)
	sendResponse(err, c)
}

func ZoneGetById(c *gin.Context) {
	zoneId := c.Param("id")
	zn, err := service.ZoneGetById(zoneId)
	sendObjResponse(zn, err, c)
}

func ZoneSnapTrackList(c *gin.Context) {
	zoneId := c.Param("id")
	req := model.ZoneSnapRequest{UserId: getUserId(c)}
	c.Bind(&req)

	err := service.ZoneSnapTrack(zoneId, &req)
	sendResponse(err, c)
}

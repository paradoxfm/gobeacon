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

func ZoneAdd(c *gin.Context) {

}

func ZoneDeleteById(c *gin.Context) {

}

func ZoneGetById(c *gin.Context) {

}

func ZoneUpdate(c *gin.Context) {

}

func ZoneSnapTrackList(c *gin.Context) {

}
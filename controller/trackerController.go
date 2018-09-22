package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/service"
)

func TrackCreate(c *gin.Context) {

}

func TrackGetById(c *gin.Context) {
	name := c.Param("id")
	result, err := service.GetTrackerById(name)
	sendObjResponse(result, err, c)
}

func TrackDeleteById(c *gin.Context) {

}

func TrackUpdate(c *gin.Context) {

}

func TrackerAvatar(c *gin.Context) {

}

func TrackerLastGeoPosition(c *gin.Context) {

}

func TrackerHistory(c *gin.Context) {

}
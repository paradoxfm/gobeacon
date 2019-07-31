package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

func BuySubscription(c *gin.Context) {
	//_ := getUserId(c)
	var req model.RegistrationRequest
	c.Bind(&req)

}

func CurrentSubscription(c *gin.Context) {
	userId := getUserId(c)
	result, err := service.CurrentSubscription(userId)
	sendObjResponse(result, err, c)
}

func AllActiveSubscription(c *gin.Context) {
	userId := getUserId(c)
	result, err := service.AllActiveSubscription(userId)
	sendObjResponse(result, err, c)
}

func Subscriptions(c *gin.Context) {
	result, err := service.Subscriptions()
	sendObjResponse(result, err, c)
}

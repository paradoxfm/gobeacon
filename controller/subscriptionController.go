package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/code"
	"gobeacon/model"
	"gobeacon/service"
)

func BuySubscription(c *gin.Context) {
	//_ := getUserId(c)
	req := model.BuySubscriptionRequest{UserId:getUserId(c)}
	if e := c.Bind(&req); e != nil {
		sendResponse([]int{code.ParseRequest}, c)
	}
	err := service.BuySubscription(&req)
	sendResponse(err, c)
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

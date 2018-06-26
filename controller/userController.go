package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var req model.RegistrationRequest
	c.Bind(&req)

	result, err := service.RegistrationUser(&req)

	if result {
		c.AbortWithStatus(http.StatusOK)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

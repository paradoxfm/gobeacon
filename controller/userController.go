package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
	"net/http"
	"github.com/appleboy/gin-jwt"
)

func UserCreate(c *gin.Context) {
	var req model.RegistrationRequest
	c.Bind(&req)

	result, err := service.RegistrationUser(&req)
	sendResponse(result, err, c)
}

func UserResetPassword(c *gin.Context) {
	req := model.ResetPasswordRequest{}
	c.Bind(&req)
	result, err := service.ResetPassword(&req)
	sendResponse(result, err, c)
}

func UserChangePassword(c *gin.Context) {
	req := model.ChangePasswordRequest{UserId: getUserId(c)}
	c.Bind(&req)
	result, err := service.ChangePassword(&req)
	sendResponse(result, err, c)
}

func UserGetProfile(c *gin.Context) {
	req := model.GetProfileRequest{UserId: getUserId(c)}
	c.Bind(&req)
}

func UserUpdateAvatar(c *gin.Context) {
	file, e := c.FormFile("avatar")
	req := model.UpdateAvatarRequest{UserId: getUserId(c), File: file}
	//cont, _ := file.Open()
	if e != nil {
		c.Bind(&req)
	}
}

func UserUpdatePushId(c *gin.Context) {
	req := model.GetProfileRequest{UserId: getUserId(c)}
	c.Bind(&req)
}

func getUserId(c *gin.Context) (string) {
	claims := jwt.ExtractClaims(c)
	return claims["id"].(string)
}

func sendResponse(goodWork bool, err []string, c *gin.Context) {
	if goodWork {
		c.AbortWithStatus(http.StatusOK)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

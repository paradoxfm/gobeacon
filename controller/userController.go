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

	_, err := service.RegistrationUser(&req)
	sendResponse(err, c)
}

func UserResetPassword(c *gin.Context) {
	req := model.ResetPasswordRequest{}
	c.Bind(&req)
	_, err := service.ResetPassword(&req)
	sendResponse(err, c)
}

func UserChangePassword(c *gin.Context) {
	req := model.ChangePasswordRequest{UserId: getUserId(c)}
	c.Bind(&req)
	_, err := service.ChangePassword(&req)
	sendResponse(err, c)
}

func UserGetProfile(c *gin.Context) {
	req := model.GetProfileRequest{UserId: getUserId(c)}
	c.Bind(&req)
	result, err := service.UserGetProfile(&req)
	sendObjResponse(result, err, c)
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
	req := model.UpdatePushRequest{UserId: getUserId(c)}
	c.Bind(&req)
	result, err := service.UserUpdatePushId(&req)
	sendObjResponse(result, err, c)
}

func getUserId(c *gin.Context) (string) {
	claims := jwt.ExtractClaims(c)
	return claims["id"].(string)
}

func sendResponse(err []int, c *gin.Context) {
	if len(err) == 0 {
		c.AbortWithStatus(http.StatusOK)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
	}
}

func sendObjResponse(obj interface{}, err []int, c *gin.Context) {
	if len(err) == 0 {
		c.JSON(http.StatusOK, obj)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
		})
	}
}

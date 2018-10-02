package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
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
	err := service.ChangePassword(&req)
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
	if e != nil {
		sendResponse([]int{}, c)
	}
	req := model.UpdateAvatarRequest{UserId: getUserId(c), File: file}
	avatar, err := service.UpdateUserAvatar(&req)
	sendObjResponse(gin.H{"url": avatar}, err, c)
}

func GetAvatar(c *gin.Context) {
	id := c.Param("id")
	result, err := service.GetAvatar(id)
	sendObjResponse(result, err, c)
}

func UserUpdatePushId(c *gin.Context) {
	req := model.UpdatePushRequest{UserId: getUserId(c)}
	c.Bind(&req)
	_, err := service.UserUpdatePushId(&req)
	sendResponse(err, c)
}

func TestPush(c *gin.Context) {
	userId := getUserId(c)
	service.SendPushNotification(userId)
	sendResponse(nil, c)
}

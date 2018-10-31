package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

// UserCreate godoc
// @Summary User registration
// @Description User registration email + password
// @Accept json
// @Produce json
// @Param request body model.RegistrationRequest true "Registration Request"
// @Router /users/signUp [post]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Users
func UserCreate(c *gin.Context) {
	var req model.RegistrationRequest
	c.Bind(&req)

	_, err := service.RegistrationUser(&req)
	sendResponse(err, c)
}

// UserResetPassword godoc
// @Summary Password reset
// @Description Reset user password by email
// @Accept json
// @Produce json
// @Param request body model.ResetPasswordRequest true "Password reset request"
// @Router /users/password/reset [post]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Users
func UserResetPassword(c *gin.Context) {
	req := model.ResetPasswordRequest{}
	c.Bind(&req)
	_, err := service.ResetPassword(&req)
	sendResponse(err, c)
}

// UserChangePassword godoc
// @Summary Change password
// @Description Change user password
// @Accept json
// @Produce json
// @Param request body model.ChangePasswordRequest true "Password change request"
// @Router /users/me/password [put]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Users
func UserChangePassword(c *gin.Context) {
	req := model.ChangePasswordRequest{UserId: getUserId(c)}
	c.Bind(&req)
	err := service.ChangePassword(&req)
	sendResponse(err, c)
}

// UserGetProfile godoc
// @Summary Profile request
// @Description Getting a user profile
// @Produce json
// @Router /users/me [get]
// @Success 200 {object} model.ProfileResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Users
func UserGetProfile(c *gin.Context) {
	req := model.GetProfileRequest{UserId: getUserId(c)}
	c.Bind(&req)
	result, err := service.UserGetProfile(&req)
	sendObjResponse(result, err, c)
}

// UserUpdateAvatar godoc
// @Summary Update user avatar
// @Description Update user avatar (jpeg 250x250)
// @Accept json
// @Produce json
// @Param avatar body file true "avatar"
// @Router /users/me/avatar [put]
// @Success 200 {object} model.AvatarIdResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Users
func UserUpdateAvatar(c *gin.Context) {
	file, e := c.FormFile("avatar")
	if e != nil {
		sendResponse([]int{}, c)
	}
	req := model.UpdateAvatarRequest{UserId: getUserId(c), File: file}
	avatar, err := service.UpdateUserAvatar(&req)
	sendObjResponse(model.AvatarIdResponse{Id: avatar}, err, c)
}

// UserUpdatePushId godoc
// @Summary Update user push id
// @Description Update user push id
// @Accept json
// @Produce json
// @Param request body model.UpdatePushRequest true "Request with push id"
// @Router /users/me/push [put]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Users
func UserUpdatePushId(c *gin.Context) {
	req := model.UpdatePushRequest{UserId: getUserId(c)}
	c.Bind(&req)
	_, err := service.UserUpdatePushId(&req)
	sendResponse(err, c)
}
package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

// UserCreate godoc
// @Summary Регистрация пользователя
// @Description Регистрация пользователя email + password
// @Accept json
// @Produce json
// @Param request body model.RegistrationRequest true "Запрос на регистрацию"
// @Router /users/signUp [post]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Пользователи
func UserCreate(c *gin.Context) {
	var req model.RegistrationRequest
	c.Bind(&req)

	_, err := service.RegistrationUser(&req)
	sendResponse(err, c)
}

// UserResetPassword godoc
// @Summary Сброс пароля
// @Description Сброс пароля пользователя по email
// @Accept json
// @Produce json
// @Param request body model.ResetPasswordRequest true "Запрос на сброс пароля"
// @Router /users/password/reset [post]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Пользователи
func UserResetPassword(c *gin.Context) {
	req := model.ResetPasswordRequest{}
	c.Bind(&req)
	_, err := service.ResetPassword(&req)
	sendResponse(err, c)
}

// UserChangePassword godoc
// @Summary Изменение пароля
// @Description Изменение пароля пользователя
// @Accept json
// @Produce json
// @Param request body model.ChangePasswordRequest true "Запрос на изменение пароля"
// @Router /users/me/password [put]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Пользователи
func UserChangePassword(c *gin.Context) {
	req := model.ChangePasswordRequest{UserId: getUserId(c)}
	c.Bind(&req)
	err := service.ChangePassword(&req)
	sendResponse(err, c)
}

// UserGetProfile godoc
// @Summary Запрос профиля
// @Description Получение профиля пользователя
// @Produce json
// @Router /users/me [get]
// @Success 200 {object} model.ProfileResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Пользователи
func UserGetProfile(c *gin.Context) {
	req := model.GetProfileRequest{UserId: getUserId(c)}
	c.Bind(&req)
	result, err := service.UserGetProfile(&req)
	sendObjResponse(result, err, c)
}

// UserUpdateAvatar godoc
// @Summary Обновление аватара пользователя
// @Description Обновление аватара пользователя (jpeg 250x250)
// @Accept json
// @Produce json
// @Param avatar body file true "avatar"
// @Router /users/me/avatar [put]
// @Success 200 {object} model.AvatarIdResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Пользователи
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
// @Summary Обновление push id пользователя
// @Description Обновление push id пользователя
// @Accept json
// @Produce json
// @Param request body model.UpdatePushRequest true "Запрос с push id"
// @Router /users/me/push [put]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Пользователи
func UserUpdatePushId(c *gin.Context) {
	req := model.UpdatePushRequest{UserId: getUserId(c)}
	c.Bind(&req)
	_, err := service.UserUpdatePushId(&req)
	sendResponse(err, c)
}
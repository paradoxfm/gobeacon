package controller

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
	"time"
)

var identityKey = "private_claim_id"

func CreateAdminJWTMiddleware() (*jwt.GinJWTMiddleware) {
	authMiddleware := &jwt.GinJWTMiddleware{// the jwt middleware
		Realm:         "jwt auth",
		Key:           []byte("z%~3povvo2tF?L3mlcvmlzoTN6i7dl"),
		Timeout:       time.Hour * 24 * 365 * 25, // время действия токена после авторизации
		MaxRefresh:    0, // время действия токена после обновления
		Authenticator: getAuthenticator,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(string); ok {
				return jwt.MapClaims{
					identityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
	}

	return authMiddleware
}

// AuthorizeUser godoc
// @Summary Авторизация пользователя
// @Description Авторизация пользователя
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Логин пароль"
// @Router /users/login [post]
// @Success 200 {object} model.LoginResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Пользователи
func getAuthenticator(c *gin.Context) (interface{}, error) {
	var cred model.LoginRequest
	if err := c.ShouldBind(&cred); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	return service.LoginUser(cred.Email, cred.Password)
}

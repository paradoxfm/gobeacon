package controller

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
	"time"
)

var identityKey = "private_claim_id"
var pwdKey = "private_claim_pwd"

func CreateAdminJWTMiddleware() (*jwt.GinJWTMiddleware) {
	authMiddleware := &jwt.GinJWTMiddleware{// the jwt middleware
		Realm:         "jwt auth",
		Key:           []byte("z%~3povvo2tF?L3mlcvmlzoTN6i7dl"),
		Timeout:       time.Hour * 24 * 365 * 25, // время действия токена после авторизации
		MaxRefresh:    0, // время действия токена после обновления
		Authenticator: getAuthenticator,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if usr, ok := data.(model.UserDb); ok {
				return jwt.MapClaims{
					identityKey: usr.Id.String(),
					pwdKey: usr.Password,
				}
			}
			return jwt.MapClaims{}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			claims := jwt.ExtractClaims(c)
			id := claims[identityKey]
			pwd := claims[pwdKey]
			exist := service.UserExist(id.(string), pwd.(string))
			//exist = false
			if !exist {
				c.AbortWithStatus(401)
				return false
			}
			return exist
		},
	}

	return authMiddleware
}

// AuthorizeUser godoc
// @Summary User Authorization
// @Description User Authorization by Login/Password
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login: Password"
// @Router /users/login [post]
// @Success 200 {object} model.LoginResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Users
func getAuthenticator(c *gin.Context) (interface{}, error) {
	var cred model.LoginRequest
	if err := c.ShouldBind(&cred); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	user, e := service.LoginUser(cred.Email, cred.Password)
	return user, e
}

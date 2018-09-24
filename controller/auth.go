package controller

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"gobeacon/service"
	"time"
)

type login struct {
	Email string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}


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

//авторизация пользователя
func getAuthenticator(c *gin.Context) (interface{}, error) {
	var cred login
	if err := c.ShouldBind(&cred); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	return service.LoginUser(cred.Email, cred.Password)
}

func CreateHeartGinJWTMiddleware() (*jwt.GinJWTMiddleware) {
	authMiddleware := &jwt.GinJWTMiddleware{// the jwt middleware
		Realm:         "megazlo.net",
		Key:           []byte("fkl;jnbLJkN;old"),
		Timeout:       time.Hour, // время действия токена после авторизации
		MaxRefresh:    time.Hour, // время действия токена после обновления
		Authenticator: func(c *gin.Context) (interface{}, error) {
			claims := jwt.ExtractClaims(c)
			if claims["id"] == "heart349023" && claims["password"] == "s156EzI07820CtsfJhu" {
				return claims["id"], nil
			}
			return "", jwt.ErrMissingLoginValues
		},
	}

	return authMiddleware
}

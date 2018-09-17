package controller

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"time"
	"gobeacon/service"
)

func CreateAdminJWTMiddleware() (*jwt.GinJWTMiddleware) {
	authMiddleware := &jwt.GinJWTMiddleware{// the jwt middleware
		Realm:         "jwt auth",
		Key:           []byte("z%~3povvo2tF?L3mlcvmlzoTN6i7dl"),
		Timeout:       time.Hour * 24 * 365 * 25, // время действия токена после авторизации
		MaxRefresh:    0, // время действия токена после обновления
		Authenticator: getAuthenticator,
	}

	return authMiddleware
}

//авторизация пользователя
func getAuthenticator(userID string, password string, c *gin.Context) (interface{}, bool) {
	return service.LoginUser(userID, password)
	/*if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return userID, true
	}
	return "", false*/
}

func CreateHeartGinJWTMiddleware() (*jwt.GinJWTMiddleware) {
	authMiddleware := &jwt.GinJWTMiddleware{// the jwt middleware
		Realm:         "megazlo.net",
		Key:           []byte("fkl;jnbLJkN;old"),
		Timeout:       time.Hour, // время действия токена после авторизации
		MaxRefresh:    time.Hour, // время действия токена после обновления
		Authenticator: func(userID string, password string, c *gin.Context) (interface{}, bool) {
			if userID == "heart349023" && password == "s156EzI07820CtsfJhu" {
				return userID, true
			}
			return "", false
		},
	}

	return authMiddleware
}

package controller

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"time"
)

func CreateAdminJWTMiddleware() (*jwt.GinJWTMiddleware) {
	authMiddleware := &jwt.GinJWTMiddleware{// the jwt middleware
		Realm:         "megazlo.net",
		Key:           []byte("fkl;jnbLJkN;old"),
		Timeout:       time.Hour, // время действия токена после авторизации
		MaxRefresh:    time.Hour, // время действия токена после обновления
		Authenticator: getAuthenticator,
		//Authorizator: getAuthorizator,// проверка на доступ к методам и тп
		Unauthorized: processUnauthorized,// ответ в случае, если не авторизован или неправильная авторизация
	}

	return authMiddleware
}

// обработка в случае, если пользователь не авторизован
func processUnauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

//авторизация пользователя
func getAuthenticator(userID string, password string, c *gin.Context) (string, bool) {
	if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
		return userID, true
	}
	return "", false
}

func CreateHeartGinJWTMiddleware() (*jwt.GinJWTMiddleware) {
	authMiddleware := &jwt.GinJWTMiddleware{// the jwt middleware
		Realm:         "megazlo.net",
		Key:           []byte("fkl;jnbLJkN;old"),
		Timeout:       time.Hour, // время действия токена после авторизации
		MaxRefresh:    time.Hour, // время действия токена после обновления
		Authenticator: func(userID string, password string, c *gin.Context) (string, bool) {
			if userID == "heart349023" && password == "s156EzI07820CtsfJhu" {
				return userID, true
			}
			return "", false
		},
		//Authorizator: getAuthorizator,// проверка на доступ к методам и тп
		Unauthorized: processUnauthorized,// ответ в случае, если не авторизован или неправильная авторизация
	}

	return authMiddleware
}



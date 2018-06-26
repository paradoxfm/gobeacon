package controller

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// ListAccounts godoc
// @Summary List accounts
// @Description get accounts
// @Accept  json
// @Produce  json
// @Param q query string false "name search by q"
// @Success 200 {array} model.UserAuth
// @Router /users [get]
// @Tags Пользователи users
func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"userID": claims["id"],
		"text":   "Hello World.",
	})
}

// ListAccounts godoc
// @Summary List accounts
// @Description get accounts
// @Accept  json
// @Produce  json
// @Param q query string false "name search by q"
// @Success 200 {array} model.UserAuth
// @Router /users [post]
// @Tags Пользователи users
func HelloHandler2(c *gin.Context) {

}

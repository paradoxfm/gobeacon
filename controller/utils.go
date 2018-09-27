package controller

import (
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getUserId(c *gin.Context) (string) {
	claims := jwt.ExtractClaims(c)
	if val, ok := claims["private_claim_id"]; ok {
		return val.(string)
	}
	return ""
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


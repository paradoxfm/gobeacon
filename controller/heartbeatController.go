package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

func HeartbeatPhone(c *gin.Context) {
	req := model.Heartbeat{}
	if e := c.Bind(&req); e != nil {
		fmt.Println(e)
		return
	}
	_, err := service.SaveHeartbeat(&req)
	if err == nil { // если нет ошибок, обновляем последние данные и выполняем проверки по трекеру
		//go service.CheckAndUpdateTracker(trk)
	}

	sendResponse(err, c)
}

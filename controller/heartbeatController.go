package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

// HeartbeatPhone godoc
// @Summary Прием координат от смартфона
// @Description Принимает координаты от смартфона, происходит: сохранение и проверка на вхождение в неозоны, обновление последних данных трекера
// @securityDefinitions.basic BasicAut
// @Accept json
// @Param request body model.Heartbeat true "Запрос с координатами и флагами устройства"
// @Success 200 "ok"
// @Router /heartbeat [post]
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Мониторинг
func HeartbeatPhone(c *gin.Context) {
	req := model.Heartbeat{}
	if e := c.Bind(&req); e != nil {
		fmt.Println(e)
		return
	}
	err := service.SaveHeartbeat(&req)

	sendResponse(err, c)
}

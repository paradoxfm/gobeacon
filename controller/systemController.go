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

// HeartbeatPhone godoc
// @Summary Загрузка аватарки по id
// @Accept json
// @Produce json
// @Param id path string true "id аватарки"
// @Router /avatar/{id} [get]
// @Success 200 {object} model.AvatarResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Система
func GetAvatar(c *gin.Context) {
	id := c.Param("id")
	result, err := service.GetAvatar(id)
	sendObjResponse(result, err, c)
}

/*func GetAvatar(c *gin.Context) {
	id := c.Param("id")
	imgData, err := service.GetAvatar(id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	header := c.Writer.Header()
	header["Content-type"] = []string{"image/jpg"}
	header["Content-Disposition"] = []string{"attachment; filename= " + id}

	io.Copy(c.Writer, bytes.NewReader(imgData))
}*/

func TestPush(c *gin.Context) {
	userId := getUserId(c)
	service.SendPushNotification(userId)
	sendResponse(nil, c)
}


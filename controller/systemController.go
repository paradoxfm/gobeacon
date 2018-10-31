package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

// HeartbeatPhone godoc
// @Summary Acceptance of coordinates from a smartphone
// @Description Accepts coordinates from the smartphone, occurs: saving and checking for entry into neozones, updating the latest tracker data
// @securityDefinitions.basic BasicAut
// @Accept json
// @Param request body model.Heartbeat true "Request with device coordinates and flags"
// @Success 200 "ok"
// @Router /heartbeat [post]
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Monitoring
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
// @Summary Loading avatar by id
// @Accept json
// @Produce json
// @Param id path string true "id avatar"
// @Router /avatar/{id} [get]
// @Success 200 {object} model.AvatarResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags System
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


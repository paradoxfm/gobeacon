package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

// ZoneAllForUser godoc
// @Summary Список зон
// @Description Получение списка зон для авторизованного пользователя
// @Accept json
// @Produce json
// @Success 200 {array} model.GeoZoneResponse
// @Router /zone/all [get]
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Зоны
func ZoneAllForUser(c *gin.Context) {
	req := model.ZoneAllRequest{UserId: getUserId(c)}
	c.Bind(&req)
	result, err := service.ZoneGetAllForUser(&req)
	sendObjResponse(result, err, c)
}

// ZoneCreate godoc
// @Summary Создание зоны
// @Description Создание зоны для авторизованного пользователя
// @Accept json
// @Produce json
// @Param request body model.ZoneCreateRequest true "Запрос настройками зоны"
// @Success 200 {object} model.GeoZoneResponse
// @Router /zone/save [post]
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Зоны
func ZoneCreate(c *gin.Context) {
	req := model.ZoneCreateRequest{UserId: getUserId(c)}
	c.Bind(&req)

	result, err := service.ZoneCreateForUser(&req)
	sendObjResponse(result, err, c)
}

// ZoneUpdate godoc
// @Summary Обновление зоны
// @Description Обновление настроек зоны для авторизованного пользователя
// @Accept json
// @Produce json
// @Param id path string true "id зоны"
// @Param request body model.ZoneCreateRequest true "Запрос настройками зоны"
// @Router /zone/update/{id} [put]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Зоны
func ZoneUpdate(c *gin.Context) {
	req := model.ZoneCreateRequest{UserId: getUserId(c)}
	c.Bind(&req)

	err := service.ZoneUpdate(&req)
	sendResponse(err, c)
}

// ZoneDeleteById godoc
// @Summary Удаление зоны
// @Description Удаление зоны для авторизованного пользователя
// @Accept json
// @Produce json
// @Param id path string true "id зоны"
// @Router /zone/delete/{id} [delete]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Зоны
func ZoneDeleteById(c *gin.Context) {
	zoneId := c.Param("id")
	err := service.ZoneDelete(zoneId)
	sendResponse(err, c)
}

// ZoneGetById godoc
// @Summary Загрузка зоны
// @Description Загрузка зоны для авторизованного пользователя
// @Accept json
// @Produce json
// @Param id path string true "id зоны"
// @Router /zone/find/{id} [get]
// @Success 200 {object} model.GeoZoneResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Зоны
func ZoneGetById(c *gin.Context) {
	zoneId := c.Param("id")
	zn, err := service.ZoneGetById(zoneId)
	sendObjResponse(zn, err, c)
}

// ZoneSnapTrackList godoc
// @Summary Привязка зоны
// @Description Привязка зоны к списку трекеров
// @Accept json
// @Produce json
// @Param id path string true "id зоны"
// @Param request body model.ZoneSnapRequest true "Запрос привязки трекеров к зоне"
// @Router /zone/snap/{id} [put]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Зоны
func ZoneSnapTrackList(c *gin.Context) {
	zoneId := c.Param("id")
	req := model.ZoneSnapRequest{UserId: getUserId(c)}
	c.Bind(&req)

	err := service.ZoneSnapTrack(zoneId, &req)
	sendResponse(err, c)
}

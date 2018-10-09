package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

// TrackCreate godoc
// @Summary Создание нового трекера
// @Description Создание нового трекера для авторизованного пользователя
// @Accept json
// @Produce json
// @Param request body model.TrackCreateRequest true "Запрос на создание треккера"
// @Router /trackers [post]
// @Success 200 {object} model.TrackCreateResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Трекеры
func TrackCreate(c *gin.Context) {
	req := model.TrackCreateRequest{UserId: getUserId(c)}
	c.Bind(&req)
	id, err := service.CreateTracker(&req)
	sendObjResponse(model.TrackCreateResponse{Id: id}, err, c)
}

// TrackGetById godoc
// @Summary Загрузить трекер
// @Description Загрузка трекера по id
// @Produce json
// @Param id path string true "id трекера"
// @Router /trackers/find/{id} [get]
// @Success 200 {object} model.Tracker
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Трекеры
func TrackGetById(c *gin.Context) {
	name := c.Param("id")
	result, err := service.GetTrackerById(name)
	sendObjResponse(result, err, c)
}

// TrackAll godoc
// @Summary Загрузить трекеры
// @Description Лист трекеров для пользователя
// @Produce json
// @Param id path string true "id трекера"
// @Router /trackers/all [get]
// @Success 200 {array} model.Tracker
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Трекеры
func TrackAll(c *gin.Context) {
	userId := getUserId(c)
	result, err := service.GetAllTrackersForUser(userId)
	sendObjResponse(result, err, c)
}


// TrackByIds godoc
// @Summary Загрузить трекеры по id
// @Description Лист выбранных трекеров для пользователя
// @Accept json
// @Produce json
// @Param request body model.TracksByIdsRequest true "Запрос на отдельные трекеры"
// @Router /trackers/custom [post]
// @Success 200 {array} model.Tracker
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Трекеры
func TrackByIds(c *gin.Context) {
	req := model.TracksByIdsRequest{}
	c.Bind(&req)
	result, err := service.GetTrackersByIds(req.Ids)
	sendObjResponse(result, err, c)
}

// ZoneDeleteById godoc
// @Summary Удаление трекера
// @Description Удаление трекера для авторизованного пользователя, если трекер привязан к нескольким пользователям, то только отвязка
// @Accept json
// @Produce json
// @Param id path string true "id трекера"
// @Router /trackers/delete/{id} [delete]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Трекеры
func TrackDeleteById(c *gin.Context) {
	trackId := c.Param("id")
	userId := getUserId(c)
	err := service.DeleteTrack(userId, trackId)
	sendResponse(err, c)
}

// TrackUpdate godoc
// @Summary Обновление трекера
// @Description Обновление свойств трекера для пользователя
// @Accept json
// @Produce json
// @Param id path string true "id трекера"
// @Param request body model.TrackPrefRequest true "Запрос на редактирование трекера"
// @Router /trackers/update/{id} [put]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Трекеры
func TrackUpdate(c *gin.Context) {
	req := model.TrackPrefRequest{TrackId: c.Param("id"), UserId: getUserId(c)}
	c.Bind(&req)
	err := service.UpdateTrackPref(&req)
	sendResponse(err, c)
}

// TrackerAvatar godoc
// @Summary Обновление аватара трекера
// @Description Обновление аватара трекера для пользователя (jpeg 250x250)
// @Accept json
// @Produce json
// @Param id path string true "id трекера"
// @Param avatar body file true "avatar"
// @Router /trackers/avatar/{id} [put]
// @Success 200 {object} model.AvatarIdResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Трекеры
func TrackerAvatar(c *gin.Context) {
	trackId := c.Param("id")
	file, e := c.FormFile("avatar")
	if e != nil {
		sendResponse([]int{}, c)
	}
	req := model.UpdateTrackAvatarRequest{UserId: getUserId(c), TrackId: trackId, File: file}
	avatar, err := service.UpdateTrackAvatar(&req)
	sendObjResponse(model.AvatarIdResponse{Id:avatar}, err, c)
}

// TrackerHistory godoc
// @Summary История трекера
// @Description Получение истории для трекера за период
// @Accept json
// @Produce json
// @Param request body model.TracksHistRequest true "Запрос истории треккера"
// @Router /trackers/geo/history [post]
// @Success 200 {array} model.TrackHistoryResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Трекеры
func TrackerHistory(c *gin.Context) {
	req := model.TracksHistRequest{}
	c.Bind(&req)
	resp, err := service.GetTrackHistory(&req)
	sendObjResponse(resp, err, c)
}

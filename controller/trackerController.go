package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

// TrackCreate godoc
// @Summary Create a new tracker
// @Description Creating a new tracker for an authorized user
// @Accept json
// @Produce json
// @Param request body model.TrackCreateRequest true "Запрос на создание треккера"
// @Router /trackers [post]
// @Success 200 {object} model.TrackCreateResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Trackers
func TrackCreate(c *gin.Context) {
	req := model.TrackCreateRequest{UserId: getUserId(c)}
	c.Bind(&req)
	id, err := service.CreateTracker(&req)
	sendObjResponse(model.TrackCreateResponse{Id: id}, err, c)
}

// TrackGetById godoc
// @Summary Load tracker
// @Description Load tracker by id
// @Produce json
// @Param id path string true "id tracker"
// @Router /trackers/find/{id} [get]
// @Success 200 {object} model.Tracker
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Trackers
func TrackGetById(c *gin.Context) {
	name := c.Param("id")
	result, err := service.GetTrackerById(name)
	sendObjResponse(result, err, c)
}

// TrackAll godoc
// @Summary Load all trackers
// @Description Load all trackers for an authorized user
// @Produce json
// @Router /trackers/all [get]
// @Success 200 {array} model.Tracker
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Trackers
func TrackAll(c *gin.Context) {
	userId := getUserId(c)
	result, err := service.GetAllTrackersForUser(userId)
	sendObjResponse(result, err, c)
}

// TrackByIds godoc
// @Summary Load custom trackers by ids
// @Description Load custom trackers by ids for an authorized user
// @Accept json
// @Produce json
// @Param request body model.TracksByIdsRequest true "Запрос на отдельные трекеры"
// @Router /trackers/custom [post]
// @Success 200 {array} model.Tracker
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Trackers
func TrackByIds(c *gin.Context) {
	req := model.TracksByIdsRequest{}
	c.Bind(&req)
	result, err := service.GetTrackersByIds(req.Ids)
	sendObjResponse(result, err, c)
}

// ZoneDeleteById godoc
// @Summary Delete tracker
// @Description Delete the tracker for an authorized user, if the tracker is associated with several users, then only unlink
// @Accept json
// @Produce json
// @Param id path string true "id tracker"
// @Router /trackers/delete/{id} [delete]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Trackers
func TrackDeleteById(c *gin.Context) {
	trackId := c.Param("id")
	userId := getUserId(c)
	err := service.DeleteTrack(userId, trackId)
	sendResponse(err, c)
}

// TrackUpdate godoc
// @Summary Update tracker
// @Description Update tracker properties for authorized user
// @Accept json
// @Produce json
// @Param id path string true "id трекера"
// @Param request body model.TrackPrefRequest true "Запрос на редактирование трекера"
// @Router /trackers/update/{id} [put]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Trackers
func TrackUpdate(c *gin.Context) {
	req := model.TrackPrefRequest{TrackId: c.Param("id"), UserId: getUserId(c)}
	c.Bind(&req)
	err := service.UpdateTrackPref(&req)
	sendResponse(err, c)
}

// TrackerAvatar godoc
// @Summary Update avatar tracker
// @Description Update avatar tracker for authorized user (jpeg 250x250)
// @Accept json
// @Produce json
// @Param id path string true "id tracker"
// @Param avatar body file true "avatar"
// @Router /trackers/avatar/{id} [put]
// @Success 200 {object} model.AvatarIdResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Trackers
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
// @Summary Tracker history
// @Description Getting history for tracker for the period
// @Accept json
// @Produce json
// @Param request body model.TracksHistRequest true "Request tracker history"
// @Router /trackers/geo/history [post]
// @Success 200 {array} model.TrackHistoryResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Trackers
func TrackerHistory(c *gin.Context) {
	req := model.TracksHistRequest{}
	c.Bind(&req)
	resp, err := service.GetTrackHistory(&req)
	sendObjResponse(resp, err, c)
}

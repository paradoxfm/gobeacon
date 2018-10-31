package controller

import (
	"github.com/gin-gonic/gin"
	"gobeacon/model"
	"gobeacon/service"
)

// ZoneAllForUser godoc
// @Summary Zone list
// @Description Getting a list of zones for an authorized user
// @Accept json
// @Produce json
// @Success 200 {array} model.GeoZoneResponse
// @Router /zone/all [get]
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Zones
func ZoneAllForUser(c *gin.Context) {
	req := model.ZoneAllRequest{UserId: getUserId(c)}
	c.Bind(&req)
	result, err := service.ZoneGetAllForUser(&req)
	sendObjResponse(result, err, c)
}

// ZoneCreate godoc
// @Summary Creating a zone
// @Description Creating a zone for an authorized user
// @Accept json
// @Produce json
// @Param request body model.ZoneCreateRequest true "Request with zone settings"
// @Success 200 {object} model.GeoZoneResponse
// @Router /zone/save [post]
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Zones
func ZoneCreate(c *gin.Context) {
	req := model.ZoneCreateRequest{UserId: getUserId(c)}
	c.Bind(&req)

	result, err := service.ZoneCreateForUser(&req)
	sendObjResponse(result, err, c)
}

// ZoneUpdate godoc
// @Summary Zone update
// @Description Updating zone settings for an authorized user
// @Accept json
// @Produce json
// @Param id path string true "id зоны"
// @Param request body model.ZoneCreateRequest true "Request with zone settings"
// @Router /zone/update/{id} [put]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Zones
func ZoneUpdate(c *gin.Context) {
	req := model.ZoneCreateRequest{UserId: getUserId(c)}
	c.Bind(&req)

	err := service.ZoneUpdate(&req)
	sendResponse(err, c)
}

// ZoneDeleteById godoc
// @Summary Deleting a zone
// @Description Deleting a zone for an authorized user
// @Accept json
// @Produce json
// @Param id path string true "id zone"
// @Router /zone/delete/{id} [delete]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Zones
func ZoneDeleteById(c *gin.Context) {
	zoneId := c.Param("id")
	err := service.ZoneDelete(zoneId)
	sendResponse(err, c)
}

// ZoneGetById godoc
// @Summary Loading zone
// @Description Loading zone for authorized user
// @Accept json
// @Produce json
// @Param id path string true "id zone"
// @Router /zone/find/{id} [get]
// @Success 200 {object} model.GeoZoneResponse
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Zones
func ZoneGetById(c *gin.Context) {
	zoneId := c.Param("id")
	zn, err := service.ZoneGetById(zoneId)
	sendObjResponse(zn, err, c)
}

// ZoneSnapTrackList godoc
// @Summary Zone binding
// @Description Linking the zone to the tracker list
// @Accept json
// @Produce json
// @Param id path string true "id зоны"
// @Param request body model.ZoneSnapRequest true "Request linking trackers to a zone"
// @Router /zone/snap/{id} [put]
// @Success 200 "ok"
// @Failure 400 "err"
// @Failure 500 "err"
// @Tags Zones
func ZoneSnapTrackList(c *gin.Context) {
	zoneId := c.Param("id")
	req := model.ZoneSnapRequest{UserId: getUserId(c)}
	c.Bind(&req)

	err := service.ZoneSnapTrack(zoneId, &req)
	sendResponse(err, c)
}

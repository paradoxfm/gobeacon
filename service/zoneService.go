package service

import (
	"gobeacon/code"
	"gobeacon/model"
)

func ZoneGetAllForUser(r *model.ZoneAllRequest) ([]model.GeoZoneResponse, []int) {
	var err []int
	zones, e := getAllZoneByUserId(r.UserId)
	if e != nil {
		err = append(err, code.DbErrorUpdateUserPush)
		return []model.GeoZoneResponse{}, err
	}
	var rez []model.GeoZoneResponse
	for _, zone := range zones {
		zn := model.GeoZoneResponse{Id:zone.Id.String(), Name:zone.Name, Points: []model.ZonePoint{}, Trackers: []model.TrackSnapZone{}}
		rez = append(rez, zn)
	}
	return rez, err
}

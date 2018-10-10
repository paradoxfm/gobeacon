package service

import (
	"gobeacon/code"
	"gobeacon/db"
	"gobeacon/model"
)

func ZoneGetAllForUser(r *model.ZoneAllRequest) ([]model.GeoZoneResponse, []int) {
	var err []int
	rez := make([]model.GeoZoneResponse, 0)
	zones, e := db.LoadZonesByUserId(r.UserId)
	if e != nil {
		err = append(err, code.DbErrorUpdateUserPush)
		return rez, err
	}
	for _, zone := range zones {
		zn := convertZoneToResponse(zone)
		rez = append(rez, zn)
	}
	return rez, err
}

func convertZoneToResponse(zone model.GeoZoneDb) (model.GeoZoneResponse) {
	zn := model.GeoZoneResponse{Id: zone.Id.String(), Name: zone.Name, Points: zone.Points}
	for key := range zone.Trackers {
		zn.Trackers = append(zn.Trackers, key.String())
	}
	return zn
}

func ZoneCreateForUser(r *model.ZoneCreateRequest) (interface{}, []int) {
	var err []int
	zn, e := db.CreateZoneForUser(r)
	if e != nil {
		err = append(err, code.DbErrorUpdateUserPush)
		return nil, err
	}

	rez := convertZoneToResponse(zn)
	rez.Trackers = make([]string, 0)
	return rez, err
}

func ZoneUpdate(r *model.ZoneCreateRequest) ([]int) {
	var err []int
	if e := db.UpdateZone(r); e != nil {
		return append(err, code.DbError)
	}
	return err
}

func ZoneGetById(id string) (interface{}, []int) {
	var err []int
	zone, e := db.LoadZoneById(id)
	if e != nil {
		return nil, append(err, code.DbError)
	}
	rez := convertZoneToResponse(zone)
	return rez, err
}

func ZoneDelete(zoneId string) ([]int) {
	var err []int
	if e := db.DeleteZoneById(zoneId); e != nil {
		return append(err, code.DbError)
	}
	return err
}

func ZoneSnapTrack(zoneId string, r *model.ZoneSnapRequest) ([]int) {
	var err []int

	track := make(map[string]bool)
	for _, v := range r.Ids {
		track[v] = true
	}

	if e := db.UpdateZoneTrackers(zoneId, track); e != nil {
		return append(err, code.DbError)
	}
	return err
}

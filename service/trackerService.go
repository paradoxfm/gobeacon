package service

import (
	"github.com/kellydunn/golang-geo"
	"gobeacon/code"
	"gobeacon/db"
	"gobeacon/model"
	"sort"
)

func GetTrackerById(id string) (interface{}, []int) {
	var err []int
	tracker, e := db.GetTrackerById(id)
	if e != nil {
		err = append(err, code.DbError)
	}
	return tracker, err
}

func CreateTracker(req *model.TrackCreateRequest) (string, []int) {
	var err []int
	// ищем может уже добавлял кто-то (ситуация с 2мя родителями)
	devId := req.DeviceId
	var e error
	var id interface{}
	if len(devId) == 15 { // если прислали imei часов
		req.Imei = req.DeviceId
		req.DeviceId = req.DeviceId[4:14]
		id, e = db.ExistTrackByDevice(req.DeviceId)
	} else {
		id, e = db.ExistTrackByDevice(devId)
	}
	if e != nil {
		return "", append(err, code.DbError)
	}
	if id == nil { // если ничего нет, то добавляем трекер
		if id, e = db.InsertNewTrack(req); e != nil {
			return "", append(err, code.DbError)
		}
	}
	// ищем текущую связь стрекером, на случай если решили повторно добавить (
	exist, e := db.ExistTrackPref(req.UserId, id.(string))
	if e != nil {
		return "", append(err, code.DbError)
	}
	if exist { // если связь уже есть отпинываем
		return "", append(err, code.TrackForUserExist)
	}
	if e = db.InsertNewTrackPref(id.(string), req); e != nil {
		return "", append(err, code.DbError)
	}
	return id.(string), err
}

func GetAllTrackersForUser(userId string) (interface{}, []int) {
	var err []int
	trackerList, e := db.GetTrackersByUserId(userId)
	if e != nil {
		err = append(err, code.DbErrorGetTracker)
	}
	if trackerList == nil {
		trackerList = make([]model.Tracker, 0)
	}
	return trackerList, err
}

func GetTrackersByIds(ids []string) (interface{}, []int) {
	var err []int
	trackerList, e := db.GetTrackerByIds(ids)
	if e != nil {
		err = append(err, code.DbErrorGetTracker)
	}
	return trackerList, err
}

func DeleteTrack(userId string, trackId string) ([]int) {
	var err []int
	e := db.DeleteTrackForUser(userId, trackId)
	if e != nil {
		err = append(err, code.DbError)
	}
	return err
}

func UpdateTrackPref(req *model.TrackPrefRequest) ([]int) {
	var err []int
	e := db.UpdateTrackPref(req)
	if e != nil {
		err = append(err, code.DbError)
	}
	return err
}

func UpdateTrackAvatar(req *model.UpdateTrackAvatarRequest) (string, []int) {
	var err []int
	data, ef := getFileData(req.File)
	if ef != nil {
		return "", append(err, ef...)
	}
	avatarId, e := db.UpdateTrackAvatar(req.UserId, req.TrackId, data)
	if e != nil {
		return "", append(err, code.DbError)
	}
	return avatarId, nil
}

func GetTrackHistory(r *model.TracksHistRequest) ([]model.TrackHistoryResponse, []int) {
	var err []int

	ping, e := db.LoadPingHistory(r)
	if e != nil {
		return nil, append(err, code.DbError)
	}
	sort.Slice(ping, func(i, j int) bool {
		return ping[i].EventTime.Before(ping[j].EventTime)
	})
	hist := make([]model.TrackHistoryResponse, 0)
	var prev model.PingDb
	for i, p := range ping {
		p1 := geo.NewPoint(float64(prev.Latitude), float64(prev.Longitude))
		p2 := geo.NewPoint(float64(p.Latitude), float64(p.Longitude))
		// фильтруем по расстоянию между точками
		dist := p1.GreatCircleDistance(p2) * 1000
		if dist >= 50 || i == len(ping)-1 { // 50 метров или последний
			h := model.TrackHistoryResponse{Date: p.EventTime, Latitude: p.Latitude, Longitude: p.Longitude}
			hist = append(hist, h)
		}
		prev = p
	}
	return hist, nil
}

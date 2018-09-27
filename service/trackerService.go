package service

import (
	"gobeacon/code"
	"gobeacon/model"
)

func GetTrackerById(id string) (interface{}, []int) {
	var err []int
	tracker, e := getTrackerById(id)
	if e != nil {
		err = append(err, code.DbError)
	}
	return tracker, err
}

func CreateTracker(req *model.TrackCreateRequest) (string, []int) {
	var err []int
	// ищем может уже добавлял кто-то (ситуация с 2мя родителями)
	id, e := existTrackByDevice(req.DeviceId)
	if e != nil {
		return "",  append(err, code.DbError)
	}
	if id == nil {// если ничего нет, то добавляем трекер
		if id, e = insertNewTrack(req); e != nil {
			return "",  append(err, code.DbError)
		}
	}
	// ищем текущую связь стрекером, на случай если решили повторно добавить (
	exist, e := existTrackPref(req.UserId, id.(string))
	if e != nil {
		return "",  append(err, code.DbError)
	}
	if exist {// если связь уже есть отпинываем
		return "",  append(err, code.TrackForUserExist)
	}
	if e = insertNewTrackPref(id.(string), req); e != nil {
		return "",  append(err, code.DbError)
	}
	return id.(string), err
}

func GetAllTrackersForUser(userId string) (interface{}, []int) {
	var err []int
	trackerList, e := getTrackersByUserId(userId)
	if e != nil {
		err = append(err, code.DbErrorGetTracker)
	}
	return trackerList, err
}

func GetTrackersByIds(ids []string) (interface{}, []int) {
	var err []int
	trackerList, e := getTrackerByIds(ids)
	if e != nil {
		err = append(err, code.DbErrorGetTracker)
	}
	return trackerList, err
}

func DeleteTrack(userId string, trackId string) ([]int) {
	var err []int
	e := deleteTrackForUser(userId, trackId)
	if e != nil {
		err = append(err, code.DbError)
	}
	return err
}

func UpdateTrackerName(req *model.TracksNameRequest) ([]int) {
	var err []int
	e := updateTrackerName(req)
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
	avatarId, e := updateTrackAvatar(req.UserId, req.TrackId, data)
	if e != nil {
		return "", append(err, code.DbError)
	}
	return avatarId, nil
}

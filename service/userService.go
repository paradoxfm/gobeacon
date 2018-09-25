package service

import (
	"gobeacon/code"
	"gobeacon/model"
)

func UserGetProfile(r *model.GetProfileRequest) (model.ProfileResponse, []int) {
	var err []int
	usr, e := getUserById(r.UserId)
	if e != nil {
		err = append(err, code.UserWithEmailNotFound) //пользователь не найден
		return model.ProfileResponse{}, err
	}
	rez := model.ProfileResponse{Id: usr.Id.String(), Email: usr.Email, Avatar: usr.Avatar}
	for id, tr := range usr.Trackers {
		rez.Trackers = append(rez.Trackers, model.UserTracker{Id: id.String(), Avatar: tr.Avatar, Name: tr.Name})
	}

	return rez, err
}

func UserUpdatePushId(r *model.UpdatePushRequest) (interface{}, []int) {
	var err []int
	e := updateUserPushId(r)
	if e != nil {
		err = append(err, code.DbErrorUpdateUserPush)
	}
	return nil, err
}

func ChangePassword(r *model.ChangePasswordRequest) ([]int) {
	var err []int
	userDb, e := getUserById(r.UserId)
	if e != nil {
		return append(err, code.DbError)
	}
	if !checkHash(r.OldPassword, userDb.Password) {
		return append(err, code.InavlidCurrentPasswords)
	}
	hash, _ := hashPassword(r.NewPassword)
	e = updateUserPassword(r.UserId, hash)
	if e != nil {
		return append(err, code.DbError)
	}
	return err
}

func SaveHeartbeat(p *model.Heartbeat) (*model.Tracker, []int) {
	t, e := getTrackerIdByDevice(p.DeviceId)
	var err []int
	if e != nil {
		return &t, append(err, code.DbError)
	}
	pingDb := model.PingDb{
		TrackerId:    t.Id,
		EventTime:    p.DateTime,
		BatteryPower: float32(p.Power),
		Latitude:     p.Latitude,
		Longitude:    p.Longitude,
		//ZoneId:       nil,
		SignalSource: getSignalId(p),
	}

	e = insertPing(&pingDb)
	if e != nil {
		return nil, append(err, code.DbError)
	}
	t.LatitudeLast = pingDb.Latitude
	t.LongitudeLast = pingDb.Longitude
	t.BatteryPowerLast = pingDb.BatteryPower
	return &t, nil
}

func getSignalId(p *model.Heartbeat) int {
	if p.IsGsm {
		return 1
	} else if p.IsWifi {
		return 3
	} else if p.IsGps {
		return 2
	}
	return 0
}

func CheckAndUpdateTracker(trk *model.Tracker) {
	tracker, e := getTrackerById(trk.Id.String())
	if e != nil {
		return
	}
	updateLastTracker(trk)
	alarmsCheck(&tracker, trk, false, false)
}

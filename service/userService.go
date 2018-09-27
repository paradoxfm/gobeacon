package service

import (
	"bytes"
	"encoding/base64"
	"gobeacon/code"
	"gobeacon/model"
	"image/jpeg"
	"io"
	"mime/multipart"
)

func UserGetProfile(r *model.GetProfileRequest) (model.ProfileResponse, []int) {
	var err []int
	usr, e := getUserById(r.UserId)
	if e != nil {
		return model.ProfileResponse{}, append(err, code.UserWithEmailNotFound) //пользователь не найден
	}
	rez := model.ProfileResponse{Id: usr.Id.String(), Email: usr.Email, Avatar: usr.Avatar}
	prefs, e := getTrackPrefs(r.UserId)
	if e != nil {
		return model.ProfileResponse{}, append(err, code.DbError)
	}
	for _, tr := range prefs {
		rez.Trackers = append(rez.Trackers, model.UserTracker{Id: tr.TrackId.String(), Avatar: tr.AvatarId, Name: tr.Name})
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

func GetAvatar(id string) (model.AvatarResponse, []int) {
	var err []int
	data, e := loadAvatar(id)
	if e != nil {
		return model.AvatarResponse{}, append(err, code.DbError)
	}
	strB64 := base64.StdEncoding.EncodeToString(data)
	return model.AvatarResponse{Data: strB64}, nil
}

func UpdateUserAvatar(req *model.UpdateAvatarRequest) (string, []int) {
	var err []int
	data, ef := getFileData(req.File)
	if ef != nil {
		return "", append(err, ef...)
	}
	avatarId, e := updateUserAvatar(req.UserId, data)
	if e != nil {
		return "", append(err, code.DbError)
	}
	return avatarId, nil
}

func getFileData(file *multipart.FileHeader) ([]byte, []int) {
	var err []int
	cont, ef := file.Open()
	if ef != nil {
		return nil, append(err, code.CantOpenFile)
	}
	buf := bytes.NewBuffer(nil)
	if _, e := io.Copy(buf, cont); e != nil {
		return nil, append(err, code.CantReadFile)
	}
	data := buf.Bytes()
	if rez, er := validateJpeg(data); !rez {
		return nil, append(err, er)
	}
	return data, nil
}

func validateJpeg(data []byte) (bool, int) {
	img, e := jpeg.Decode(bytes.NewReader(data))
	if e != nil {
		return false, code.InvalidImage
	}
	b := img.Bounds()
	if b.Dx() != 250 && b.Dy() != 250 {
		return false, code.InvalidImageSize
	}
	return true, -1
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

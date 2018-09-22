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
	rez := model.ProfileResponse{Id:usr.Id.String(), Email: usr.Email, Avatar: usr.Avatar}
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

package service

import (
	"gobeacon/code"
	"gobeacon/model"
)

func UserGetProfile(r *model.GetProfileRequest) (model.ProfileResponse, []int) {
	var err []int
	rez := model.ProfileResponse{}
	usr, e := getUserByEmail(r.UserId)
	if e != nil {
		err = append(err, code.UserWithEmailNotFound)//пользователь не найден
	} else {
		rez.Id = usr.Id.String()
		rez.Email = usr.Email
		rez.Avatar = usr.Avatar
		for id, tr := range usr.Trackers {
			rez.Trackers = append(rez.Trackers, model.UserTracker{Id: id.String(), Avatar: tr.Avatar, Name: tr.Name})
		}
	}

	return rez, err
}

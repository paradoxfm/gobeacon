package service

import "gobeacon/code"

func GetTrackerById(id string) (interface{}, []int) {
	var err []int
	tracker, e := getTrackerById(id)
	if e != nil {
		err = append(err, code.DbErrorGetTracker)
	}
	return tracker, err
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

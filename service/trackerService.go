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
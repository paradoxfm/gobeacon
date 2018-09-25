package service

import (
	"github.com/kellydunn/golang-geo"
	"gobeacon/model"
)

func alarmsCheck(prev *model.Tracker, curr *model.Tracker, lowPowerAlarm bool, sosAlarm bool) {
	var userPush map[string][]string
	for _, userId := range prev.Users {
		ids, e := getUserPushIds(userId)
		if e != nil && len(ids) > 0 {
			userPush[userId] = ids
		}
	}
	if len(userPush) == 0 {
		return
	}

	// LOW POWER ALARM
	if (curr.DeviceType == 1 && lowPowerAlarm) || (prev.BatteryPowerLast >= 20 && curr.BatteryPowerLast < 20) {
		data := map[string]interface{}{
			"message": "From iGurkin",
			"details": map[string]string{
				"name":  "Name",
				"user":  "Admin",
				"thing": "none",
			},
		}
		SendPushForUsers(userPush, data)
	}
	checkZones(prev, curr)
}

func checkZones(prev *model.Tracker, curr *model.Tracker) {
	if (prev.LatitudeLast == 0 && prev.LongitudeLast == 0) || (curr.LatitudeLast == 0 && curr.LongitudeLast == 0) {
		return
	}
	pOld := geo.NewPoint(float64(prev.LatitudeLast), float64(prev.LongitudeLast))
	pNew := geo.NewPoint(float64(curr.LatitudeLast), float64(curr.LongitudeLast))

	zones, err := getZonesByTrackId(curr.Id.String())
	if err != nil {
		return
	}
	for _, geoZone := range zones {
		var points []*geo.Point

		for _, gp := range geoZone.Points {
			points = append(points, geo.NewPoint(float64(gp.Latitude), float64(gp.Longitude)))
		}

		zone := geo.NewPolygon(points)

		//fmt.Printf("%#v \n", geoTokens)
		if zone.Contains(pOld) != zone.Contains(pNew) {
			ids, e := getUserPushIds(geoZone.UserId.String())
			if e != nil || len(ids) == 0 {
				break
			}
			userPush := map[string][]string{geoZone.UserId.String(): ids}
			//name, avatar := GetTrackerInfo(geoUser, trackerId)
			data := map[string]interface{}{
				"message": "From iGurkin",
				"details": map[string]string{
					"name":  "Name",
					"user":  "Admin",
					"thing": "none",
				},
			}
			SendPushForUsers(userPush, data)
		}
	}
}

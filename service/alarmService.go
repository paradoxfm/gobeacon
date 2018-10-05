package service

import (
	"github.com/kellydunn/golang-geo"
	"gobeacon/model"
)

const (
	lowPowerMsgId = 2001
	zoneMsgId     = 2002
)

type AlarmConf struct {
	UserId  string
	PushIds []string
	Zones   []model.GeoZoneDb
	Pref    model.TrackPref
}

func createPushData(trackId string) (map[string]AlarmConf) {
	rez := make(map[string]AlarmConf)
	trackPref, _ := getTrackPrefsByTrack(trackId)
	for _, tr := range trackPref {
		usrId := tr.UserId.String()
		ids, _ := getUserPushIds(usrId)
		if len(ids) != 0 {
			rez[usrId] = AlarmConf{UserId: usrId, PushIds: ids, Pref: tr}
		}
	}
	zones, _ := getZonesByTrackId(trackId)
	for _, z := range zones {
		usrId := z.UserId.String()
		alarmConf := rez[usrId]
		alarmConf.Zones = append(alarmConf.Zones, z)
	}
	return rez
}

func alarmsCheck(prev *model.Tracker, curr *model.Tracker, lowPowerAlarm bool, sosAlarm bool) {
	confList := createPushData(curr.Id.String())

	// LOW POWER ALARM
	if (curr.DeviceType == 1 && lowPowerAlarm) || (prev.BatteryPowerLast >= 20 && curr.BatteryPowerLast < 20) {
		data := map[string]interface{}{
			"message":      lowPowerMsgId,
			"tracker_id":   curr.Id.String(),
			"tracker_name": curr.Id.String(),
		}
		for userId, conf := range confList {
			SendPushForUser(userId, conf.PushIds, data)
		}
	}
	checkZones(prev, curr, confList)
}

func checkZones(prev *model.Tracker, curr *model.Tracker, confList map[string]AlarmConf) {
	if (prev.LatitudeLast == 0 && prev.LongitudeLast == 0) || (curr.LatitudeLast == 0 && curr.LongitudeLast == 0) {
		return
	}
	pOld := geo.NewPoint(float64(prev.LatitudeLast), float64(prev.LongitudeLast))
	pNew := geo.NewPoint(float64(curr.LatitudeLast), float64(curr.LongitudeLast))

	for _, conf := range confList {
		checkZonesForUser(pOld, pNew, conf)
	}
}

func checkZonesForUser(pOld *geo.Point, pNew *geo.Point, conf AlarmConf) {
	for _, geoZone := range conf.Zones {
		var points []*geo.Point

		for _, gp := range geoZone.Points {
			points = append(points, geo.NewPoint(float64(gp.Latitude), float64(gp.Longitude)))
		}

		zone := geo.NewPolygon(points)

		if zone.Contains(pOld) != zone.Contains(pNew) {
			data := map[string]interface{}{
				"message":      zoneMsgId,
				"tracker_id":   conf.Pref.TrackId.String(),
				"tracker_name": conf.Pref.Name,
				"zone_name":    geoZone.Name,
			}
			SendPushForUser(conf.UserId, conf.PushIds, data)
			return
		}
	}
}

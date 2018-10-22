package service

import (
	"github.com/kellydunn/golang-geo"
	"gobeacon/db"
	"gobeacon/model"
)

const (
	lowPowerMsgId = 2001
	zoneMsgIn     = 2002
	zoneMsgOut    = 2003
)

type AlarmConf struct {
	UserId    string
	TrackName string
	PushIds   []string
	Zones     []model.GeoZoneDb
	Pref      model.TrackPref
}

func createPushData(trackId string) (map[string]AlarmConf) {
	rez := make(map[string]AlarmConf)
	trackPref, _ := db.GetTrackPrefsByTrack(trackId)
	zones, _ := db.LoadZonesByTrackId(trackId)
	for _, tr := range trackPref {
		usrId := tr.UserId.String()
		trName := tr.Name
		zoneList := make([]model.GeoZoneDb, 0)
		for _, zn := range zones {
			if zn.UserId.String() == usrId {
				zoneList = append(zoneList, zn)
			}
		}
		ids, _ := db.LoadUserPushIds(usrId)
		if len(ids) != 0 {
			rez[usrId] = AlarmConf{UserId: usrId, TrackName: trName, PushIds: ids, Pref: tr, Zones: zoneList}
		}
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
			data["tracker_name"] = conf.TrackName
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

		inZoneNew := zone.Contains(pNew)
		if zone.Contains(pOld) != inZoneNew {
			msgId := zoneMsgOut
			if inZoneNew {
				msgId = zoneMsgIn
			}
			data := map[string]interface{}{
				"message":      msgId,
				"tracker_id":   conf.Pref.TrackId.String(),
				"tracker_name": conf.TrackName,
				"zone_name":    geoZone.Name,
			}
			SendPushForUser(conf.UserId, conf.PushIds, data)
			return
		}
	}
}

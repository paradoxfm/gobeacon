package db

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"gobeacon/model"
	"time"
)

func GetTrackPrefsByUser(userId string) ([]model.TrackPref, error) {
	stmt, names := qb.Select(tTrackPref).Where(qb.Eq("user_id")).ToCql()
	var pref []model.TrackPref

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId})
	err := q.SelectRelease(&pref)
	return pref, err
}

func GetTrackPrefsByTrack(trackId string) ([]model.TrackPref, error) {
	stmt, names := qb.Select(tTrackPref).Columns("user_id", "track_id", "track_name").Where(qb.Eq("track_id")).ToCql()
	var pref []model.TrackPref

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"track_id": trackId})
	err := q.SelectRelease(&pref)
	return pref, err
}

func GetTrackerById(id string) (model.Tracker, error) {
	builder := qb.Select(tTrackers)
	stmt, names := builder.Where(qb.Eq("id")).ToCql()

	var track model.Tracker
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": id})
	err := q.GetRelease(&track)
	return track, err
}

func GetTrackerByIds(ids []string) ([]model.Tracker, error) {
	builder := qb.Select(tTrackers)
	stmt, names := builder.Where(qb.In("id")).ToCql()

	var track []model.Tracker
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": ids})
	err := q.SelectRelease(&track)
	return track, err
}

func GetTrackersByUserId(userId string) ([]model.Tracker, error) {
	prefs, er := GetTrackPrefsByUser(userId)
	if er != nil {
		return nil, er
	}
	var trackers []model.Tracker
	if len(prefs) == 0 {
		return trackers, nil
	}
	var ids []string
	for _, v := range prefs {
		ids = append(ids, v.TrackId.String())
	}
	stmt, names := qb.Select(tTrackers).Where(qb.In("id")) /*.AllowFiltering()*/ .ToCql()

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": ids})
	err := q.SelectRelease(&trackers)
	return trackers, err
}

func ExistTrackByDevice(deviceId string) (interface{}, error) {
	stmt, names := qb.Select(tTrackers).Columns("id").Where(qb.Eq("device_id")).AllowFiltering().ToCql()

	var t []model.Tracker
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"device_id": deviceId})
	err := q.SelectRelease(&t)
	if len(t) > 0 {
		return t[0].Id.String(), err
	}
	return nil, err
}

func ExistTrackPref(userId string, trackId string) (bool, error) {
	stmt, names := qb.Select(tTrackPref).Columns("user_id", "track_id").Where(qb.Eq("user_id"), qb.Eq("track_id")).AllowFiltering().ToCql()

	var t []model.TrackPref
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId, "track_id": trackId})
	err := q.SelectRelease(&t)
	return len(t) > 0, err
}

func InsertNewTrackPref(trackId string, req *model.TrackCreateRequest) (error) {
	stmt, names := qb.Insert(tTrackPref).Columns("user_id", "track_id", "track_name").ToCql()

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": req.UserId, "track_id": trackId, "track_name": req.Name})
	return q.ExecRelease()
}

func InsertNewTrack(t *model.TrackCreateRequest) (interface{}, error) {
	stmt, names := qb.Insert(tTrackers).Columns("id", "device_id", "created_at").ToCql()
	uuid, _ := gocql.RandomUUID()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": uuid, "device_id": t.DeviceId, "created_at": time.Now()})
	err := q.ExecRelease()
	return uuid.String(), err
}

func GetTrackerIdByDevice(deviceId string) (model.Tracker, error) {
	stmt, names := qb.Select(tTrackers).Columns("id", "device_id", "latitude_last", "longitude_last", "battery_power_last").Where(qb.Eq("device_id")).AllowFiltering().ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"device_id": deviceId})
	var track model.Tracker
	err := q.GetRelease(&track)
	return track, err
}

func GetTrackUserIds(trackId string) ([]string, error) {
	stmt, names := qb.Select(tTrackPref).Columns("user_id").Where(qb.Eq("track_id")).ToCql()

	var ids []string

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"track_id": trackId})
	err := q.SelectRelease(&ids)
	return ids, err
}

func GetTrackPrefForUser(userId string, trackId string) (model.TrackPref, error) {
	stmt, names := qb.Select(tTrackPref).Where(qb.Eq("user_id"), qb.Eq("track_id")).Limit(1).ToCql()

	var p model.TrackPref

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId, "track_id": trackId})
	err := q.GetRelease(&p)
	return p, err
}

func UpdateTrackAvatar(userId string, trackId string, blob []byte) (string, error) {
	trackPref, _ := GetTrackPrefForUser(userId, trackId)

	avatarLink := trackPref.AvatarId
	avaId, e := insertNewAvatar(blob)
	if e != nil {
		return "", e
	}

	// обновляем ссылку на аватар
	stmt, names := qb.Update(tTrackPref).Set("track_ava").Where(qb.Eq("user_id"), qb.Eq("track_id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId, "track_id": trackId, "track_ava": avaId})
	if err := q.ExecRelease(); err != nil {
		return "", err
	}

	if len(avatarLink) > 0 {
		if e = deleteAvatar(avatarLink); e != nil {
			return "", e
		}
	}

	return avaId, nil
}

func DeleteTrackForUser(userId string, trackId string) (error) {
	stmt, names := qb.Delete(tTrackPref).Where(qb.Eq("user_id"), qb.Eq("track_id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId, "track_id": trackId})
	if e := q.ExecRelease(); e != nil {
		return e
	}
	strings, e := GetTrackUserIds(trackId)
	if e != nil {
		return e
	}
	if len(strings) > 0 { // если остались еще связи, оставляем трекер
		return nil
	}
	stmt, names = qb.Delete(tTrackers).Where(qb.Eq("id")).ToCql()
	q = gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": trackId})
	if e := q.ExecRelease(); e != nil {
		return e
	}
	stmt, names = qb.Delete(tPings).Where(qb.Eq("tracker_id")).ToCql()
	q = gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"tracker_id": trackId})
	if e := q.ExecRelease(); e != nil {
		return e
	}
	// отвязать от зоны
	stmt, names = qb.Update(tZones).Remove("trackers").Set("updated_at").Where(qb.Eq("user_id")).ToCql()
	q = gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId, "trackers": trackId, "updated_at": time.Now()})
	q.ExecRelease()
	return nil
}

func UpdateTrackPref(req *model.TrackPrefRequest) (error) {
	stmt, names := qb.Update(tTrackPref).Set("track_name", "track_offs").Where(qb.Eq("user_id"), qb.Eq("track_id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"track_name": req.Name, "track_offs": req.Offset, "user_id": req.UserId, "track_id": req.TrackId})
	err := q.ExecRelease()
	return err
}

func UpdateLastTracker(tr *model.Tracker, dt time.Time) (error) {
	stmt, names := qb.Update(tTrackers).Set("latitude_last", "longitude_last", "battery_power_last", "updated_at", "signal_timestamp_last").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"latitude_last": tr.LatitudeLast, "longitude_last": tr.LongitudeLast, "battery_power_last": tr.BatteryPowerLast, "updated_at": time.Now(), "id": tr.Id.String(), "signal_timestamp_last": dt})
	err := q.ExecRelease()
	return err
}

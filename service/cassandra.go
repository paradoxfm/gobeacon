package service

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"gobeacon/model"
	"log"
	"time"
)

var session *gocql.Session

var tUsers = "watch.users"
var tTrackPref = "watch.user_track_prefs"
var tAvatars = "watch.files"
var tTrackers = "watch.trackers"
var tZones = "watch.geozones"
var tPings = "track_ping"

func init() {
	var err error

	cluster := gocql.NewCluster(Config().CassandraIp)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: Config().CassandraUser,
		Password: Config().CassandraPassword,
	}
	cluster.Keyspace = Config().CassandraKey
	cluster.ProtoVersion = 4
	cluster.ReconnectInterval = 10 * time.Second
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Cassandra init done")
}

func insertNewUser(email string, password string) (error) {
	stmt, _ := qb.Insert(tUsers).Columns("id", "email", "password", "created_at").ToCql()
	q := session.Query(stmt)
	err := q.Bind(gocql.TimeUUID(), email, password, time.Now()).Exec()

	return err
}

func getUserByEmail(email string) (model.UserDb, error) {
	stmt, names := qb.Select(tUsers).Columns("id", "email", "password").Where(qb.Eq("email")).Limit(1).ToCql()
	var u model.UserDb

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"email": email,})
	err := q.GetRelease(&u)
	return u, err
}

func getUserById(id string) (model.UserDb, error) {
	stmt, names := qb.Select(tUsers).Columns("id", "email", "password", "avatar").Where(qb.Eq("id")).Limit(1).ToCql()
	var u model.UserDb

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": id,})
	err := q.GetRelease(&u)
	return u, err
}

func getTrackPrefs(userId string) ([]model.TrackPref, error) {
	stmt, names := qb.Select(tTrackPref).Where(qb.Eq("user_id")).ToCql()
	var pref []model.TrackPref

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId})
	err := q.SelectRelease(&pref)
	return pref, err
}

func updateUserPushId(r *model.UpdatePushRequest) (error) {
	stmt, names := qb.Update(tUsers).Add("push_id").Set("updated_at").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": r.UserId, "push_id": []string{r.PushId}, "updated_at": time.Now()})
	err := q.ExecRelease()
	return err
}

func updateUserPassword(userId string, hash string) (error) {
	stmt, names := qb.Update(tUsers).Set("password").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": userId, "password": hash})

	err := q.ExecRelease()
	return err
}

func loadAvatar(id string) ([]byte, error) {
	stmt, names := qb.Select(tAvatars).Where(qb.Eq("id")).Limit(1).ToCql()
	var u model.BlobDb

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": id,})
	err := q.GetRelease(&u)
	return u.Content, err
}

func updateUserAvatar(userId string, blob []byte) (string, error) {
	userDb, _ := getUserById(userId)

	avatarLink := userDb.Avatar
	avaId, e := insertNewAvatar(blob)
	if e != nil {
		return "", e
	}

	// обновляем ссылку на аватар
	stmt, names := qb.Update(tUsers).Set("avatar").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": userId, "avatar": avaId})
	if err := q.ExecRelease(); err != nil {
		return "", err
	}

	if e = deleteAvatar(avatarLink); e != nil {
		return "", e
	}

	return avaId, nil
}

// удаляем старый аватар
func deleteAvatar(uid string) (error) {
	stmt, names := qb.Delete(tAvatars).Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": uid})
	err := q.ExecRelease()
	return err
}

// вставляем новый аватар
func insertNewAvatar(blob []byte) (string, error) {
	uuid, _ := gocql.RandomUUID()
	ava := model.BlobDb{Id: uuid, Content: blob}
	// вставляем новый аватар
	stmt, names := qb.Insert(tAvatars).Columns("id", "avatar").ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindStruct(&ava)
	if err := q.ExecRelease(); err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func getUserPushIds(userId string) ([]string, error) {
	stmt, names := qb.Select(tUsers).Columns("id", "push_id").Where(qb.Eq("id")).Limit(1).ToCql()
	var u model.UserDb

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": userId,})
	err := q.GetRelease(&u)
	return u.PushId, err
}

func removeInvalidPush(userId string, push []string) error {
	stmt, names := qb.Update(tUsers).Remove("push_id").Set("updated_at").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": userId, "push_id": push, "updated_at": time.Now()})
	err := q.ExecRelease()
	return err
}

func getTrackerById(id string) (model.Tracker, error) {
	builder := qb.Select(tTrackers)
	stmt, names := builder.Where(qb.Eq("id")).ToCql()

	var track model.Tracker
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": id})
	err := q.GetRelease(&track)
	return track, err
}

func getTrackerByIds(ids []string) ([]model.Tracker, error) {
	builder := qb.Select(tTrackers)
	stmt, names := builder.Where(qb.In("id")).ToCql()

	var track []model.Tracker
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": ids})
	err := q.SelectRelease(&track)
	return track, err
}

func getTrackersByUserId(userId string) ([]model.Tracker, error) {
	prefs, er := getTrackPrefs(userId)
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

func existTrackByDevice(deviceId string) (interface{}, error) {
	stmt, names := qb.Select(tTrackers).Columns("id").Where(qb.Eq("device_id")).AllowFiltering().ToCql()

	var t []model.Tracker
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"device_id": deviceId})
	err := q.SelectRelease(&t)
	if len(t) > 0 {
		return t[0].Id.String(), err
	}
	return nil, err
}

func existTrackPref(userId string, trackId string) (bool, error) {
	stmt, names := qb.Select(tTrackPref).Columns("user_id", "track_id").Where(qb.Eq("user_id"), qb.Eq("track_id")).AllowFiltering().ToCql()

	var t []model.TrackPref
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId, "track_id": trackId})
	err := q.SelectRelease(&t)
	return len(t) > 0, err
}

func insertNewTrackPref(trackId string, req *model.TrackCreateRequest) (error) {
	stmt, names := qb.Insert(tTrackPref).Columns("user_id", "track_id", "track_name").ToCql()

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": req.UserId, "track_id": trackId, "track_name": req.Name})
	return q.ExecRelease()
}

func insertNewTrack(t *model.TrackCreateRequest) (interface{}, error) {
	stmt, names := qb.Insert(tTrackers).Columns("id", "device_id", "created_at").ToCql()
	uuid, _ := gocql.RandomUUID()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": uuid, "device_id": t.DeviceId, "created_at": time.Now()})
	err := q.ExecRelease()
	return uuid.String(), err
}

func getTrackerIdByDevice(deviceId string) (model.Tracker, error) {
	stmt, names := qb.Select(tTrackers).Columns("id", "device_id").Where(qb.Eq("device_id")).AllowFiltering().ToCql()

	var track model.Tracker
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"device_id": deviceId})
	err := q.GetRelease(&track)
	return track, err
}

func getTrackUserIds(trackId string) ([]string, error) {
	stmt, names := qb.Select(tTrackPref).Columns("user_id").Where(qb.Eq("track_id")).ToCql()

	var ids []string

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"track_id": trackId})
	err := q.SelectRelease(&ids)
	return ids, err
}

func getTrackPrefForUser(userId string, trackId string) (model.TrackPref, error) {
	stmt, names := qb.Select(tTrackPref).Where(qb.Eq("user_id"), qb.Eq("track_id")).Limit(1).ToCql()

	var p model.TrackPref

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId, "track_id": trackId})
	err := q.GetRelease(&p)
	return p, err
}

func updateTrackAvatar(userId string, trackId string, blob []byte) (string, error) {
	trackPref, _ := getTrackPrefForUser(userId, trackId)

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

	if e = deleteAvatar(avatarLink); e != nil {
		return "", e
	}

	return avaId, nil
}

func deleteTrackForUser(userId string, trackId string) (error) {
	stmt, names := qb.Delete(tTrackPref).Where(qb.Eq("user_id"), qb.Eq("track_id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId, "track_id": trackId})
	if e := q.ExecRelease(); e != nil {
		return e
	}
	strings, e := getTrackUserIds(trackId)
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

func updateTrackerName(req *model.TracksNameRequest) (error) {
	/*u, e := getUserById(req.UserId)
	if e != nil {
		return e
	}
	u.Trackers
	stmt, names := qb.Update(tUsers).Set("latitude_last", "longitude_last", "battery_power_last", "updated_at").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"latitude_last": tr.LatitudeLast, "longitude_last": tr.LongitudeLast, "battery_power_last": tr.BatteryPowerLast, "updated_at": time.Now(), "id": tr.Id.String()})
	err := q.ExecRelease()
	return err*/
	return nil
}

func updateLastTracker(tr *model.Tracker, dt time.Time) (error) {
	stmt, names := qb.Update(tTrackers).Set("latitude_last", "longitude_last", "battery_power_last", "updated_at", "signal_timestamp_last").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"latitude_last": tr.LatitudeLast, "longitude_last": tr.LongitudeLast, "battery_power_last": tr.BatteryPowerLast, "updated_at": time.Now(), "id": tr.Id.String(), "signal_timestamp_last": dt})
	err := q.ExecRelease()
	return err
}

func insertPing(ping *model.PingDb) error {
	stmt, names := qb.Insert(tPings).Columns("tracker_id", "event_time", "battery_power", "latitude", "longitude", "signal_source").ToCql()
	err := gocqlx.Query(session.Query(stmt), names).BindStruct(&ping).ExecRelease()
	return err
}

func getPingOnInterval(trackId string, tFrom time.Time, tTo time.Time) ([]model.PingDb, error) {
	stmt, names := qb.Select(tPings).Columns("event_time", "latitude", "longitude").Where(qb.EqNamed("tracker_id", "id"), qb.GtNamed("event_time", "tfrom"), qb.LtNamed("event_time", "tto")).AllowFiltering().ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": trackId, "tfrom": tFrom, "tto": tTo})

	var list []model.PingDb
	err := q.Select(&list)
	return list, err
}

func getAllZoneByUserId(userId string) ([]model.GeoZoneDb, error) {
	stmt, names := qb.Select(tZones).Where(qb.Eq("user_id")).ToCql()

	var zones []model.GeoZoneDb
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId})
	err := q.SelectRelease(&zones)
	return zones, err
}

func createZoneForUser(r *model.ZoneCreateRequest) (model.GeoZoneDb, error) {
	id, _ := gocql.RandomUUID()
	usr, _ := gocql.ParseUUID(r.UserId)
	db := model.GeoZoneDb{Id: id, UserId: usr, Name: r.Name, CreatedAt: time.Now(), UpdatedAt: time.Now(), Points: r.Points, Trackers: make(map[gocql.UUID]bool)}

	stmt, names := qb.Insert(tZones).Columns("id", "user_id", "name", "created_at", "updated_at", "points", "trackers").ToCql()
	e := gocqlx.Query(session.Query(stmt), names).BindStruct(&db).ExecRelease()
	return db, e
}

func updateZoneProp(r *model.ZoneCreateRequest) (error) {
	stmt, names := qb.Update(tZones).Set("name", "points").Where(qb.Eq("id")).ToCql()
	e := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": r.Id, "name": r.Name, "points": r.Points}).ExecRelease()
	return e
}

func findZoneById(id string) (model.GeoZoneDb, error) {
	stmt, names := qb.Select(tZones).Where(qb.Eq("id")).ToCql()
	var zone model.GeoZoneDb
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": id})
	err := q.GetRelease(&zone)
	return zone, err
}

func getZonesByTrackId(trackId string) ([]model.GeoZoneDb, error) {
	stmt, names := qb.Select(tZones).Where(qb.ContainsKey("trackers")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"trackers": trackId})

	var zones []model.GeoZoneDb
	err := q.SelectRelease(&zones)
	return zones, err
}

func deleteZoneById(zoneId string) (error) {
	stmt, names := qb.Delete(tZones).Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": zoneId})
	return q.ExecRelease()
}

func updateZoneTrackers(zoneId string, track map[string]bool) (error) {
	stmt, names := qb.Update(tZones).Set("trackers").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"trackers": track, "id": zoneId})

	return q.ExecRelease()
}

func loadTrackHistory(r *model.TracksHistRequest) ([]model.PingDb, error) {
	stmt, names := qb.Select(tPings).Where(qb.Eq("tracker_id"), qb.GtNamed("event_time", "from"), qb.LtOrEqNamed("event_time", "to")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"tracker_id": r.TrackId, "from": r.DateFrom, "to": r.DateTo})

	var ping []model.PingDb
	err := q.SelectRelease(&ping)
	return ping, err
}

/*func Exception(err error) {
	Error(err.Error())
}*/

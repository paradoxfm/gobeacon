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
	stmt, names := qb.Select(tUsers).Columns("id", "email", "password", "trackers").Where(qb.Eq("email")).Limit(1).ToCql()
	var u model.UserDb

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"email": email,})
	err := q.GetRelease(&u)
	return u, err
}

func getUserById(id string) (model.UserDb, error) {
	stmt, names := qb.Select(tUsers).Columns("id", "email", "password", "avatar", "trackers").Where(qb.Eq("id")).Limit(1).ToCql()
	var u model.UserDb

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": id,})
	err := q.GetRelease(&u)
	return u, err
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
	builder := qb.Select(tTrackers)
	stmt, names := builder.Where(qb.Contains("users")).AllowFiltering().ToCql()

	var trackers []model.Tracker
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"users": userId})
	err := q.SelectRelease(&trackers)
	return trackers, err
}

func getTrackerIdByDevice(deviceId string) (model.Tracker, error) {
	builder := qb.Select(tTrackers).Columns("id", "device_id")
	stmt, names := builder.Where(qb.Eq("device_id")).AllowFiltering().ToCql()

	var track model.Tracker
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"device_id": deviceId})
	err := q.GetRelease(&track)
	return track, err
}

func getAllUsers() ([]model.UserDb) {
	stmt, names := qb.Select(tUsers).Columns("id", "trackers").ToCql()
	var users []model.UserDb
	q := gocqlx.Query(session.Query(stmt), names)

	e := q.SelectRelease(&users)
	if e != nil {
		return nil
	}
	return users
}

func insertTrackSettings(uuid gocql.UUID, trkId gocql.UUID, trkName string) (error) {
	stmt, names := qb.Insert(tTrackPref).ToCql()
	pref := model.TrackPref{UserId: uuid, TrackId: trkId, Name: trkName}
	err := gocqlx.Query(session.Query(stmt), names).BindStruct(&pref).ExecRelease()
	return err
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

func updateLastTracker(tr *model.Tracker) (error) {
	stmt, names := qb.Update(tTrackers).Set("latitude_last", "longitude_last", "battery_power_last", "updated_at").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"latitude_last": tr.LatitudeLast, "longitude_last": tr.LongitudeLast, "battery_power_last": tr.BatteryPowerLast, "updated_at": time.Now(), "id": tr.Id.String()})
	err := q.ExecRelease()
	return err
}

func insertPing(ping *model.PingDb) error {
	stmt, names := qb.Insert(tPings).ToCql()
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

func getZonesByTrackId(trackId string) ([]model.GeoZoneDb, error) {
	stmt, names := qb.Select(tZones).Where(qb.ContainsKey("trackers")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"trackers": trackId})

	var zones []model.GeoZoneDb
	err := q.SelectRelease(&zones)
	return zones, err
}

/*func Exception(err error) {
	Error(err.Error())
}*/

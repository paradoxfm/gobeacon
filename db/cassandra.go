package db

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
var tPings = "watch.track_ping"

func init() {
	conf := Config()
	cluster := gocql.NewCluster(conf.CassandraIp)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: conf.CassandraUser,
		Password: conf.CassandraPassword,
	}
	cluster.Keyspace = conf.CassandraKey
	cluster.ProtoVersion = 4
	cluster.ReconnectInterval = 10 * time.Second
	//go createSession(cluster)
	createSession(cluster)
}

func createSession(cluster *gocql.ClusterConfig) {
	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Cassandra init done")
}

func LoadAvatar(id string) ([]byte, error) {
	stmt, names := qb.Select(tAvatars).Where(qb.Eq("id")).Limit(1).ToCql()
	var u model.BlobDb

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": id,})
	err := q.GetRelease(&u)
	return u.Content, err
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

func InsertPing(ping *model.PingDb) error {
	stmt, names := qb.Insert(tPings).Columns("tracker_id", "event_time", "battery_power", "latitude", "longitude", "signal_source").ToCql()
	err := gocqlx.Query(session.Query(stmt), names).BindStruct(&ping).ExecRelease()
	return err
}

func LoadPingHistory(r *model.TracksHistRequest) ([]model.PingDb, error) {
	stmt, names := qb.Select(tPings).Where(qb.Eq("tracker_id"), qb.GtNamed("event_time", "from"), qb.LtOrEqNamed("event_time", "to")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"tracker_id": r.TrackId, "from": r.DateFrom, "to": r.DateTo})

	var ping []model.PingDb
	err := q.SelectRelease(&ping)
	return ping, err
}

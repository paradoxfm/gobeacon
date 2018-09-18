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
var tTrackers = "watch.trackers"

func init() {
	var err error

	cluster := gocql.NewCluster(Config().Cassandra_ip)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: Config().Cassandra_user,
		Password: Config().Cassandra_password,
	}
	cluster.Keyspace = Config().Cassandra_keyspace
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

func updateUserPassword(userId gocql.UUID, password string) (error) {
	stmt, names := qb.Update(tUsers).Set("password").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": userId, "password": password})

	err := q.Exec()
	return err
}

func getTrackersByUserId(userId gocql.UUID) ([]model.Tracker, error) {
	stmt, names := qb.Select(tTrackers)/*.Where(qb.In("users")).AllowFiltering()*/.ToCql()

	var trackers []model.Tracker
	q := gocqlx.Query(session.Query(stmt), names)/*.BindMap(qb.M{"users": []gocql.UUID {userId},})*/
	err := q.GetRelease(&trackers)
	return trackers, err
}

func getTrackerhistory(timeFrom time.Time, timeTo time.Time, trackerId string) ([]interface{}, error) {
	stmt, names := qb.Select("watch.ping").Where(qb.GtNamed("event_time", "from"), qb.LtNamed("event_time", "to"), qb.EqNamed("tracker_id", "id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"from": timeFrom, "to": timeTo, "id": trackerId})

	var list []interface{}
	err := q.SelectRelease(&list)
	return list, err
}

func addTrackerhistory(ping *model.PingDb) {
	qb.Insert("watch.ping").Columns("event_time").Timestamp(time.Now())
}

/*func Exception(err error) {
	Error(err.Error())
}*/

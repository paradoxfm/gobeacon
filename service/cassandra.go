package service

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/qb"
	"log"
	"gobeacon/model"
)

var session *gocql.Session

func init() {
	var err error
	if err == nil {
		return
	}

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

func insertNewUser(u model.UserDb) {
	stmt, _ := qb.Insert("watch.users").Columns("id", "email", "password", "created_at").ToCql()
	q := session.Query(stmt)
	err := q.Bind(u.Id, u.Email, u.Password).Exec()
	if err != nil {

	}
}

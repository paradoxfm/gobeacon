package model

import (
	"github.com/gocql/gocql"
	"time"
)

/*> CREATE TABLE users (
	id UUID,
	email varchar,
	password text,
	trackers map<UUID, frozen<tuple<varchar, varchar>>>,
	geozones set<UUID>,
	avatar varchar,
	push_id set<text>,
	created_at timestamp,
	updated_at timestamp,
	PRIMARY KEY (id)
);*/
type UserDb struct {
	Id       gocql.UUID `cql:"id" json:"id"`
	Email    string     `cql:"email" json:"email"`
	Password string     `cql:"password" json:"password"`
	Avatar   string     `cql:"avatar" json:"avatar"`
	PushId   []string   `cql:"push_id" json:"push_id"`

	//Trackers map[gocql.UUID]UserTrackers `cql:"trackers" json:"trackers"`
	//ZoneList []gocql.UUID                `cql:"geozones" json:"geozones"`
}

/*> CREATE TABLE watch.geozones (
	id UUID,
	trackers map<UUID, boolean>,
	user_id UUID,
	name varchar,
	points list<frozen<tuple<float, float>>>,
	created_at timestamp,
	updated_at timestamp,
	PRIMARY KEY (id)
);*/
type GeoZoneDb struct {
	Id      gocql.UUID  `cql:"id" json:"id" description:"Id геозоны"`
	UserId  gocql.UUID  `cql:"user_id" json:"user_id" description:"Id пользователя"`
	Name    string      `cql:"name" json:"name" description:"Имя геозоны"`
	Created time.Time   `cql:"created_at" json:"createTime" description:"Дата создания"`
	Points  []ZonePoint `cql:"points" json:"points" description:"Точки полигона"`
}

type ZonePoint struct {
	Latitude  float32 `json:"latitude,required" description:"Широта"`
	Longitude float32 `json:"longitude,required" description:"Долгота"`
}

/*> CREATE TABLE trackers (
	id UUID,
	device_id bigint,
	imei bigint,
	device_type int,
	signal_source int,
	latitude_last float,
	longitude_last float,
	battery_power_last float,
	signal_timestamp_last timestamp,
	users set<UUID>,
	created_at timestamp,
	updated_at timestamp,
	PRIMARY KEY (id)
);*/

/*> CREATE TABLE avatars (
	id UUID,
	content blob,
	PRIMARY KEY (id)
);*/
type BlobDb struct {
	id      gocql.UUID `cql:"id" json:"id"`
	content []byte     `cql:"content" json:"content"`
}

/*> CREATE TABLE watch.ping (
	event_time timestamp,
	tracker_id UUID,
	d bigint,
	PRIMARY KEY (time_ins, tracker_id)
);*/
type PingDb struct {
	EventTime time.Time `cql:"id" json:"id"`
}

type UserTrackers struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

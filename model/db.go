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
	Id       gocql.UUID `db:"id" json:"id"`
	Email    string     `db:"email" json:"email"`
	Password string     `db:"password" json:"password"`
	Avatar   string     `db:"avatar" json:"avatar"`
	PushId   []string   `db:"push_id" json:"push_id"`
	//Trackers map[gocql.UUID]UserTrackers `db:"trackers" json:"trackers"`

	//Trackers map[gocql.UUID]UserTrackers `cql:"trackers" json:"trackers"`
	//ZoneList []gocql.UUID                `cql:"geozones" json:"geozones"`
}

type UserTrackers struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
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
	Id        gocql.UUID          `db:"id" json:"id" description:"Id геозоны"`
	UserId    gocql.UUID          `db:"user_id" json:"user_id" description:"Id пользователя"`
	Name      string              `db:"name" json:"name" description:"Имя геозоны"`
	CreatedAt time.Time           `db:"created_at" json:"createTime" description:"Дата создания"`
	UpdatedAt time.Time           `db:"updated_at" json:"updateTime" description:"Дата обновления"`
	Points    []ZonePoint         `db:"points" json:"points" description:"Точки полигона"`
	Trackers  map[gocql.UUID]bool `db:"trackers" json:"points" description:"Список привязанных трекеров"`
}

type ZonePoint struct {
	Latitude  float32 `json:"latitude,required" description:"Широта"`
	Longitude float32 `json:"longitude,required" description:"Долгота"`
}

type Counter struct {
	Count int64
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
type Tracker struct {
	Id                  gocql.UUID `db:"id" json:"id"`
	DeviceId            string     `db:"device_id" json:"device_id"`
	Imei                string     `db:"imei" json:"-"`
	DeviceType          int        `db:"device_type" json:"device_type"`
	SignalSource        int        `db:"signal_source" json:"signal_source"`
	LatitudeLast        float32    `db:"latitude_last" json:"latitude_last"`
	LongitudeLast       float32    `db:"longitude_last" json:"longitude_last"`
	BatteryPowerLast    float32    `db:"battery_power_last" json:"battery_power_last"`
	Users               []string   `db:"users"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time  `db:"updated_at" json:"updated_at"`
	SignalTimestampLast time.Time  `db:"signal_timestamp_last" json:"signal_timestamp_last"`
}

func (trk Tracker) Copy() Tracker {
	rez := Tracker{Id: trk.Id, DeviceId: trk.DeviceId, DeviceType: trk.DeviceType, SignalSource: trk.SignalSource}
	rez.LatitudeLast = trk.LatitudeLast
	rez.LongitudeLast = trk.LongitudeLast
	rez.BatteryPowerLast = trk.BatteryPowerLast
	rez.Users = trk.Users
	rez.SignalTimestampLast = trk.SignalTimestampLast
	return rez
}

/*CREATE TABLE watch.user_track_prefs (
    user_id uuid,
    track_id uuid,
    track_ava uuid,
    track_name text,
    PRIMARY KEY (user_id, track_id)
)*/
type TrackPref struct {
	UserId   gocql.UUID `db:"user_id"`
	TrackId  gocql.UUID `db:"track_id"`
	AvatarId string     `db:"track_ava"`
	Name     string     `db:"track_name"`
	Offset   int        `db:"track_offs"`
}

/*> CREATE TABLE watch.files (
	id uuid PRIMARY KEY,
	link_id uuid,
	avatar blob
);*/
type BlobDb struct {
	Id      gocql.UUID `db:"id" json:"id"`
	Content []byte     `db:"avatar" json:"content"`
	//LinkId  gocql.UUID `db:"link_id" json:"id"`
}

/*> create table watch.track_ping (
	tracker_id uuid,
	event_time timestamp,
	battery_power float,
	latitude float,
	longitude float,
	zone_id uuid,
	signal_source int,
	primary key (tracker_id, event_time)
);*/
type PingDb struct {
	TrackerId    gocql.UUID `db:"tracker_id"`
	EventTime    time.Time  `db:"event_time" json:"datetime"`
	BatteryPower float32    `db:"battery_power"`
	Latitude     float32    `db:"latitude" json:"latitude"`
	Longitude    float32    `db:"longitude" json:"longitude"`
	SignalSource int        `db:"signal_source"`
	//ZoneId       gocql.UUID `db:"zone_id"`
}

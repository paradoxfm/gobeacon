package db

import (
	"github.com/gocql/gocql"
	"gobeacon/model"
	"time"
)

type TrackDbInterface interface {
	GetTrackPrefsByUser(userId string) ([]model.TrackPref, error)
	GetTrackPrefsByTrack(trackId string) ([]model.TrackPref, error)
	GetTrackerById(id string) (model.Tracker, error)
	GetTrackerByIds(ids []string) ([]model.Tracker, error)
	GetTrackersByUserId(userId string) ([]model.Tracker, error)
	ExistTrackByDevice(deviceId string) (interface{}, error)
	ExistTrackPref(userId string, trackId string) (bool, error)
	InsertNewTrackPref(trackId string, req *model.TrackCreateRequest) (error)
	InsertNewTrack(t *model.TrackCreateRequest) (interface{}, error)
	GetTrackerIdByDevice(deviceId string) (model.Tracker, error)
	GetTrackUserIds(trackId string) ([]string, error)
	GetTrackPrefForUser(userId string, trackId string) (model.TrackPref, error)
	UpdateTrackAvatar(userId string, trackId string, blob []byte) (string, error)
	DeleteTrackForUser(userId string, trackId string) (error)
	UpdateTrackPref(req *model.TrackPrefRequest) (error)
	UpdateLastTracker(tr *model.Tracker, dt time.Time) (error)
}

type TrackDataBase struct {
	session *gocql.Session
}

type ZoneDbInterface interface {
	LoadZonesByUserId(userId string) ([]model.GeoZoneDb, error)
	CreateZoneForUser(r *model.ZoneCreateRequest) (model.GeoZoneDb, error)
	UpdateZone(r *model.ZoneCreateRequest) (error)
	LoadZoneById(id string) (model.GeoZoneDb, error)
	LoadZonesByTrackId(trackId string) ([]model.GeoZoneDb, error)
	DeleteZoneById(zoneId string) (error)
	UpdateZoneTrackers(zoneId string, track map[string]bool) (error)
}

type ZoneDataBase struct {
	session *gocql.Session
}

type UserDbInterface interface {
	InsertNewUser(email string, password string) (error)
	LoadUserByEmail(email string) (model.UserDb, error)
	LoadUserById(id string) (model.UserDb, error)
	UpdateUserPassword(userId string, hash string) (error)
	UpdateUserPushId(r *model.UpdatePushRequest) (error)
	LoadUserPushIds(userId string) ([]string, error)
	RemoveUserPush(userId string, push []string) (error)
	UpdateUserAvatar(userId string, blob []byte) (string, error)
}

type UserDataBase struct {
	session *gocql.Session
}

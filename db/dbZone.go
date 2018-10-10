package db

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"gobeacon/model"
	"time"
)

func LoadZonesByUserId(userId string) ([]model.GeoZoneDb, error) {
	stmt, names := qb.Select(tZones).Where(qb.Eq("user_id")).ToCql()

	var zones []model.GeoZoneDb
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId})
	err := q.SelectRelease(&zones)
	return zones, err
}

func CreateZoneForUser(r *model.ZoneCreateRequest) (model.GeoZoneDb, error) {
	id, _ := gocql.RandomUUID()
	usr, _ := gocql.ParseUUID(r.UserId)
	db := model.GeoZoneDb{Id: id, UserId: usr, Name: r.Name, CreatedAt: time.Now(), UpdatedAt: time.Now(), Points: r.Points, Trackers: make(map[gocql.UUID]bool)}

	stmt, names := qb.Insert(tZones).Columns("id", "user_id", "name", "created_at", "updated_at", "points", "trackers").ToCql()
	e := gocqlx.Query(session.Query(stmt), names).BindStruct(&db).ExecRelease()
	return db, e
}

func UpdateZone(r *model.ZoneCreateRequest) (error) {
	stmt, names := qb.Update(tZones).Set("name", "points").Where(qb.Eq("id")).ToCql()
	e := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": r.Id, "name": r.Name, "points": r.Points}).ExecRelease()
	return e
}

func LoadZoneById(id string) (model.GeoZoneDb, error) {
	stmt, names := qb.Select(tZones).Where(qb.Eq("id")).ToCql()
	var zone model.GeoZoneDb
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": id})
	err := q.GetRelease(&zone)
	return zone, err
}

func LoadZonesByTrackId(trackId string) ([]model.GeoZoneDb, error) {
	stmt, names := qb.Select(tZones).Where(qb.ContainsKey("trackers")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"trackers": trackId})

	var zones []model.GeoZoneDb
	err := q.SelectRelease(&zones)
	return zones, err
}

func DeleteZoneById(zoneId string) (error) {
	stmt, names := qb.Delete(tZones).Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": zoneId})
	return q.ExecRelease()
}

func UpdateZoneTrackers(zoneId string, track map[string]bool) (error) {
	stmt, names := qb.Update(tZones).Set("trackers").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"trackers": track, "id": zoneId})

	return q.ExecRelease()
}

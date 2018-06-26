package model

import "github.com/gocql/gocql"

type UserDb struct {

	Id gocql.UUID

	Email string

	Password string

	Trackers map[gocql.UUID]UserTrackers

	ZoneList []gocql.UUID

	Avatar   string

	Push_id  []string
}


type UserTrackers struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
package model

import "time"

type HeartbeatResponse struct {
	Code int `json:"code"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type ProfileResponse struct {
	Id       string        `json:"id"`
	Email    string        `json:"email"`
	Avatar   string        `json:"avatar"`
	Trackers []UserTracker `json:"trackers"`
}

type AvatarResponse struct {
	Data string `json:"data"`
}

type GeoZoneResponse struct {
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Points   []ZonePoint `json:"points"`
	Trackers []string    `json:"trackers"`
}

type TrackCreateResponse struct {
	Id string `json:"id"`
}

type TrackHistoryResponse struct {
	Date      time.Time `json:"datetime"`
	Latitude  float32   `json:"latitude"`
	Longitude float32   `json:"longitude"`
}

type TrackSnapZone struct {
	Id     string `json:"id,required" description:"Id трекера"`
	Inside bool   `json:"inside,required" description:"Отслеживать вход или выход из зоны (true - вход, false - выход)"`
}

type IBaseResponse interface {
	ToSerialize() string
}

type BaseResponse struct {
	Manufacter  string
	EquipmentId int64
	Type        MessageType
}

type LKResponse struct {
	BaseResponse
}

type ALResponse struct {
	BaseResponse
}

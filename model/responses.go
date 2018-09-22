package model

type HeartbeatResponse struct {
	Code int `json:"code"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type ProfileResponse struct {
	Id       string        `json:"id"`
	Email    string        `json:"name"`
	Avatar   string        `json:"avatar"`
	Trackers []UserTracker `json:"trackers"`
}

type GeoZoneResponse struct {
	Id       string          `json:"id"`
	Name     string          `json:"name"`
	Points   []ZonePoint     `json:"points"`
	Trackers []TrackSnapZone `json:"trackers"`
}

type TrackSnapZone struct {
	Id     string `json:"id,required" description:"Id трекера"`
	Inside bool   `json:"inside,required" description:"Отслеживать вход или выход из зоны (true - вход, false - выход)"`
}

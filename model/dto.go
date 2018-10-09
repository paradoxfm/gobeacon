package model

import "time"

type UserAuth struct {
	UserName  string `json: "id"`
	FirstName string
	LastName  string
}

type UserTracker struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Offset int    `json:"offset"`
}

type UserSession struct {
	Id    string
	Email string
}

type Heartbeat struct {
	IsGps           bool      `json:"is_gps_source"`
	IsGsm           bool      `json:"is_gsm_source"`
	IsWifi          bool      `json:"is_wifi_source"`
	Latitude        float32   `json:"latitude"`
	Longitude       float32   `json:"longitude"`
	Power           int       `json:"power"`
	DateTime        time.Time `json:"datetime" example:"RFC3339"` //RFC3339
	DeviceId        string    `json:"device_id"`
	IsLowPowerAlarm bool      `json:"is_low_power_alarm"`
	IsSOSAlarm      bool      `json:"is_sos_alarm"`
}

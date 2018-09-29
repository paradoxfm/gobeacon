package model

import (
	"mime/multipart"
	"time"
)

type RegistrationRequest struct {
	Email    string `json:"email,required" description:"Email пользователя"`
	Password string `json:"password,required," description:"Пароль пользователя"`
	Confirm  string `json:"сonfirm,required" description:"Подтверждение пароля"`
}

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	UserId      string
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

type GetProfileRequest struct {
	UserId string
}

type UpdateAvatarRequest struct {
	UserId string
	File   *multipart.FileHeader
}

type UpdateTrackAvatarRequest struct {
	UserId  string
	TrackId string
	File    *multipart.FileHeader
}

type UpdatePushRequest struct {
	UserId string
	PushId string `json:"push_id"`
}

type ZoneAllRequest struct {
	UserId string
}

type ZoneSnapRequest struct {
	UserId string
	Ids    []string `json:"trackers"`
}

type ZoneCreateRequest struct {
	UserId string
	Id     string      `json:"id"`
	Name   string      `json:"name"`
	Points []ZonePoint `json:"points"`
}

type TracksByIdsRequest struct {
	Ids []string `json:"ids"`
}

type TrackCreateRequest struct {
	UserId   string
	Name     string `json:"name"`
	DeviceId string `json:"equipment_id"`
}

type TracksNameRequest struct {
	TrackId string
	UserId  string
	Name    string `json:"name"`
}

type TracksHistRequest struct {
	TrackId  string    `json:"tracker_id"`
	DateFrom time.Time `json:"date_start"`
	DateTo   time.Time `json:"date_end"`
}

type HeartbeatRequest struct {
	Datetime     int64   `json:"datetime"`
	IsGPSSource  bool    `json:"is_gps_source"`
	IsGSMSource  bool    `json:"is_gsm_source"`
	IsWifiSource bool    `json:"is_wifi_source"`
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
	Power        float32 `json:"power"`
	DeviceId     string  `json:"device_id"`
}

package model

import "mime/multipart"

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

type UpdatePushRequest struct {
	UserId string
	PushId string `json:"push_id"`
}

type ZoneAllRequest struct {
	UserId string
}

type TracksByIdsRequest struct {
	Ids []string `json:"ids"`
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

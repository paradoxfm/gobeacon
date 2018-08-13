package model

import "mime/multipart"

type RegistrationRequest struct {
	Email    string `json:"email,required" description:"Email пользователя" valid:"email~Неправильный формат email,required~Поле email обязательно для заполнения"`
	Password string `json:"password,required," description:"Пароль пользователя" valid:"required~Не заполнен пароль,length(6|14)~Длина пароля от 6 до 14 символов"`
	Confirm  string `json:"сonfirm,required" description:"Подтверждение пароля" valid:"required~Не заполнено подтверждение пароля"`
}

type ResetPasswordRequest struct {
}

type ChangePasswordRequest struct {
	UserId string
}

type GetProfileRequest struct {
	UserId string
}

type UpdateAvatarRequest struct {
	UserId string
	Avatar interface{} `json:"avatar"`
	File   *multipart.FileHeader
}

type UpdatePushRequest struct {
	UserId string
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

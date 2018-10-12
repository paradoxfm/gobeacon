package model

import (
	"fmt"
	"mime/multipart"
	"strings"
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
	UserId string `json:"-"`
	PushId string `json:"push_id"`
}

type LoginRequest struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type ZoneAllRequest struct {
	UserId string
}

type ZoneSnapRequest struct {
	UserId string   `json:"-"`
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
	UserId   string `json:"-"`
	Imei     string `json:"-"`
	Name     string `json:"name"`
	DeviceId string `json:"equipment_id" example:"device id"`
}

type TrackPrefRequest struct {
	TrackId string `json:"-"`
	UserId  string `json:"-"`
	Name    string `json:"name"`
	Offset  int    `json:"offset"`
}

type TracksHistRequest struct {
	TrackId  string    `json:"tracker_id"`
	DateFrom time.Time `json:"date_start" example:"RFC3339"`
	DateTo   time.Time `json:"date_end" example:"RFC3339"`
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

type PositionData struct {
	Date               string
	Time               string
	WhetherTheLocation string
	Latitude           float32
	MarkOfLatitude     string
	Longitude          float32
	MarkOfLongitude    string
	Power              float32
	TerminalState      int16
}

type IBaseRequest interface {
	GetType() MessageType
	GetBase() BaseRequest
	GetPositionData() PositionData
}

type BaseRequest struct {
	Manufacter  string
	EquipmentId int64
	Type        MessageType
}

type LKRequest struct {
	BaseRequest
}

type UDRequest struct {
	BaseRequest
	PositionData
}

type UD2Request struct {
	BaseRequest
	PositionData
}

type ALRequest struct {
	BaseRequest
	PositionData
}

func (m BaseRequest) GetType() MessageType {
	return m.Type
}

func (m BaseRequest) GetBase() BaseRequest {
	return m
}

func (m LKRequest) GetPositionData() PositionData {
	return PositionData{}
}

func (m PositionData) GetPositionData() PositionData {
	return m
}

type MessageType int

const (
	LK MessageType = 1 + iota
	UD
	UD2
	AL
	TKQ
	TKQ2
	None
)

func ToMessageType(p string) MessageType {

	switch strings.ToUpper(strings.TrimSpace(p)) {
	case "LK":
		return LK
	case "UD":
		return UD
	case "UD2":
		return UD2
	case "AL":
		return AL
	case "TKQ":
		return TKQ
	case "TKQ2":
		return TKQ2
	}
	return None
}

func (p MessageType) ToString() string {
	switch p {
	case LK:
		return "LK"
	case UD:
		return "UD"
	case UD2:
		return "UD2"
	case AL:
		return "AL"
	case TKQ:
		return "TKQ"
	case TKQ2:
		return "TKQ2"
	}
	return fmt.Sprintf("MessageType(%d)", p)
}

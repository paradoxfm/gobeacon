package model

import (
	"fmt"
	"time"
)

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

type AvatarIdResponse struct {
	Id string `json:"url"`
}

type AvatarResponse struct {
	Data string `json:"data" example:"base64"`
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

type LoginResponse struct {
	Token string `json:"token"`
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

func (response BaseResponse) ToSerialize() string {
	return "[" + response.Manufacter + "*" + fmt.Sprintf("%v", response.EquipmentId) + "*" + string(lengthMessage(response.Type.ToString())) + "*" + response.Type.ToString() + "]"
}

// Считаем длину строки
// Examples:
//[3G*1208178692*0009*UPLOAD,30]
// input: UPLOAD
// output: 0009 (HEX)
// --------------
//[3G*1208178692*000E*MESSAGE,GOHOME]
//input: MESSAGE
// output: 000E (HEX)
func lengthMessage(mes string) []byte {

	size := len([]byte(mes))
	hex := fmt.Sprintf("%X", size)
	lenHex := len(string(size))

	if lenHex == 0 {
		return []byte("0000")
	} else if lenHex == 1 {
		return []byte("000" + hex)
	} else if lenHex == 2 {
		return []byte("00" + hex)
	} else if lenHex == 3 {
		return []byte("0" + hex)
	} else if lenHex == 4 {
		return []byte(hex)
	} else {
		return []byte("0000")
	}
}

type LKResponse struct {
	BaseResponse
}

type ALResponse struct {
	BaseResponse
}

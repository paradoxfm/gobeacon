package service

import (
	"bytes"
	"fmt"
	"gobeacon/model"
	"log"
	"strconv"
	"strings"
	"time"
)

func WatchHandleMessage2(msg string) (interface{}) {
	commands := strings.FieldsFunc(msg, splitByChars)
	var resp interface{}
	for _, str := range commands {
		str = strings.TrimSpace(str)
		rez := WatchHandleMessage([]byte(str))
		if resp == nil && rez != nil && len(rez.(string)) > 0 {
			resp = rez
		}
	}
	return resp
}

func splitByChars(r rune) bool {
	return r == '[' || r == ']'
}

func WatchHandleMessage(buff []byte) (interface{}) {
	message := parseMessage(buff)
	if message == nil {
		return nil
	}
	switch message.GetType() {
	case model.UD, model.UD2:
		hb := convertToHeartbeat(message)
		SaveHeartbeat(&hb)
	case model.LK:
		return model.LKResponse{BaseResponse: createBase(message.GetBase())}.ToSerialize()
	case model.AL:
		return model.ALResponse{BaseResponse: createBase(message.GetBase())}.ToSerialize()
	}
	return nil
}

func createBase(message model.BaseRequest) (model.BaseResponse) {
	return model.BaseResponse{Manufacter: message.Manufacter, EquipmentId: message.EquipmentId, Type: message.Type}
}

// гавно и лютая копипаста
func parseMessage(buff []byte) (model.IBaseRequest) {
	// обработчик сообщений
	// [3G*1208178692*0009*UPLOAD,30]
	splitBuf := bytes.Split(buff, []byte(",")) // разбиваем по запятым
	// [ [3G*1208178692*0009*UPLOAD], [30] ]
	header := bytes.Replace(splitBuf[0], []byte("["), []byte(""), -1) // убираем [
	// [ 3G*1208178692*0009*UPLOAD] ]
	splitHeader := bytes.Split(header, []byte("*")) // header разбиваем по *
	// [3G, 1208178692, 0009, UPLOAD]
	commandBytes := bytes.Replace(splitHeader[len(splitHeader)-1], []byte("]"), []byte(""), -1) // убираем ] получаем последний элемент

	command := string(bytes.Trim(commandBytes, "\x00"))
	command = strings.TrimSpace(command)

	// UPLOAD
	//Info(fmt.Sprintf("split_buf: ", string(split_buf[0])))
	//Info(fmt.Sprintf("header: ", string(header)))
	//Info(fmt.Sprintf("command: ", command))

	switch command {
	// Terminal to send commands
	case "LK": // 1. Link Keep
		return model.LKRequest{BaseRequest: createBaseRequest(splitHeader, command)}
		// TODO: проверить нужно ли, что-то запросить у часов
	case "UD": // 2. Position data report
		return model.UDRequest{BaseRequest: createBaseRequest(splitHeader, command), PositionData: createPositionData(splitBuf)}
	case "UD2": // 3. Blind spot Data Supplements
		return model.UD2Request{BaseRequest: createBaseRequest(splitHeader, command), PositionData: createPositionData(splitBuf)}
	case "AL": // 4. Alarm data report
		// Пришол сигнал тревоги
		return model.ALRequest{BaseRequest: createBaseRequest(splitHeader, command), PositionData: createPositionData(splitBuf)}
		// TODO: если нужно действие при получении сигнала тревоги
		//case bytes.Compare(command, []byte("WAD")) == 0: // 5. Requested address instruction
		//	Info("WAD")
		//case bytes.Compare(command, []byte("WG")) == 0: // 6. Requests the latitude and longitude of instruction
		//	Info("WG")
		//// Send command platform - Terminal response
		//case bytes.Compare(command, []byte("UPLOAD")) == 0: // 1. Data upload interval setting
		//	Info("UPLOAD")
		//case bytes.Compare(command, []byte("TS")) == 0: // 18. Query parameters
		//	Info("TS")
	default:
		log.Printf("unknown watch command %s", command)
		return nil
	}
}

func createBaseRequest(header [][]byte, command string) model.BaseRequest {
	manufacturer := header[0]
	equipmentId := header[1]
	equipmentId64, e := strconv.ParseInt(string(equipmentId), 10, 64)
	CheckError(e, "equipmentId64.Parse failed")
	b := model.BaseRequest{
		Manufacter:  string(manufacturer),
		EquipmentId: equipmentId64,
		Type:        model.ToMessageType(command),
	}
	return b
}

func createPositionData(posData [][]byte) model.PositionData {
	lat64, e := strconv.ParseFloat(string(posData[4]), 32)
	CheckError(e, "lat64.Parse failed")
	lat := float32(lat64)

	lon64, e := strconv.ParseFloat(string(posData[6]), 32)
	CheckError(e, "lon64.Parse failed")
	lon := float32(lon64)

	date := string(posData[1])
	tm := string(posData[2])

	power64, e := strconv.ParseFloat(string(posData[13]), 32)
	CheckError(e, "power64.Parse failed")
	power := float32(power64)

	state64, e := strconv.ParseFloat(string(posData[16]), 16)
	CheckError(e, "state64.Parse failed")
	state := int16(state64)

	return model.PositionData{Date: date, Time: tm, WhetherTheLocation: string(posData[3]), Latitude: lat,
		MarkOfLatitude: string(posData[5]), Longitude: lon, MarkOfLongitude: string(posData[7]), Power: power, TerminalState: state}
}

func CheckError(err error, msq string) {
	if err != nil {
		//Error(fmt.Sprintf("server '%s' error: %s", msq, err))
	}
}

// Пдоготовка сообщения для отправки в RabbitMQ
func convertToHeartbeat(req model.IBaseRequest) (model.Heartbeat) {
	p := req.GetPositionData()
	b := req.GetBase()
	dt := p.Date
	tm := p.Time
	//rfcStr := fmt.Sprintf("a %s", "string")
	tdStr := fmt.Sprintf("20%c%c-%c%c-%c%c", dt[4], dt[5], dt[2], dt[3], dt[0], dt[1])
	timeStr := fmt.Sprintf("%c%c:%c%c:%c%cZ", tm[0], tm[1], tm[2], tm[3], tm[4], tm[5])
	rfcStr := tdStr + "T" + timeStr
	datetime, e := time.Parse(time.RFC3339, rfcStr)

	CheckError(e, "time.Parse failed")
	lat := p.Latitude
	if p.MarkOfLatitude == "S" { // N expresses the north latitude, S expresses the south latitude.
		lat = -lat // Если юг, меняем знак
	}

	lon := p.Longitude
	if p.MarkOfLongitude == "W" { // E expresses the east longitude, W expresses the west longitude
		lon = -lon // Если запад, меняем знак
	}
	devId := strconv.FormatInt(b.EquipmentId, 10)

	return model.Heartbeat{
		DateTime:        datetime,
		IsGps:           p.WhetherTheLocation == "A",
		IsGsm:           p.WhetherTheLocation != "A",
		IsWifi:          false,
		Latitude:        lat,
		Longitude:       lon,
		Power:           int(p.Power),
		DeviceId:        devId,
		IsLowPowerAlarm: p.TerminalState == 17,
		IsSOSAlarm:      p.TerminalState == 16,
		//State:           p.TerminalState,
		//DeviceType:      "watch",
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kellydunn/golang-geo"
	"gobeacon/model"
	"gobeacon/service"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestBasicAuth(t *testing.T) {
	str := "1234"
	val := fmt.Sprintf("|%q|%q|\n", str[1], str[2])
	println(val)
}

func TestApple(t *testing.T) {
	codeB64 := "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"
	values := map[string]interface{}{"password": "cad4c4cd404644e78e4eef0cd35b907f", "receipt-data": codeB64, "exclude-old-transactions": true}
	jsonValue, _ := json.Marshal(values)
	resp, e := http.Post("https://sandbox.itunes.apple.com/verifyReceipt", "application/json", bytes.NewBuffer(jsonValue))
	if e != nil {
		return
	}
	if resp == nil {
		return
	}
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	res := model.AppleReceiptResponse{}
	if e = json.Unmarshal(body, &res); e != nil {
		return
	}
	if res.Status != 0 {
		return
	}
}

func TestZoneAlarm(t *testing.T) {
	// работа la: 56.813248 lo: 60.59087
	req := model.Heartbeat{IsGps: true, Latitude: 56.813248, Longitude: 60.59087, Power: 99, DateTime: time.Now(), DeviceId: "c60050f8255acc10"}
	service.SaveHeartbeat(&req)
	// вне зоны 56.814017, 60.592747
	req2 := model.Heartbeat{IsGps: true, Latitude: 56.814017, Longitude: 60.592747, Power: 99, DateTime: time.Now(), DeviceId: "c60050f8255acc10"}
	service.SaveHeartbeat(&req2)
}

func TestCopyTrack(t *testing.T) {
	t1 := model.Tracker{DeviceId: "asdfasdf", LatitudeLast: 56.814017, LongitudeLast: 60.592747}
	t2 := new(model.Tracker)
	*t2 = *&t1
	t1.DeviceId = "dfgh"
	println(t1.DeviceId)
	println(t2.DeviceId)
	println(t1.LatitudeLast)
	println(t2.LatitudeLast)
}

func TestPointDist(t *testing.T) {
	// Make a few points
	p := geo.NewPoint(56.813554, 60.590319)
	p2 := geo.NewPoint(56.812955, 60.590383)

	// find the great circle distance between them
	dist := p.GreatCircleDistance(p2) * 1000
	fmt.Printf("great circle distance: %f\n", dist)
}

func TestDateFormat(t *testing.T) {
	dt := "071018"
	println(fmt.Sprintf("20%c%c-%c%c-%c%c", dt[4], dt[5], dt[2], dt[3], dt[0], dt[1]))
}

/*func TestParseWatchData(t *testing.T) {
	dat := "[3G*1208178692*000A*LK,0,0,100][3G*1208178692*00C5*UD,071018,145708,A,56.822265,N,60.6324567,E,3.55,231.9,0.0,4,100,100,0,0,00000008,7,255,250,1,6624,501,158,6624,15231,149,6624,1301,146,6624,15232,143,6624,3003,141,6624,182,141,6624,185,136,0,46.2]"
	service.WatchHandleMessage2(dat)
}*/

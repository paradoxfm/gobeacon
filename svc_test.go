package main

import (
	"fmt"
	"github.com/kellydunn/golang-geo"
	"gobeacon/model"
	"gobeacon/service"
	"testing"
	"time"
)

func TestBasicAuth(t *testing.T) {
	str := "1234"
	val := fmt.Sprintf("|%q|%q|\n", str[1], str[2])
	println(val)
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

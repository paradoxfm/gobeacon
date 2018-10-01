package main

import (
	"fmt"
	"github.com/kellydunn/golang-geo"
	"gobeacon/model"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	str := "1234"
	val := fmt.Sprintf("|%q|%q|\n", str[1], str[2])
	println(val)
}


func TestCopyTrack(t *testing.T) {
	t1 := model.Tracker{DeviceId:"asdfasdf"}
	t2 := new(model.Tracker)
	*t2 = *&t1
	t1.DeviceId = "dfgh"
	println(t1.DeviceId)
	println(t2.DeviceId)
}

func TestPointDist(t *testing.T) {
	// Make a few points
	p := geo.NewPoint(56.813554, 60.590319)
	p2 := geo.NewPoint(56.812955, 60.590383)

	// find the great circle distance between them
	dist := p.GreatCircleDistance(p2) * 1000
	fmt.Printf("great circle distance: %f\n", dist)
}

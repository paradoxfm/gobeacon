package main

import (
	"encoding/json"
	"github.com/douglasmakey/go-fcm"
	"gobeacon/model"
	"log"
	"testing"
)

func TestJson(t *testing.T) {
	res := model.AppleReceiptResponse{}
	body := []byte("{}")
	if e := json.Unmarshal(body, &res); e != nil {
		return
	}
	log.Print(res)
}

func TestPush(t *testing.T) {
	client := fcm.NewClient("")

	//
	pushId := []string{
		"f9dUgUE7YBI:APA91bGg5Eb9biUvVDQx8nOMNaQhAekgPa4w49cHeTwWGgTeAYREjSA8BYecDZyEmWgscPWFVZDmeglC7LWrn6z68wMbQae1S94RmNogeX4y51QsGSjm-kFHfNN4VIaHFcOefemt7fOW",
	}

	data := map[string]interface{}{
		"notification": map[string]interface{}{
			"title": "Тест пуш",
			"body": "Тест пуш",
		},
	}

	// You can use PushMultiple or PushSingle
	client.PushMultiple(pushId, data)
	badRegistrations := client.CleanRegistrationIds()

	log.Print("bad reg ")
	log.Println(badRegistrations)

	status, err := client.Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println(status.Results)
}

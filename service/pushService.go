package service

import (
	"github.com/douglasmakey/go-fcm"
	"gobeacon/db"
	"log"
)

func SendPushForUser(userId string, pushIds []string, data interface{}) {
	client := fcm.NewClient(Config().ServerKey)
	if len(pushIds) != 0 {
		sendPush(userId, pushIds, data, client)
	}
}

func sendPush(userId string, pushIds []string, data interface{}, client *fcm.Client) {
	client.PushMultiple(pushIds, data)
	badRegistrations := client.CleanRegistrationIds()
	if len(badRegistrations) > 0 {
		db.RemoveUserPush(userId, badRegistrations)
		//log.Println(badRegistrations)
	}
	status, err := client.Send()
	if err != nil {
		log.Printf("error: %v", err)
	}

	log.Println(status.Results)
}

func SendPushNotification(userId string) {
	ids, e := db.LoadUserPushIds(userId)
	if e != nil {
		return
	}
	client := fcm.NewClient(Config().ServerKey)

	data := map[string]interface{}{
		"message":      lowPowerMsgId,
		"tracker_id":   "Трекер тест",
		"tracker_name": "Трекер тест",
	}

	client.PushMultiple(ids, data)

	badRegistrations := client.CleanRegistrationIds()
	if len(badRegistrations) > 0 {
		db.RemoveUserPush(userId, badRegistrations)
		//log.Println(badRegistrations)
	}

	status, err := client.Send()
	if err != nil {
		log.Printf("error: %v", err)
	}

	log.Println(status.Results)
}

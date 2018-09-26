package service

import (
	"github.com/douglasmakey/go-fcm"
	"log"
)

func SendPushForUsers(userPush map[string][]string, data interface{}) {
	client := fcm.NewClient(Config().ServerKey)
	for userId, pushIds := range userPush {
		if len(pushIds) != 0 {
			sendPush(userId, pushIds, data, client)
		}
	}
}

func sendPush(userId string, pushIds []string, data interface{}, client *fcm.Client) {
	client.PushMultiple(pushIds, data)
	badRegistrations := client.CleanRegistrationIds()
	if len(badRegistrations) > 0 {
		removeInvalidPush(userId, badRegistrations)
		//log.Println(badRegistrations)
	}
	status, err := client.Send()
	if err != nil {
		log.Printf("error: %v", err)
	}

	log.Println(status.Results)
}

func MoveTrackerSettings() {
	users := getAllUsers()
	for _, usr := range users {
		for trkId, trk := range usr.Trackers {
			insertTrackSettings(usr.Id, trkId, trk.Name)
		}
	}
}

func SendPushNotification(userId string) {
	ids, e := getUserPushIds(userId)
	if e != nil {
		return
	}
	client := fcm.NewClient(Config().ServerKey)

	data := map[string]interface{}{
		"message": "Тестовое оповещение при выходе трекера из зоны, возможны ураганы и шквалистый ветер",
		"tracker_id": "",
	}

	client.PushMultiple(ids, data)

	badRegistrations := client.CleanRegistrationIds()
	if len(badRegistrations) > 0 {
		removeInvalidPush(userId, badRegistrations)
		//log.Println(badRegistrations)
	}

	status, err := client.Send()
	if err != nil {
		log.Printf("error: %v", err)
	}

	log.Println(status.Results)
}

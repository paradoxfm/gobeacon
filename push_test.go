package main

import (
	"github.com/douglasmakey/go-fcm"
	"gobeacon/service"
	"log"
	"testing"
)

func TestPush(t *testing.T) {
	client := fcm.NewClient(service.Config().ServerKey)

	pushId := []string{
		"efucSs6EJuE:APA91bFgtVcsjDxzphJROz8czIssVUy0r53addXJAuuaienMMttNpdDKO4ofKczB0e0BKuLEbCaEZGBSA7ynojmwYExcoVir_BwIX1r30GgYzGS3BnQoxIAyI1D3pDttOZ8Rm_UoDiOl",
	}

	data := map[string]interface{}{
		"message":      2001,
		"content_available": true,
		"tracker_id":   "Трекер тест",
		"tracker_name": "Трекер тест",
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

package service

import (
	"gobeacon/code"
	"gobeacon/db"
	"gobeacon/model"
)

func CurrentSubscription(userId string) ([]model.UserSubscription, []int) {
	var err []int
	activeSubscriptions, e := db.LoadUserCurrentSubscriptions(userId)
	if e != nil {
		return nil, append(err, code.DbError)
	}
	subMap, err := getSubscriptionsMap(err)
	if err != nil {
		return nil, err
	}
	rez := make([]model.UserSubscription, len(activeSubscriptions))
	for i := 0; i < len(activeSubscriptions); i++ {
		bs := activeSubscriptions[i]
		rez[i] = model.UserSubscription{Title: subMap[bs.Item.String()], DateFrom: bs.EnableFrom, DateTo: bs.EnableTo}
	}
	return rez, err
}

func AllActiveSubscription(userId string) ([]model.UserSubscription, []int) {
	var err []int
	buySubscriptions, e := db.LoadUserSubscriptions(userId)
	if e != nil {
		return nil, append(err, code.DbError)
	}
	subMap, err := getSubscriptionsMap(err)
	if err != nil {
		return nil, err
	}

	rez := make([]model.UserSubscription, len(buySubscriptions))
	for i := 0; i < len(buySubscriptions); i++ {
		bs := buySubscriptions[i]
		rez[i] = model.UserSubscription{Title: subMap[bs.Item.String()], DateFrom: bs.EnableFrom, DateTo: bs.EnableTo}
	}
	return rez, err
}

func getSubscriptionsMap(err []int) (map[string]string, []int) {
	subscriptions, err := Subscriptions()
	if err != nil {
		return nil, err
	}
	subMap := make(map[string]string)
	for i := 0; i < len(subscriptions); i++ {
		sub := subscriptions[i]
		subMap[sub.Id.String()] = sub.Title
	}
	return subMap, nil
}

func Subscriptions() ([]model.Subscription, []int) {
	var err []int
	subscriptions, e := db.LoadSubscriptions()
	if e != nil {
		return nil, append(err, code.DbError)
	}
	return subscriptions, err
}

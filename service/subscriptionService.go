package service

import (
	"github.com/gocql/gocql"
	"gobeacon/code"
	"gobeacon/db"
	"gobeacon/model"
	"time"
)

const maxFamilyCount = 5

func BuySubscription(req *model.BuySubscriptionRequest) []int {
	var err []int
	if len(req.Accounts) > maxFamilyCount - 1 {
		return append(err, code.MaxSubscriptionCount)
	}
	users := make([]model.UserDb, len(req.Accounts) + 1)
	var e error
	users[len(users)], e = db.LoadUserById(req.UserId)
	if e != nil {
		return append(err, code.DbError)
	}
	for idx, acc := range req.Accounts {
		users[idx], e = db.LoadUserByEmail(acc)
		if e != nil {
			return append(err, code.InvalidUserAccount)
		}
	}

	sub, e := db.LoadSubscriptionById(req.SubId)
	if e != nil {
		return append(err, code.DbError)
	}
	if !sub.Enabled {
		return append(err, code.DisabledSubscription)
	}
	dateTo := req.DateFrom.Add(time.Hour * time.Duration(sub.Length*24))
	var subs = make([]model.BuySubscription, len(users))
	for idx, usr := range users {
		uuid, _ := gocql.RandomUUID()
		subs[idx] = model.BuySubscription{Id: uuid, User: usr.Id, Item: sub.Id, BuyDate: time.Now(), EnableFrom: req.DateFrom, EnableTo: dateTo}
	}
	e = db.SaveSubscriptions(subs)
	if e != nil {
		return append(err, code.DbError)
	}
	return nil
}

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

package service

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"gobeacon/code"
	"gobeacon/db"
	"gobeacon/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const maxFamilyCount = 5

func ExtendSubscription(userId string) []int {
	var err []int
	active, e := db.LoadUserCurrentSubscriptions(userId)
	if e != nil {
		return append(err, code.DbError)
	}
	if active == nil || len(active) == 0 {
		if active, e = db.LoadUserLastSubscriptions(userId); e != nil {
			return append(err, code.DbError)
		}
		if active == nil || len(active) == 0 {
			return append(err, code.NoActiveSubscription)
		}
	}
	buySub, e := db.LoadBuySubscriptionByGroup(active[0].GroupId.String())
	newBuySubs := make([]model.BuySubscription, len(buySub))
	groupId, _ := gocql.RandomUUID()
	for ind, bs := range buySub {
		uuid, _ := gocql.RandomUUID()
		now := time.Now()
		dateTo := now.Add(time.Hour * time.Duration(30*24))
		newBuySubs[ind] = model.BuySubscription{GroupId: groupId, Item: bs.Item, User: bs.User, Id: uuid, Trial: false, BuyDate: now, EnableFrom: now, EnableTo: dateTo}
	}
	if e = db.SaveSubscriptions(newBuySubs); e != nil {
		return append(err, code.DbError)
	}
	return nil
}

func SendQueryApple(codeB64 string) (*model.AppleReceiptResponse, []int) {
	cfg := Config()
	resp, e := http.PostForm(cfg.AppleValidationUrl, url.Values{
		"password":     {cfg.AppleValidationKey},
		"receipt-data": {codeB64},
	})
	if e != nil {
		return nil, []int{code.ErrorHttpSubscription}
	}
	if resp == nil {
		return nil, []int{code.ErrorHttpResponseSubscription}
	}
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	res := model.AppleReceiptResponse{}
	if e = json.Unmarshal(body, &res); e != nil {
		return nil, []int{code.ErrorHttpResponseSubscription}
	}
	if res.Status != 0 {
		return nil, []int{code.ErrorHttpResponseSubscription}
	}
	return &res, nil
}

func AddUserToMySubscription(req *model.AddSubscriptionRequest) []int {
	var err []int
	active, e := db.LoadUserCurrentSubscriptions(req.UserId)
	if e != nil {
		return append(err, code.DbError)
	}
	if active == nil || len(active) == 0 {
		return append(err, code.NoActiveSubscription)
	}
	buySubs, e := db.LoadBuySubscriptionByGroup(active[0].GroupId.String())
	if len(req.Accounts)+len(buySubs) > maxFamilyCount {
		return append(err, code.MaxSubscriptionCount)
	}
	users := make([]model.UserDb, len(req.Accounts))
	for idx, acc := range req.Accounts {
		users[idx], e = db.LoadUserByEmail(acc)
		if e != nil {
			return append(err, code.InvalidUserAccount)
		}
	}
	var subs = make([]model.BuySubscription, len(users))
	for idx, usr := range users {
		uuid, _ := gocql.RandomUUID()
		subs[idx] = model.BuySubscription{Id: uuid, User: usr.Id, Item: active[0].Item, BuyDate: active[0].BuyDate,
			EnableFrom: active[0].EnableFrom, EnableTo: active[0].EnableTo, GroupId: active[0].GroupId}
	}
	e = db.SaveSubscriptions(subs)
	if e != nil {
		return append(err, code.DbError)
	}
	return nil
}

func GetAllAccountWithMySubscription(userId string) ([]model.UserBuySubResponse, []int) {
	var err []int
	active, e := db.LoadUserCurrentSubscriptions(userId)
	if e != nil {
		return nil, append(err, code.DbError)
	}
	if active == nil || len(active) == 0 {
		return nil, append(err, code.NoActiveSubscription)
	}
	buySub, e := db.LoadBuySubscriptionByGroup(active[0].GroupId.String())
	resp := make([]model.UserBuySubResponse, len(buySub))
	for i, bs := range buySub {
		usr, e := db.LoadUserById(bs.User.String())
		if e != nil {
			return nil, append(err, code.DbError)
		}
		resp[i] = model.UserBuySubResponse{Email: usr.Email}
	}
	return resp, nil
}

func BuySubscription(req *model.BuySubscriptionRequest) []int {
	var err []int
	var e error

	sub, e := db.LoadSubscriptionById(req.SubId)
	if e != nil {
		return append(err, code.DbError)
	}
	if len(req.Accounts) > maxFamilyCount-1 {
		return append(err, code.MaxSubscriptionCount)
	}
	var users []model.UserDb
	if owner, e := db.LoadUserById(req.UserId); e == nil {
		users = append(users, owner)
	} else {
		return append(err, code.DbError)
	}
	if sub.Payable {
		for idx, acc := range req.Accounts {
			users[idx], e = db.LoadUserByEmail(acc)
			if e != nil {
				return append(err, code.InvalidUserAccount)
			}
		}
	}

	if !sub.Enabled {
		return append(err, code.DisabledSubscription)
	}
	groupId, _ := gocql.RandomUUID()
	dateTo := req.DateFrom.Add(time.Hour * time.Duration(sub.Length*24))
	var subs = make([]model.BuySubscription, len(users))
	for idx, usr := range users {
		uuid, _ := gocql.RandomUUID()
		subs[idx] = model.BuySubscription{Id: uuid, User: usr.Id, Item: sub.Id, BuyDate: time.Now(), EnableFrom: req.DateFrom, EnableTo: dateTo, GroupId: groupId}
	}
	e = db.SaveSubscriptions(subs)
	if e != nil {
		return append(err, code.DbError)
	}
	if !sub.Payable {
		if e := db.UpdateUserUsedTrial(req.UserId); e != nil {
			return append(err, code.DbError)
		}
	}
	return nil
}

func CurrentSubscription(userId string) ([]model.UserSubscription, []int) {
	var err []int
	activeSubscriptions, e := db.LoadUserCurrentSubscriptions(userId)
	if e != nil {
		return nil, append(err, code.DbError)
	}
	subMap, err := getSubscriptionsMap(userId, err)
	if err != nil {
		return nil, err
	}
	rez := make([]model.UserSubscription, len(activeSubscriptions))
	for i := 0; i < len(activeSubscriptions); i++ {
		bs := activeSubscriptions[i]
		rez[i] = model.UserSubscription{Title: subMap[bs.Item.String()], DateFrom: bs.EnableFrom, DateTo: bs.EnableTo, Trial: bs.Trial}
	}
	return rez, err
}

func AllActiveSubscription(userId string) ([]model.UserSubscription, []int) {
	var err []int
	buySubscriptions, e := db.LoadUserSubscriptions(userId)
	if e != nil {
		return nil, append(err, code.DbError)
	}
	subMap, err := getSubscriptionsMap(userId, err)
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

func getSubscriptionsMap(userId string, err []int) (map[string]string, []int) {
	subscriptions, err := Subscriptions(userId)
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

func Subscriptions(userId string) ([]model.Subscription, []int) {
	var err []int
	userDb, eu := db.LoadUserById(userId)
	if eu != nil {
		return nil, append(err, code.DbError)
	}
	subscriptions, e := db.LoadSubscriptions()
	if e != nil {
		return nil, append(err, code.DbError)
	}
	if userDb.UsedTrial {
		var rez []model.Subscription
		for i := 0; i < len(subscriptions); i++ {
			s := subscriptions[i]
			if !s.Payable {
				continue
			}
			rez = append(rez, s)
		}
		return rez, err
	} else {
		return subscriptions, err
	}
}

package db

import (
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"gobeacon/model"
	"time"
)

func SaveSubscriptions(data []model.BuySubscription) error {
	stmt, names := qb.Insert(tBuySubscription).Columns("id", "user_id", "subscription_id", "buy_date", "enable_from", "enable_to").ToCql()
	for _, bsub := range data {
		e := gocqlx.Query(session.Query(stmt), names).BindStruct(&bsub).ExecRelease()
		if e != nil {
			return e
		}
	}
	return nil
}

func LoadSubscriptionById(id string) (model.Subscription, error) {
	stmt, names := qb.Select(tSubscription).Where(qb.Eq("id")).Limit(1).ToCql()
	var s model.Subscription

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": id,})
	err := q.GetRelease(&s)
	return s, err
}

func LoadUserCurrentSubscriptions(userId string) ([]model.BuySubscription, error) {
	stmt, names := qb.Select(tBuySubscription).Where(qb.Eq("user_id")).Where(qb.Lt("enable_to")).Where(qb.Gt("enable_from")).ToCql()
	var sub []model.BuySubscription

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId, "enable_to": time.Now(), "enable_from": time.Now()})
	err := q.SelectRelease(&sub)
	return sub, err
}

func LoadUserSubscriptions(userId string) ([]model.BuySubscription, error) {
	stmt, names := qb.Select(tBuySubscription).Where(qb.Eq("user_id")).Where(qb.Gt("enable_to")).ToCql()
	var sub []model.BuySubscription

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"user_id": userId, "enable_to": time.Now()})
	err := q.SelectRelease(&sub)
	return sub, err
}

func LoadSubscriptions() ([]model.Subscription, error) {
	stmt, names := qb.Select(tSubscription).Where(qb.Eq("enabled")).ToCql()
	var sub []model.Subscription

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"enabled": true})
	err := q.SelectRelease(&sub)
	return sub, err
}

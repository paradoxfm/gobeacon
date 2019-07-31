package db

import (
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"gobeacon/model"
	"time"
)

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

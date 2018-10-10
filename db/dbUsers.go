package db

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"gobeacon/model"
	"time"
)

func InsertNewUser(email string, password string) (error) {
	stmt, _ := qb.Insert(tUsers).Columns("id", "email", "password", "created_at").ToCql()
	q := session.Query(stmt)
	err := q.Bind(gocql.TimeUUID(), email, password, time.Now()).Exec()

	return err
}

func LoadUserByEmail(email string) (model.UserDb, error) {
	stmt, names := qb.Select(tUsers).Columns("id", "email", "password").Where(qb.Eq("email")).Limit(1).ToCql()
	var u model.UserDb

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"email": email,})
	err := q.GetRelease(&u)
	return u, err
}

func LoadUserById(id string) (model.UserDb, error) {
	stmt, names := qb.Select(tUsers).Columns("id", "email", "password", "avatar").Where(qb.Eq("id")).Limit(1).ToCql()
	var u model.UserDb

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": id,})
	err := q.GetRelease(&u)
	return u, err
}

func UpdateUserPassword(userId string, hash string) (error) {
	stmt, names := qb.Update(tUsers).Set("password").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": userId, "password": hash})

	err := q.ExecRelease()
	return err
}

func UpdateUserPushId(r *model.UpdatePushRequest) (error) {
	stmt, names := qb.Update(tUsers).Add("push_id").Set("updated_at").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": r.UserId, "push_id": []string{r.PushId}, "updated_at": time.Now()})
	err := q.ExecRelease()
	return err
}

func LoadUserPushIds(userId string) ([]string, error) {
	stmt, names := qb.Select(tUsers).Columns("id", "push_id").Where(qb.Eq("id")).Limit(1).ToCql()
	var u model.UserDb

	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": userId})
	err := q.GetRelease(&u)
	return u.PushId, err
}

func RemoveUserPush(userId string, push []string) error {
	stmt, names := qb.Update(tUsers).Remove("push_id").Set("updated_at").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": userId, "push_id": push, "updated_at": time.Now()})
	err := q.ExecRelease()
	return err
}

func UpdateUserAvatar(userId string, blob []byte) (string, error) {
	userDb, _ := LoadUserById(userId)

	avatarLink := userDb.Avatar
	avaId, e := insertNewAvatar(blob)
	if e != nil {
		return "", e
	}

	// обновляем ссылку на аватар
	stmt, names := qb.Update(tUsers).Set("avatar").Where(qb.Eq("id")).ToCql()
	q := gocqlx.Query(session.Query(stmt), names).BindMap(qb.M{"id": userId, "avatar": avaId})
	if err := q.ExecRelease(); err != nil {
		return "", err
	}

	if e = deleteAvatar(avatarLink); e != nil {
		return "", e
	}

	return avaId, nil
}

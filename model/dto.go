package model

type UserAuth struct {
	UserName  string `json: "id"`
	FirstName string
	LastName  string
}

type UserTracker struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

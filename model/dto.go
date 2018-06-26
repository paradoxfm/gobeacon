package model

type UserAuth struct {
	UserName  string `json: "id"`
	FirstName string
	LastName  string
}

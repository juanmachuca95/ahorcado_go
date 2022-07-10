package models

type Login struct {
	Username string `bson:"username" json:"username,omitempty"`
	Password string `bson:"password" json:"password,omitempty"`
}

func NewLogin(username, password string) *Login {
	return &Login{Username: username, Password: password}
}

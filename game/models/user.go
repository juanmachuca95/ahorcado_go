package models

type User struct {
	Username string `bson:"username" json:"username,omitempty"`
	Password string `bson:"password" json:"password,omitempty"`
}

func NewUser(username string) *User {
	return &User{
		Username: username,
	}
}

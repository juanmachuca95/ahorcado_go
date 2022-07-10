package jsmodels

import "github.com/gopherjs/gopherjs/js"

type RequestLogin struct {
	*js.Object
	Username string `js:"username"`
	Password string `js:"password"`
}

type ResponseLogin struct {
	*js.Object
	Token string `js:"token"`
}

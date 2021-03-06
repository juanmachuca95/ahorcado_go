package jsmodels

import "github.com/gopherjs/gopherjs/js"

type Ranking struct {
	*js.Object
	Id       string `js:"id"`
	Username string `js:"username"`
	Won      int32  `js:"won"`
}

type Rankings struct {
	*js.Object
	Rankings []Ranking `js:"rankings"`
}

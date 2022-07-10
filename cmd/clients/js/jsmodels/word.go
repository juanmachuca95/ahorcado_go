package jsmodels

import "github.com/gopherjs/gopherjs/js"

type Word struct {
	*js.Object
	Game_id string `js:"game_id"`
	Word    string `js:"word"`
	User    string `js:"user"`
}

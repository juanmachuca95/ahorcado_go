package jsmodels

import "github.com/gopherjs/gopherjs/js"

type Game struct {
	*js.Object
	Id          string   `js:"id"`
	Word        string   `js:"word"`
	Winner      string   `js:"winner"`
	Encontrados []string `js:"encontrados"`
	Finalizada  bool     `js:"finalizada"`
	Error       string   `js:"error"`
	UserSend    string   `js:"user_send"`
	WordSend    string   `js:"word_send"`
	Status      string   `js:"status"`
}

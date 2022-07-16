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
	Usersend    string   `js:"usersend"`
	Wordsend    string   `js:"wordsend"`
	Status      int      `js:"status"`
}

package main

//go:generate gopherjs build -m jsclient/client.go -o jsclient/html/index.js
//go:generate go-bindata -pkg compiled -nometadata -o jsclient/compiled/client.go -prefix html ./jsclient/html
//go:generate bash -c "rm clientjs/html/*.js*"

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/juanmachuca95/ahorcado_go/cmd/clients/js/jsmodels"
	vue "github.com/oskca/gopherjs-vue"
)

// Model is the state keeper of the app.

func main() {
	m := &jsmodels.Model{
		Object: js.Global.Get("Object").New(),
	}

	// These must be set after the struct has been initialised
	// so that the values can be mirrored into the internal JS Object.
	m.Game = &jsmodels.Game{}
	m.Word = &jsmodels.Word{}

	m.Username = ""
	m.Status = ""
	m.Error = ""
	m.FoundLetters = ""
	m.Winner = ""
	m.GameData = []*jsmodels.Game{}
	m.ConnOpen = false
	m.Tries = 6

	m.InputWord = ""
	m.InputUser = ""

	// GetGame retorna el juego
	m.GetGame()
	m.Connect()
	m.Received()

	// Create the VueJS viewModel using a struct pointer
	vue.New("#app", m)
}

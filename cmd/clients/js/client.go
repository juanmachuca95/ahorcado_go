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

	// Game
	m.Game = &jsmodels.Game{}
	m.Word = &jsmodels.Word{}

	// Login
	m.LoginUser = true
	m.Username = ""
	m.Password = ""
	m.Token = ""
	m.ErrorLogin = ""

	// Register
	m.RegisterUser = false
	m.ErrorRegister = ""
	m.Status = ""
	m.Error = ""
	m.FoundLetters = ""
	m.Winner = ""
	m.GameData = []*jsmodels.Game{}
	m.RankingData = &jsmodels.Rankings{}
	m.ConnOpen = false
	m.Tries = 6

	m.InputWord = ""
	m.InputUser = ""

	// Game
	m.Jugando = false
	m.Ranking = false
	m.Lost = ""
	m.Win = ""

	// Create the VueJS viewModel using a struct pointer
	vue.New("#app", m)
}

package main

//go:generate gopherjs build -m jsclient/client.go -o jsclient/html/index.js
//go:generate go-bindata -pkg compiled -nometadata -o jsclient/compiled/client.go -prefix html ./jsclient/html
//go:generate bash -c "rm clientjs/html/*.js*"

import (
	"github.com/gopherjs/gopherjs/js"
	json "github.com/johanbrandhorst/gopherjs-json"
	vue "github.com/oskca/gopherjs-vue"
	"honnef.co/go/js/xhr"
)

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
}

// Model is the state keeper of the app.
type Model struct {
	*js.Object
	DataGame *Game `js:"game"`
}

func (m *Model) GetGame() {
	req := xhr.NewRequest("GET", "http://localhost:8090/api/v1/game")
	req.SetRequestHeader("Content-Type", "application/json")

	// Wrap call in goroutine to use blocking code
	go func() {
		// Blocks until reply received
		err := req.Send(nil)
		if err != nil {
			panic(err)
		}

		if req.Status != 200 {
			panic(req.ResponseText)
		}

		rObj, err := json.Unmarshal(req.ResponseText)
		if err != nil {
			panic(err)
		}

		msg := &Game{
			Object: rObj,
		}

		m.DataGame = msg
	}()
}

func main() {
	m := &Model{
		Object: js.Global.Get("Object").New(),
	}

	// These must be set after the struct has been initialised
	// so that the values can be mirrored into the internal JS Object.
	m.DataGame = nil

	// Create the VueJS viewModel using a struct pointer
	vue.New("#app", m)
}

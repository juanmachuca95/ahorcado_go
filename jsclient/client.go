package main

//go:generate gopherjs build -m jsclient/client.go -o jsclient/html/index.js
//go:generate go-bindata -pkg compiled -nometadata -o jsclient/compiled/client.go -prefix html ./jsclient/html
//go:generate bash -c "rm clientjs/html/*.js*"

import (
	"log"
	"net"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
	json "github.com/johanbrandhorst/gopherjs-json"
	help "github.com/juanmachuca95/ahorcado_go/helpers"
	vue "github.com/oskca/gopherjs-vue"
	"honnef.co/go/js/xhr"
)

var WSConn net.Conn

type Word struct {
	*js.Object
	Game_id string `js:"game_id"`
	Word    string `js:"word"`
	User    string `js:"user"`
}

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
	Game         *Game   `js:"game"`
	Word         *Word   `js:"word"`
	FoundLetters string  `js:"found_letters"`
	Input_user   string  `js:"input_user"`
	Input_word   string  `js:"input_word"`
	GameData     []*Game `js:"game_data"`
	Status       string  `js:"status"`
	ConnOpen     bool    `js:"ws_conn"`
	BidiMessages []*Game `js:"bidi_messages"`
}

func main() {
	m := &Model{
		Object: js.Global.Get("Object").New(),
	}

	// These must be set after the struct has been initialised
	// so that the values can be mirrored into the internal JS Object.
	m.Game = &Game{}
	m.Word = &Word{}
	m.Input_word = ""
	m.Input_user = "Juancete"
	m.FoundLetters = ""
	m.Status = ""
	m.GameData = []*Game{}
	m.ConnOpen = false
	m.BidiMessages = []*Game{}

	// GetGame retorna el juego
	m.GetGame()
	m.Connect()
	m.Received()

	// Create the VueJS viewModel using a struct pointer
	vue.New("#app", m)
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

		m.Game = msg
		m.FoundLetters = help.ShowWord(m.Game.Word, m.Game.Encontrados)
	}()
}

func (m *Model) Connect() {
	// Wrap call in goroutine to use blocking code
	go func() {
		// Blocks until connection is established
		var err error
		WSConn, err = websocket.Dial("ws://localhost:8090/api/v1/playing")
		if err != nil {
			panic(err)
		}

		m.ConnOpen = true
	}()
}

func (m *Model) Close() {
	err := WSConn.Close()
	if err != nil {
		panic(err)
	}

	m.ConnOpen = false
	m.Game = &Game{}
	m.Word = &Word{}
}

func getStreamMessage(msg string) *Game {
	rObj, err := json.Unmarshal(msg)
	if err != nil {
		panic(err.Error())
	}

	// The actual message is wrapped in a "result" key,
	// and there might be an error returned as well.
	// See https://github.com/grpc-ecosystem/grpc-gateway/blob/b75dbe36289963caa453a924bd92ddf68c3f2a62/runtime/handler.go#L163
	aux := &struct {
		*js.Object
		msg *Game `js:"result"`
	}{
		Object: rObj,
	}

	// The most reliable way I've found to check whether
	// an error was returned.
	if rObj.Get("error").Bool() {
		panic(msg)
	}

	return aux.msg
}

func (m *Model) Send() {
	msg := &Word{
		Object: js.Global.Get("Object").New(),
	}

	msg.Game_id = m.Game.Id
	msg.Word = m.Input_word
	msg.User = m.Input_user

	s, err := json.Marshal(msg.Object)
	if err != nil {
		panic(err)
	}

	_, err = WSConn.Write([]byte(s))
	if err != nil {
		panic(err)
	}

	m.Received()
}

func (m *Model) Received() {
	buf := make([]byte, 1024)
	// Wrap call in goroutine to use blocking code
	go func() {
		// Blocks until a WebSocket frame is received
		for m.ConnOpen {
			n, err := WSConn.Read(buf)
			if err != nil {
				panic(err.Error())
			}

			game := getStreamMessage(string(buf[:n]))
			log.Println("game .", game.UserSend)
			m.Status = game.Error
			m.FoundLetters = help.ShowWord(game.Word, game.Encontrados)
			m.Input_word = ""
			m.GameData = append(m.GameData, game)
		}

	}()
}

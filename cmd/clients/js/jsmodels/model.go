package jsmodels

import (
	"net"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
	json "github.com/johanbrandhorst/gopherjs-json"
	help "github.com/juanmachuca95/ahorcado_go/pkg/helpers"
	"honnef.co/go/js/xhr"
)

var WSConn net.Conn

type Model struct {
	*js.Object

	// Login
	LoginUser  bool   `js:"login"`
	Username   string `js:"username"`
	Password   string `js:"password"`
	ErrorLogin string `js:"error_login"`
	Token      string `js:"token"`

	// Register
	RegisterUser  bool   `js:"register"`
	ErrorRegister string `js:"error_register"`

	// Websocket
	ConnOpen bool `js:"ws_conn"`

	// Inputs
	InputWord string `js:"input_word"`
	InputUser string `js:"input_user"`

	// Game
	Jugando      bool    `js:"jugando"`
	FoundLetters string  `js:"found_letters"`
	GameData     []*Game `js:"game_data"`

	// outputs
	Winner string `js:"winner"`
	Status string `js:"status"`
	Error  string `js:"error"`

	Word  *Word `js:"word"`
	Game  *Game `js:"game"`
	Tries int   `js:"tries"`
}

func (m *Model) Login() {
	req := xhr.NewRequest("POST", "http://localhost:8090/api/v1/login")
	req.SetRequestHeader("Content-Type", "application/json")
	login := &RequestLogin{Object: js.Global.Get("Object").New()}
	login.Username = m.Username
	login.Password = m.Password

	s, err := json.Marshal(login.Object)
	if err != nil {
		panic(err)
	}

	go func() {
		err := req.Send([]byte(s))
		if err != nil {
			panic(err)
		}

		rObj, err := json.Unmarshal(req.ResponseText)
		if err != nil {
			m.ErrorLogin = err.Error()
			return
		}

		if req.Status != 200 {
			m.ErrorLogin = rObj.Get("message").String()
			return
		}

		m.Token = rObj.Get("token").String()
		if m.Token != "" {
			m.Username = ""
			m.Password = ""
		}
	}()
}

func (m *Model) Register() {
	req := xhr.NewRequest("POST", "http://localhost:8090/api/v1/register")
	req.SetRequestHeader("Content-Type", "application/json")
	login := &RequestLogin{Object: js.Global.Get("Object").New()}
	login.Username = m.Username
	login.Password = m.Password

	s, err := json.Marshal(login.Object)
	if err != nil {
		panic(err)
	}

	go func() {
		err := req.Send([]byte(s))
		if err != nil {
			panic(err)
		}

		rObj, err := json.Unmarshal(req.ResponseText)
		if err != nil {
			m.ErrorRegister = err.Error()
			return
		}

		if req.Status != 200 {
			m.ErrorRegister = rObj.Get("message").String()
			return
		}

		m.Token = rObj.Get("token").String()
		if m.Token != "" {
			m.Username = ""
			m.Password = ""
			m.GetGame()
		}
	}()
}

func (m *Model) Jugar() {
	if m.Token == "" {
		m.Error = "Ups! Parece que hoy no será posible jugar, lo sentimos."
		return
	}

	m.Jugando = true
	m.GetGame()
}

func (m *Model) GetGame() {
	req := xhr.NewRequest("GET", "http://localhost:8090/api/v1/game")
	req.SetRequestHeader("Content-Type", "application/json")
	req.SetRequestHeader("Authorization", m.Token)

	go func() {
		err := req.Send(nil)
		if err != nil {
			panic(err)
		}

		if req.Status != 200 {
			m.Error = "No hay juegos disponibles."
			return
		}

		rObj, err := json.Unmarshal(req.ResponseText)
		if err != nil {
			m.Error = err.Error()
			return
		}

		msg := &Game{
			Object: rObj,
		}

		m.Game = msg
		m.FoundLetters = help.ShowWord(m.Game.Word, m.Game.Encontrados)
	}()
}

func (m *Model) Connect() {
	go func() {
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

	aux := &struct {
		*js.Object
		msg *Game `js:"result"`
	}{
		Object: rObj,
	}

	if rObj.Get("error").Bool() {
		panic(msg)
	}

	return aux.msg
}

func (m *Model) Send() {
	msg := &Word{Object: js.Global.Get("Object").New()}

	msg.Game_id = m.Game.Id
	msg.Word = strings.ToUpper(m.InputWord)
	msg.User = m.Username

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
			m.Status = game.Status
			m.FoundLetters = help.ShowWord(game.Word, game.Encontrados)
			m.GameData = append(m.GameData, game)
			m.InputWord = ""

			if m.Game.Finalizada {
				m.Winner = "🏆 " + game.Winner + " - Ha ganado el juego"
				m.Reset()
			}

		}

	}()
}

func (m *Model) SetRegister() {
	m.RegisterUser = true
	m.LoginUser = false
}

func (m *Model) SetLogin() {
	m.LoginUser = true
	m.RegisterUser = false
}

func (m *Model) Reset() {
	m.Username = ""

	// Inputs
	m.InputUser = ""
	m.InputWord = ""

	// Game
	m.Status = ""
	m.Error = ""
	m.Winner = ""
	m.Game = &Game{}
	m.Word = &Word{}
	m.ConnOpen = false
	m.FoundLetters = ""
	m.Tries = 6

}
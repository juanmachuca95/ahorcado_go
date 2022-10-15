package game

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/juanmachuca95/ahorcado_go/pkg/frames"
	ah "github.com/juanmachuca95/ahorcado_go/protos/ahorcado"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	PLAY              = "--> ? Join The Game "
	RANTING           = "--> ? See Top Players"
	EXIT              = "--> ? Exit"
	_codeFound        = 1
	_codeNotFound     = 2
	_codeAlreadyFound = 3
	_codeWinner       = 4
	_codeUnexpected   = 5
)

var (
	TRIES = 6
)

type Game interface {
	Init() string
	Start(string)
	PanelOptions()
	SeeRanking()
	Exit()
}

type game struct {
	gameID     string
	conn       *grpc.ClientConn
	client     ah.AhorcadoClient
	stream     ah.Ahorcado_AhorcadoClient
	username   string
	serverAddr string
}

func NewGame(serverAddr, username string) Game {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	cli := ah.NewAhorcadoClient(conn)
	return &game{
		conn:       conn,
		username:   username,
		client:     cli,
		serverAddr: serverAddr,
	}
}

func (g *game) Init() string {
	juego, err := g.GetGame()
	if err != nil {
		panic("error to get game")
	}

	// Show game
	g.ShowWord(juego.Word, juego.Encontrados, len(juego.Word))
	g.PrintInput()
	pterm.Println()
	return juego.GetId()
}

func (g *game) Start(gameID string) {
	ctx := context.Background()
	stream, err := g.client.Ahorcado(ctx)
	if err != nil {
		panic("unexpected error")
	}
	g.stream = stream
	g.gameID = gameID

	// Escribiendo los mensajes recibidos.
	var quit = make(chan int)
	go g.PrintMessageGame(quit)
	g.InputGame(quit)
}

func (g *game) PrintMessageGame(c chan int) {
	for {
		game, err := g.stream.Recv()
		if err != nil {
			return
		}

		status, code := g.MessageStatus(game)
		g.Attempts(code)
		g.ShowInfo(game, status, code)
		if g.checkWin(game, g.username) {
			c <- 1
			close(c)
		}
	}
}

func (g *game) InputGame(c chan int) {
	var input string
	leave := false
	for {
		select {
		case <-c: // win game
			leave = true
			pterm.Info.Println("Game finished ", leave)
			TRIES = 6
			return
		default:
			if leave {
				return
			}

			fmt.Scan(&input)
			input = strings.ToTitle(input)
			if err := g.stream.Send(&ah.Word{GameId: g.gameID, User: g.username, Word: input}); err != nil {
				pterm.Warning.Println("Cannot send message to server - error: ", err.Error())
				return
			}
		}
	}
}

// Get the current game in room
func (g *game) GetGame() (*ah.Game, error) {
	game, err := g.client.GetGame(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal(err.Error())
		return game, err
	}
	return game, nil
}

// Show an status to the current game
func (g *game) ShowInfo(game *ah.Game, status string, code int) {
	pterm.Println()
	// Print main game of the word
	g.ShowWord(game.Word, game.Encontrados, len(game.Word))
	pterm.Info.Println("El usuario", game.Usersend,
		" ha jugado: ", game.Wordsend,
		"\nLetras encontradas: ", game.Encontrados,
		"\nIntentos: ", TRIES,
		"\nEstatus: ", status,
	)
	pterm.Println()
	g.PrintInput()
}

// Helper función to show word
func (g *game) ShowWord(clave string, letras []string, length int) {
	chars := []rune(clave)
	var (
		wordPositions []string
		results       []string
		encontrados   []string = letras
	)
	for i := 0; i < len(chars); i++ {
		char := string(chars[i])
		results = append(results, "_")
		wordPositions = append(wordPositions, char)
	}

	for _, encontrado := range encontrados {
		for i, worPos := range wordPositions {
			if encontrado == worPos {
				results[i] = worPos
			}
		}
	}

	pterm.DefaultSection.Println("Palabra: ", results, " - (", length, ") Letras")
}

func (g *game) checkWin(game *ah.Game, user string) bool {
	pterm.Println()
	if game.Finalizada {
		switch game.Usersend {
		case g.username:
			g.VictoryMessage(game)
			return true
		default:
			g.DefeatMessage(game)
			g.Exit()
			return true
		}
	}

	return false
}

func (g *game) PanelOptions() {
	style := pterm.NewStyle(pterm.FgGreen, pterm.Bold)
	style.Print("\n")
	result, err := pterm.DefaultInteractiveSelect.WithOptions([]string{PLAY, RANTING, EXIT}).
		Show()
	if err != nil {
		panic("unexpected error")
	}

	// Instancia de game
	switch result {
	case PLAY:
		g.Start(g.Init())
		g.PanelOptions()
	case RANTING:
		g.SeeRanking()
		g.PanelOptions()
	case EXIT:
		g.Exit()
	default:
		g.PanelOptions()
	}
}

func (g *game) SeeRanking() {
	ranking, err := g.client.GetTop(context.Background(), &emptypb.Empty{})
	if err != nil {
		g.conn.Close()
		panic("Cannot consulting to server about top players")
	}

	pterm.DefaultSection.Println("Top players Ahorcado 🏆")
	style := pterm.NewStyle(pterm.FgWhite, pterm.Bold)
	for _, rank := range ranking.Rankings {
		style.Printf("\t🏁 %s \t\t%d\n", rank.Username, rank.Won)
	}

	pterm.Println()
}

func (g *game) Exit() {
	g.conn.Close()
	pterm.FgYellow.Println("Thanks for playing my game. ¡Bye!")
	os.Exit(1)
}

func (g *game) PrintInput() {
	// Input after every info showed
	pStyle := pterm.NewStyle(pterm.FgGreen, pterm.Bold)
	pStyle.Print("Letter or word 👇")
}

func (g *game) Attempts(code int) {
	if code > 0 {
		TRIES--
		frames.Frames(TRIES)
		if TRIES == 0 {
			pterm.FgRed.Println("😥 Sorry you lost, better luck next time")
			pterm.Println()
			// informar que ha perdido el game en el servidor
			g.Exit()
		}
	}
}

func (g *game) DefeatMessage(game *ah.Game) {
	pterm.Println()
	pterm.FgYellow.Println("👎 Han estado más rapido que tú, la proxima será. La palabra era: ", game.Word)
	pterm.FgCyan.Println("Ha ganado el usuario: ", game.Usersend, "✅.")
	pterm.Println()
}

func (g *game) VictoryMessage(game *ah.Game) {
	panel := pterm.DefaultBox.WithTitle("🎉 Has Ganado! 🎉").Sprint("\nFelicidades 🏆 ", game.Usersend, "\nLa palabra era: 💀 ", game.Word)
	panels, err := pterm.DefaultPanel.WithPanels(pterm.Panels{{{Data: panel}}}).Srender()
	if err != nil {
		panic("Cannot implement panel pterm")
	}
	pterm.DefaultBox.WithTitle("Ahorcado | Made by @juanmachuca95").WithTitleBottomRight().WithRightPadding(0).WithBottomPadding(0).Println(panels)
}

// game.Usersend, g.username, game.Wordsend, game.Word, int(game.Status)
func (g *game) MessageStatus(game *ah.Game) (string, int) {
	messageMe := false
	if game.Usersend == g.username {
		messageMe = true
	}
	switch game.Status {
	case _codeFound:
		if messageMe {
			return fmt.Sprintf("¡Has encontrado la letra %s! 👍", game.Wordsend), 0
		}
		return fmt.Sprintf("El usario %s ha encontrado la letra %s 🔥", g.username, game.Wordsend), 0
	case _codeNotFound:
		if messageMe {
			return fmt.Sprintf("La letra %s no tiene coincidencias! 👎", game.Wordsend), 1
		}
		return fmt.Sprintf("¡El usuario %s ha intentado la letra %s sin exito! 🖊️", g.username, game.Wordsend), 0
	case _codeAlreadyFound:
		if messageMe {
			return fmt.Sprintf("La letra %s ya ha sido encontrada previamente! 🐢", game.Wordsend), 1
		}
		return fmt.Sprintf("¡El usuario %s ha intentado la letra %s pero ya ha sido encontrada! 🐢", g.username, game.Wordsend), 0
	case _codeWinner:
		if messageMe {
			return fmt.Sprintf("🎉 Has ganado! 🎊 La palabra era %s 🏆", game.Word), 0
		}
		return fmt.Sprintf("¡El usuario %s ha ganado! 🐢 La palabra era: %s", g.username, game.Word), 0
	case _codeUnexpected:
		return "Ha ocurrido un error inesperado", 0
	default:
		return "", 0
	}
}

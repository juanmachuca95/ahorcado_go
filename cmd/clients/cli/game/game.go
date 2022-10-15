package game

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/juanmachuca95/ahorcado_go/pkg/frames"
	"github.com/juanmachuca95/ahorcado_go/pkg/helpers"
	ah "github.com/juanmachuca95/ahorcado_go/protos/ahorcado"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	PLAY    = "--> ? Join The Game "
	RANTING = "--> ? See Top Players"
	EXIT    = "--> ? Exit"
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
	conn       *grpc.ClientConn
	client     ah.AhorcadoClient
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
	// Stream playing
	stream, err := g.client.Ahorcado(ctx)
	if err != nil {
		panic("unexpected error")
	}

	// Escribiendo los mensajes recibidos.
	var quit = make(chan int)
	go func() {
		for {
			game, err := stream.Recv()
			if err != nil {
				return
			}

			g.ShowInfo(game)
			if g.checkWin(game, g.username) {
				quit <- 1
				close(quit)
			}
		}
	}()

	var input string
	var leave bool
	for {
		select {
		case <-quit: // win game
			TRIES = 6
			leave = true
			pterm.Info.Println("Game finished ", leave)
			return
		default:
			if leave {
				if err := stream.CloseSend(); err != nil {
					log.Fatal(err)
				}
				return
			}

			fmt.Scan(&input)
			input = strings.ToTitle(input)
			if err = stream.Send(&ah.Word{GameId: gameID, User: g.username, Word: input}); err != nil {
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
func (g *game) ShowInfo(game *ah.Game) {
	pterm.Println()
	status, code := helpers.MessageStatus(game.Usersend, g.username, game.Wordsend, game.Word, int(game.Status))
	if code > 0 {
		TRIES--
		frames.Frames(TRIES)
		if TRIES == 0 {
			pterm.FgRed.Println("😥 Sorry you lost, better luck next time. you exhausted the number of possible attempts")
			pterm.Println()
			// informar que ha perdido el game en el servidor
			return
		}
	}

	// Print main game of the word
	g.ShowWord(game.Word, game.Encontrados, len(game.Word))
	pterm.Info.Println("El usuario", game.Usersend, " ha jugado: ", game.Wordsend,
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

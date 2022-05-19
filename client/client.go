package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/juanmachuca95/ahorcado_go/frames"
	"github.com/juanmachuca95/ahorcado_go/generated"
	"github.com/pterm/pterm"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:8080", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

var encontrados = []string{}

var collection *mongo.Collection
var ctx = context.TODO()
var user string = "default"

func main() {
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	/* Ahorcado Client gRPC Go */
	client := generated.NewAhorcadoClient(conn)

	Header()
	var leaveGame bool = false
	var myInput string
	var myUser string

	/* Inside Stream */
	var input string
	var finded bool = false
	var inGame generated.Game
	var activateStream bool = false
	var tries int = 6
	pterm.Print("\n")

	for !leaveGame {
		if activateStream {
			waitc := make(chan struct{})
			stream, err := client.Ahorcado(context.Background())
			if err != nil {
				pterm.Warning.Println("Fallo al activar stream del juego - error: ", err.Error())
				activateStream = false
				leaveGame = true
				return
			}

			go func() {
				for {
					in, err := stream.Recv()
					if err == io.EOF {
						close(waitc)
						return
					}
					if err != nil {
						log.Fatalf("Failed to receive a note : %v", err)
					}

					inGame.Id = in.Id
					inGame.Word = in.Word
					inGame.Winner = in.Winner
					inGame.Encontrados = in.Encontrados
					inGame.Finalizada = in.Finalizada
					inGame.Error = in.Error

					if inGame.Finalizada {
						GameWin(&inGame, myUser)
						activateStream = false
						finded = true
						pterm.Println()
						Panel()
					} else {
						pterm.Println()
						pterm.DefaultSection.Println("Encontrados: ", inGame.Encontrados)
						pterm.Info.Println("En juego\nFinalizada: ", inGame.Finalizada, " \nIntentos: ", tries, "\nUsuario ", in.UserSend, "\nFallos: ", inGame.Error)
						if inGame.Error != "" && in.UserSend == myUser {
							tries--
							frames.Frames(tries)
						}

						if tries == 0 {
							pterm.FgRed.Println("😥 Has perdido, lo sentimos. Inténtalo de nuevo más tarde.")
							activateStream = false
							finded = true
						}
						pterm.Println()
					}
				}
			}()

			if activateStream {
				for !finded {
					fmt.Scan(&input)
					input = strings.ToTitle(input)

					req := &generated.Word{
						GameId: inGame.Id,
						Word:   input,
						User:   myUser,
					}
					if err := stream.Send(req); err != nil {
						log.Fatalf("Failed to send a note: %v", err)
					}
				}
			}

			stream.CloseSend()
			<-waitc
		}

		/* Entradas del usuario*/
		Panel()
		pterm.FgDefault.Print("Ingresa una opcion: \n")
		fmt.Scan(&myInput)
		switch myInput {
		case "1":
			pterm.FgYellow.Print("[+] Ingresa un username: ")
			fmt.Scan(&myUser) // Set user
			game, err := GetRandomGame(client)
			if err != nil {
				pterm.Warning.Println("No fue posible obtener un juego activo.")
				leaveGame = true
				return
			}

			inGame.Id = game.Id
			inGame.Word = game.Word
			inGame.Encontrados = game.Encontrados
			inGame.Finalizada = game.Finalizada
			inGame.Error = game.Error
			activateStream = true
			finded = false
			tries = 6

			if inGame.Error != "" {
				pterm.Println()
				pterm.FgLightRed.Println("No hay juegos disponibles")
				pterm.Println()
				return
			}

			pterm.Println()
			pterm.FgGreen.Println("**** El juego ha iniciado 💀 - Tenes 6 intentos ****")
			pterm.Println()
		case "2":
			GetRanking(client)
		default:
			Panel()
		}
	}
}

func Header() {
	// Color: primary  "255, 215, 0"
	header := pterm.DefaultHeader.WithBackgroundStyle(pterm.DefaultHeader.BackgroundStyle)
	pterm.DefaultCenter.Println(header.Sprintf("Made by @juanmachuca95"))

	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithRGB("Ahorcado", pterm.NewRGB(255, 215, 0))).
		Render()
}

func Panel() {
	panels := pterm.Panels{
		{{Data: "Opciones del juego"}},
		{{Data: "[1] Unirse al Game"}, {Data: "\n Ingresa tu nombre de usuario\n Posteriormente se inicia el juego\n Comienza"}},
		/* 		{{Data: "[2] Ranking Podio"}, {Data: "\n Ingresa 2 para ver el Ranking actual del juego y tu número de posición\n"}},
		 */}

	_ = pterm.DefaultPanel.WithPanels(panels).Render()
}

func GameWin(inGame *generated.Game, myUser string) {
	if myUser == inGame.Winner {
		panel := pterm.DefaultBox.WithTitle("🎉 Has Ganado! 🎉").Sprint("\nFelicidades 🏆", inGame.Winner, "\nLa palabra era: 💀 ", inGame.Word)

		panels, _ := pterm.DefaultPanel.WithPanels(pterm.Panels{{{Data: panel}}}).Srender()
		pterm.DefaultBox.WithTitle("Ahorcado | Made by @juanmachuca95").WithTitleBottomRight().WithRightPadding(0).WithBottomPadding(0).Println(panels)
	} else {
		pterm.Println()
		pterm.FgYellow.Println("👎 Han estado más rapido que tú, la proxima será")
		pterm.FgCyan.Println("Ha ganado el usuario: ", inGame.Winner, "✅.")
		pterm.Println()
	}
}

func ShowInfo(inGame *generated.Game, tries int, userSend, wordSend string) {
	pterm.Println()
	chars := []rune(inGame.Word)
	var wordPositions []string

	for i := 0; i < len(chars); i++ {
		char := string(chars[i])
		wordPositions = append(wordPositions, char)
	}

	fmt.Println(wordPositions)

	pterm.DefaultSection.Println("Encontrados: ", inGame.Encontrados)
	pterm.Info.Println("En juego\nFinalizada: ", inGame.Finalizada, " \nIntentos: ", tries, "\nUsuario ", userSend, "\nFallos: ", inGame.Error)

}

func GetRandomGame(client generated.AhorcadoClient) (*generated.Game, error) {
	game, err := client.GetRandomGame(context.Background(), &generated.Empty{})
	if err != nil {
		log.Fatal(err.Error())
		return game, err
	}

	return game, nil
}

func GetRanking(client generated.AhorcadoClient) {
	log.Println("Función no programada")
}

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/juanmachuca95/ahorcado_go/pkg/frames"
	ah "github.com/juanmachuca95/ahorcado_go/protos/ahorcado"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	// var tls bool
	// caFile     string
	serverAddr string
	// var serverHostOverride string
)

func init() {
	//tls = *flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	// caFile = *flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr = *flag.String("addr", "localhost:8080", "The server address in the format of host:port")
	//serverHostOverride = *flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
}

func main() {
	flag.Parse()
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	/* Ahorcado Client gRPC Go */
	client := ah.NewAhorcadoClient(conn)

	Header()
	var leaveGame bool
	var myInput string
	var myUser string

	/* Inside Stream */
	var input string
	var finded bool = false
	var inGame ah.Game
	var activateStream bool = false
	var tries int = 6
	pterm.Print("\n")

	for !leaveGame {
		if activateStream {
			pterm.DefaultSection.Println("Palabra: ", ShowWord(inGame.Word, inGame.Encontrados))

			waitc := make(chan struct{})
			stream, err := client.Ahorcado(context.Background())
			if err != nil {
				pterm.Warning.Println("Fallo al activar stream del juego - error: ", err.Error())
				activateStream = false
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
					inGame.Wordsend = in.Wordsend
					inGame.Usersend = in.Usersend

					if inGame.Finalizada {
						GameWin(&inGame, myUser)
						activateStream = false
						finded = true
						input = ""
						pterm.Println()
						Panel()
					} else {
						ShowInfo(&inGame, tries)
						if inGame.Error != "" && in.Usersend == myUser {
							tries--
							frames.Frames(tries)
						}

						if tries == 0 {
							pterm.FgRed.Println("???? Has perdido, lo sentimos. Int??ntalo de nuevo m??s tarde.")
							activateStream = false
							finded = true
							input = ""
						}
						pterm.Println()
					}
				}
			}()

			if activateStream {
				for !finded {
					fmt.Scan(&input)
					input = strings.ToTitle(input)

					req := &ah.Word{
						GameId: inGame.Id,
						Word:   input,
						User:   myUser,
					}
					if err := stream.Send(req); err != nil {
						log.Fatalf("Failed to send a note: %v", err)
					}
				}
			}

			err = stream.CloseSend()
			if err != nil {
				panic(err)
			}
			<-waitc
		}

		/* Entradas del usuario*/
		switch myInput {
		case "1":
			pterm.FgYellow.Print("[+] Ingresa un username: ")
			fmt.Scan(&myUser) // Set user
			game, err := GetRandomGame(client)
			if err != nil {
				pterm.Warning.Println("No fue posible obtener un juego activo.")
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
			pterm.FgGreen.Println("**** El juego ha iniciado ???? - Tenes 6 intentos ****")
			pterm.Println()
		case "2":
			GetRanking(client)
		default:
			Panel()
			pterm.FgDefault.Print("Ingresa una opcion: \n")
			fmt.Scan(&myInput)
		}
	}
}

func Header() {
	// Color: primary  "255, 215, 0"
	header := pterm.DefaultHeader.WithBackgroundStyle(pterm.DefaultHeader.BackgroundStyle)
	pterm.DefaultCenter.Println(header.Sprintf("Made by @juanmachuca95"))

	err := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromStringWithRGB("Ahorcado", pterm.NewRGB(255, 215, 0))).Render()
	if err != nil {
		log.Fatal("Cannot render information print")
	}
}

func Panel() {
	panels := pterm.Panels{
		{{Data: "Opciones del juego"}},
		{{Data: "[1] Unirse al Game"}, {Data: "\n Ingresa tu nombre de usuario\n Posteriormente se inicia el juego\n Comienza"}},
		/* 		{{Data: "[2] Ranking Podio"}, {Data: "\n Ingresa 2 para ver el Ranking actual del juego y tu n??mero de posici??n\n"}},
		 */}

	_ = pterm.DefaultPanel.WithPanels(panels).Render()
}

func GameWin(inGame *ah.Game, myUser string) {
	if myUser == inGame.Winner {
		panel := pterm.DefaultBox.WithTitle("???? Has Ganado! ????").Sprint("\nFelicidades ????", inGame.Winner, "\nLa palabra era: ???? ", inGame.Word)

		panels, _ := pterm.DefaultPanel.WithPanels(pterm.Panels{{{Data: panel}}}).Srender()
		pterm.DefaultBox.WithTitle("Ahorcado | Made by @juanmachuca95").WithTitleBottomRight().WithRightPadding(0).WithBottomPadding(0).Println(panels)
	} else {
		pterm.Println()
		pterm.FgYellow.Println("???? Han estado m??s rapido que t??, la proxima ser??. La palabra era: ", inGame.Word)
		pterm.FgCyan.Println("Ha ganado el usuario: ", inGame.Winner, "???.")
		pterm.Println()
	}
}

func ShowInfo(inGame *ah.Game, tries int) {
	pterm.Println()

	var status string
	if inGame.Error != "" {
		status = inGame.Error
	} else {
		status = fmt.Sprintf("La letra (%s) Ha sido encontrada ????", inGame.Wordsend)
	}

	results := ShowWord(inGame.Word, inGame.Encontrados)
	pterm.DefaultSection.Println("Palabra: ", results, " - (", len(inGame.Word), ") Letras")
	pterm.Info.Println("El usuario", inGame.Usersend, " ha jugado: ", inGame.Wordsend, "\nLetras encontradas: ", inGame.Encontrados, " \nIntentos: ", tries, "\nEstatus: ", status)
	pterm.Println()
}

func GetRandomGame(client ah.AhorcadoClient) (*ah.Game, error) {
	game, err := client.GetGame(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal(err.Error())
		return game, err
	}

	return game, nil
}

func GetRanking(client ah.AhorcadoClient) {
	log.Println("Funci??n no programada")
}

func ShowWord(clave string, letras []string) []string {
	chars := []rune(clave)
	var wordPositions []string

	var encontrados []string = letras
	var results []string
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

	return results
}

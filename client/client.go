package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/juanmachuca95/ahorcado_go/generated"
	"github.com/juanmachuca95/ahorcado_go/services/game/models"
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
	Panel()
	var leaveGame bool = false
	var myInput string
	var inGame generated.Game
	var activateStream bool = false
	var mux sync.Mutex
	pterm.Print("\n")
	waitc := make(chan struct{})
	for !leaveGame {

		if activateStream {
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

					mux.Lock()
					inGame.Id = in.Id
					inGame.Word = in.Word
					inGame.Winner = in.Winner
					inGame.Encontrados = in.Encontrados
					inGame.Finalizada = in.Finalizada
					inGame.Error = in.Error

					if in.Finalizada {
						log.Println("🏆 HA GANADO EL USUARIO: ", in.Winner, " - AL DESCUBRIR LA PALABRA: ", in.Word)
						leaveGame = true
						activateStream = false
					}

					pterm.Println()
					pterm.DefaultSection.Println("Encontrados: ", in.Encontrados)
					pterm.Info.Println("En juego\nGanador: ", in.Winner)
					pterm.Println()
					mux.Unlock()
				}
			}()

			var input string
			var finded bool = false
			var failed bool = false
			for !finded && !failed {
				fmt.Scan(&input)
				input = strings.ToTitle(input)

				req := &generated.Word{
					GameId: inGame.Id,
					Word:   input,
					User:   user,
				}
				if err := stream.Send(req); err != nil {
					log.Fatalf("Failed to send a note: %v", err)
				}

			}

			stream.CloseSend()
			<-waitc
		}

		/* Entradas del usuario*/
		fmt.Scan(&myInput)
		switch myInput {
		case "1":
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
			activateStream = true

			pterm.Println()
			pterm.FgBlue.Println("ID DEL JUEGO ES: ", inGame.Id)
			pterm.FgYellow.Println("Letras encontradas: ", inGame.GetEncontrados(), " PALABRA: ", inGame.GetWord())
			pterm.FgGreen.Println("Juego: ACTIVO - Stream: ", activateStream, " Finalizado: ", inGame.GetFinalizada())
		case "2":
			GetRanking(client)
		default:
			Panel()
		}
	}

	// Declare panels in a two dimensional grid system.
}

func getGame() models.Game {
	return models.Game{}
}

func win(clave string) bool {
	var fin bool = false
	var lengthClave int = len(clave)
	var lengthEncontrados int = 0
	for _, encontrado := range encontrados {
		count := strings.Count(clave, encontrado)
		lengthEncontrados += count
	}

	log.Println("Cantidad de encontrados = ", lengthEncontrados, " Cantidad total de la clave = ", lengthClave)
	if lengthEncontrados == lengthClave {
		fin = true
	}
	return fin
}

func alreadyFound(character string) bool {
	var result bool = false
	for _, encontrado := range encontrados {
		if character == encontrado {
			result = true
		}
	}
	return result
}

/* if clave == input {
	log.Println("Arriesgaste", input, " HAS GANADO 🏆 - coincidencias: TOTAL - palabra: ", clave)
	finded = true
} else if len(input) > 1 {
	tries--
	frames.Frames(tries)
	log.Println("Encontrados hasta el momento: ", encontrados)
	log.Println("Arriesgaste", input, "HAS FALLADO 👎 - coincidencias: 0 - Intentos: ", tries)
} else if alreadyFound(input) {
	log.Println("El caracter", input, "YA HA SIDO ENCONTRADO - coincidencias: ", strings.Count(clave, input))
} else if strings.Contains(clave, input) {
	log.Println("El caracter", input, " SI esta 👍 - coincidencias: ", strings.Count(clave, input))
	encontrados = append(encontrados, input)
	log.Println("Encontrados: ", encontrados)
	if win(clave) {
		finded = true
		log.Println("🏆 Has ganado el juego ", encontrados)
	}
} else {
	tries--
	frames.Frames(tries)
	log.Println("Encontrados hasta el momento: ", encontrados)
	log.Println("El caracter", input, " (NO) esta 👎 - coincidencias: ", strings.Count(clave, input), " - Intentos: ", tries)
}

if tries == 0 {
	failed = true

	log.Println("Lo siento has perdido. ")
}
*/

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
		{{Data: "[+] Unirse al juego"}, {Data: "\n Ingresa 1 para comenzar a juagar\n"}},
		{{Data: "[*] Ranking Podio"}, {Data: "\n Ingresa 2 para ver el Ranking actual del juego y tu número de posición\n"}},
	}

	_ = pterm.DefaultPanel.WithPanels(panels).Render()
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

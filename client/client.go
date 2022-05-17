package main

import (
	"context"
	"flag"
	"log"
	"strings"

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

	/* Titúlo */
	log.Println("Death By Hanging - Go 👍")

	/* Ahorcado Client gRPC Go */
	client := generated.NewAhorcadoClient(conn)
	game, err := client.GetRandomGame(context.Background(), &generated.Empty{})

	if game.Error != "" {
		log.Println("🔥 - ", game.Error)
		return
	}

	// Color
	// primary  "255, 215, 0"
	header := pterm.DefaultHeader.WithBackgroundStyle(pterm.DefaultHeader.BackgroundStyle)
	pterm.DefaultCenter.Println(header.Sprintf("Made by @juanmachuca95"))

	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithRGB("Ahorcado", pterm.NewRGB(255, 215, 0))).
		Render()
		// Declare panels in a two dimensional grid system.
	panels := pterm.Panels{
		{{Data: "[+] Empezar juego"}, {Data: "\n Ingresa 1 para comenzar a juagar\n"}},
		{{Data: "[*] Ranking Podio"}, {Data: "\n Ingresa 2 para ver el Ranking actual del juego y tu número de posición\n"}},
	}

	// Print panels.
	_ = pterm.DefaultPanel.WithPanels(panels).Render()
	pterm.NewRGB(255, 215, 0).Println("GAME ID: ", game.Id, " - WORD: ", game.Word)
	log.Println("*******************************")

	/* stream, err := client.Ahorcado(context.Background())
	if err != nil {
		log.Fatal(err)
	} */

	/*
		waitc := make(chan struct{})
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

				log.Println("***************************************")
				log.Println("Encontrados", in.Encontrados)
				log.Println("***************************************")

				if in.Finalizada {
					log.Println("🏆 HA GANADO EL USUARIO: ", in.Winner, " - AL DESCUBRIR LA PALABRA: ", in.Word)
				}
			}
		}()

		log.Println("Death By Hanging - Juan Gabriel Machuca")

		clave := "Developer"
		clave = strings.ToTitle(clave)

		var input string
		var finded bool = false
		var failed bool = false

		for !finded && !failed {
			fmt.Scan(&input)
			input = strings.ToTitle(input)

			req := &generated.Word{Word: input, User: user}
			if err := stream.Send(req); err != nil {
				log.Fatalf("Failed to send a note: %v", err)
			}

		}

		stream.CloseSend()
		<-waitc */

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

package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/juanmachuca95/ahorcado_go/generated"
	"github.com/juanmachuca95/ahorcado_go/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var encontrados = []string{}

func main() {
	addr := fmt.Sprintf("0.0.0.0:%d", 8080)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}
	log.Println("The server is running successfully on port ", 8080)

	ahorcado := server.NewAhorcadoServer()

	var opts []grpc.ServerOption
	serv := grpc.NewServer(opts...)

	/* Registro de servicios */
	generated.RegisterAhorcadoServer(serv, ahorcado) // Register Services Cliente

	/* Enable reflection */
	reflection.Register(serv)
	if err = serv.Serve(listener); err != nil {
		panic("cannot initialize the server" + err.Error())
	}

	/* log.Println("Death By Hanging - Juan Gabriel Machuca")

	var tries int = 6
	clave := "Developer"
	clave = strings.ToTitle(clave)

	var input string
	var finded bool = false
	var failed bool = false

	for !finded && !failed {
		fmt.Scan(&input)
		input = strings.ToTitle(input)

		if clave == input {
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
	} */
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

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/juanmachuca95/ahorcado_go/frames"
)

var encontrados = []string{}

func main() {
	log.Println("MyGo - Game 'Death By Hanging'")

	var tries int = 6
	clave := "Laura"
	clave = strings.ToTitle(clave)

	var input string
	var finded bool = false
	var failed bool = false

	for !finded && !failed {
		fmt.Scan(&input)
		input = strings.ToTitle(input)

		if alreadyFound(input) {
			log.Println("El caracter", input, " YA HA SIDO ENCONTRADO - coincidencias: ", strings.Count(clave, input))
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
	}

	/* for i := 0; i < tries; i++ {
		frames.Inicial(i)
		frames.FirstAttempt(i)
		frames.SecondAttempt(i)
		frames.ThirdAttempt(i)
		frames.FourthAttempt(i)
		frames.FifthAttempt(i)
		frames.SixthAttempt(i)
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

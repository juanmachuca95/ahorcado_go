package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	log.Println("MyGo - Game 'death by hanging'")

	//var tries int = 8
	clave := "Aprendiendo con Golang"
	clave = strings.ToTitle(clave)

	log.Println("Find: ", clave)
	var input string
	var finded bool = false

	for !finded {
		fmt.Scan(&input)
		input = strings.ToTitle(input)

		if strings.Contains(clave, input) {
			log.Println("El caracter", input, " SI esta 👍 - coincidencias: ", strings.Count(clave, input))
		} else {
			log.Println("El caracter", input, " (NO) esta 👎 - coincidencias: ", strings.Count(clave, input))
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

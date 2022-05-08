package main

import (
	"log"

	frames "github.com/juanmachuca95/ahorcado_go/frames"
)

func main() {
	log.Println("Ahorcado Golang")

	var tries int = 8

	for i := 0; i < tries; i++ {
		frames.Inicial(i)
		frames.FirstAttempt(i)
		frames.SecondAttempt(i)
		frames.ThirdAttempt(i)
		frames.FourthAttempt(i)
		frames.FifthAttempt(i)
		frames.SixthAttempt(i)
	}
}

package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/juanmachuca95/ahorcado_go/generated"
	service "github.com/juanmachuca95/ahorcado_go/services/game/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	/* Mis variables de entorno */
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("SERVER_PORT"))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}
	log.Println("The server is running successfully on port ", os.Getenv("SERVER_PORT"))

	ahorcado := service.NewAhorcadoServer()

	var opts []grpc.ServerOption
	serv := grpc.NewServer(opts...)

	/* Registro de servicios */
	generated.RegisterAhorcadoServer(serv, ahorcado) // Register Services Cliente

	/* Enable reflection */
	reflection.Register(serv)
	if err = serv.Serve(listener); err != nil {
		panic("cannot initialize the server" + err.Error())
	}
}

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/juanmachuca95/ahorcado_go/generated"
	"github.com/juanmachuca95/ahorcado_go/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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
}

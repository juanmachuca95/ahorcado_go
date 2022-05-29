package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	gmService "github.com/juanmachuca95/ahorcado_go/game/handler"
	database "github.com/juanmachuca95/ahorcado_go/internal/database/mongo"
	ah "github.com/juanmachuca95/ahorcado_go/protos/ahorcado"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	/* Mis variables de entorno */
	LoadEnv()

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("SERVER_PORT"))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	// Database
	db := database.Connect()

	// Service
	game := gmService.NewGameService(db)
	var opts []grpc.ServerOption
	serv := grpc.NewServer(opts...)

	/* Registro de servicios */
	ah.RegisterAhorcadoServer(serv, game) // Register Services Cliente

	/* Enable reflection */
	reflection.Register(serv)
	// Serve gRPC server

	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(serv.Serve(listener))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()

	// Register AhorcadoHandler
	err = ah.RegisterAhorcadoHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: wsproxy.WebsocketProxy(gwmux),
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

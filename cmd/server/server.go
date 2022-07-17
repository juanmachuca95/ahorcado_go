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
	"github.com/juanmachuca95/ahorcado_go/game/handler"
	database "github.com/juanmachuca95/ahorcado_go/pkg/database/mongo"
	"github.com/juanmachuca95/ahorcado_go/pkg/interceptor"
	ah "github.com/juanmachuca95/ahorcado_go/protos/ahorcado"
	au "github.com/juanmachuca95/ahorcado_go/protos/auth"

	"github.com/rs/cors"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	MethodOverrideParam = "method"
	TokenCookieName     = "token"
)

func main() {
	LoadEnv()

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("SERVER_PORT"))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic("cannot create tcp connection" + err.Error())
	}

	// Database
	db := database.Connect()

	// Services
	authServ := handler.NewAuthService(db)
	gameServ := handler.NewGameService(db)

	// Middleware
	authInterceptor := interceptor.NewAuthInterceptor()
	serv := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor()),
		grpc.StreamInterceptor(authInterceptor.StreamInterceptor()),
	)

	// Registro de servicios
	ah.RegisterAhorcadoServer(serv, gameServ)
	au.RegisterAuthServer(serv, authServ)

	// Enable reflection
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
	if err := ah.RegisterAhorcadoHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	if err := au.RegisterAuthHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	handler := cors.AllowAll().Handler(wsproxy.WebsocketProxy(gwmux))
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: handler,
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

package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/juanmachuca95/ahorcado_go/game/handler"
	database "github.com/juanmachuca95/ahorcado_go/pkg/database/mongo"
	"github.com/juanmachuca95/ahorcado_go/pkg/interceptor"
	ah "github.com/juanmachuca95/ahorcado_go/protos/ahorcado"
	au "github.com/juanmachuca95/ahorcado_go/protos/auth"
	"github.com/rs/cors"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	port string
)

func main() {
	LoadEnv()
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := Start(); err != nil {
		log.Fatal(err)
	}
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			fmt.Println("GRPC")
			grpcServer.ServeHTTP(w, r)
		} else {
			fmt.Println("REST")
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func makeGrpcServer(address string) (*grpc.ClientConn, *grpc.Server) {
	authInterceptor := interceptor.NewAuthInterceptor()
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor()),
		grpc.StreamInterceptor(authInterceptor.StreamInterceptor()),
		grpc.Creds(insecure.NewCredentials()),
	}

	// Database
	grpcServer := grpc.NewServer(opts...)
	db := database.Connect()
	// Services
	authServ := handler.NewAuthService(db)
	gameServ := handler.NewGameService(db)

	// Registro de servicios
	ah.RegisterAhorcadoServer(grpcServer, gameServ)
	au.RegisterAuthServer(grpcServer, authServ)

	conn, err := grpc.DialContext(
		context.Background(),
		address,
		//grpc.WithTransportCredentials(dcreds)
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Println("Failed to dial server", err.Error())
	}

	// Enable reflection
	reflection.Register(grpcServer)
	return conn, grpcServer
}

func makeHttpServer(conn *grpc.ClientConn) *runtime.ServeMux {
	gwmux := runtime.NewServeMux()
	if err := ah.RegisterAhorcadoHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	if err := au.RegisterAuthHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	return gwmux
}

func Start() error {
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	_, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err.Error())
	}

	conn, grpcServer := makeGrpcServer(addr)
	router := makeHttpServer(conn)

	log.Println("Starting server on addr : " + addr)
	handler := cors.Default().Handler(wsproxy.WebsocketProxy(router))

	err = http.ListenAndServe(":"+port, grpcHandlerFunc(grpcServer, handler))
	return err
}

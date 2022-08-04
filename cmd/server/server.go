package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func LoadTLSCredentials() (credentials.TransportCredentials, error) {
	// load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("./cert/server-cert.pem", "./cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	LoadEnv()
	// context
	ctx := context.Background()

	// Database
	db := database.Connect()
	// Services
	authServ := handler.NewAuthService(db)
	gameServ := handler.NewGameService(db)

	// load TLS credentials
	tlsCredentials, err := LoadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	// Middleware
	authInterceptor := interceptor.NewAuthInterceptor()
	servGrpc := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor()),
		grpc.StreamInterceptor(authInterceptor.StreamInterceptor()),
	)

	// Registro de servicios
	ah.RegisterAhorcadoServer(servGrpc, gameServ)
	au.RegisterAuthServer(servGrpc, authServ)

	// Enable reflection
	reflection.Register(servGrpc)

	conn, err := grpc.DialContext(ctx, "0.0.0.0:8080", grpc.WithTransportCredentials(tlsCredentials))
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
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("Serving gRPC & gRPC-Gateway on 0.0.0.0:8080")
	log.Fatalln(http.ListenAndServe(":8080", grpcHandlerFunc(servGrpc, gwServer.Handler)))
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
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

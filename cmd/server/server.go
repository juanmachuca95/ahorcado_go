package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/juanmachuca95/ahorcado_go/cmd/server/utils"
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
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

var log grpclog.LoggerV2
var (
	cert string = "cert/server-cert.pem"
	key  string = "cert/server-key.pem"
)

func init() {
	log = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)
}

func main() {
	LoadEnv()
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
	certPool, err := utils.GetCertPool("cert/ca-cert.pem")
	if err != nil {
		panic(err)
	}

	serverCert, err := utils.ServerCert()
	if err != nil {
		panic(err)
	}

	configTLS := utils.ConfigTLS(*serverCert, certPool)
	tlsCredentials := credentials.NewTLS(configTLS)

	authInterceptor := interceptor.NewAuthInterceptor()
	opts := []grpc.ServerOption{
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor()),
		grpc.StreamInterceptor(authInterceptor.StreamInterceptor()),
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

	addr := fmt.Sprintf("localhost:%d", 8080)
	dcreds := credentials.NewTLS(&tls.Config{
		ServerName: addr,
		RootCAs:    certPool,
	})
	conn, err := grpc.DialContext(
		context.Background(),
		address,
		grpc.WithTransportCredentials(dcreds),
	)
	if err != nil {
		log.Error("Failed to dial server", err.Error())
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
	addr := fmt.Sprintf("localhost:%d", 8080)

	_, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to listen: %v", err.Error())
	}

	conn, grpcServer := makeGrpcServer(addr)
	router := makeHttpServer(conn)

	log.Info("Starting server on addr : " + addr)
	certPool, err := utils.GetCertPool("cert/ca-cert.pem")
	if err != nil {
		panic(err)
	}

	serverCert, err := utils.ServerCert()
	if err != nil {
		panic(err)
	}

	configTLS := utils.ConfigTLS(*serverCert, certPool)
	configTLS.NextProtos = []string{"h2"}
	handler := cors.Default().Handler(wsproxy.WebsocketProxy(router))

	err = http.ListenAndServeTLS(":8080", cert, key, grpcHandlerFunc(grpcServer, handler))
	return err
}

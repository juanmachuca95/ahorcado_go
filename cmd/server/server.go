package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/juanmachuca95/ahorcado_go/game/handler"
	database "github.com/juanmachuca95/ahorcado_go/pkg/database/mongo"
	"github.com/juanmachuca95/ahorcado_go/pkg/interceptor"
	ah "github.com/juanmachuca95/ahorcado_go/protos/ahorcado"
	au "github.com/juanmachuca95/ahorcado_go/protos/auth"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

var (
	KeyPair  *tls.Certificate
	CertPool *x509.CertPool
	logger   *logrus.Logger
)

func init() {
	logger = logrus.StandardLogger()
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.Kitchen,
		DisableSorting:  true,
	})
	grpclog.SetLogger(logger)
}

func LoadTLSCredentials() (credentials.TransportCredentials, error) {
	// load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	// read ca's cert, verify to client's certificate
	caPem, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	// create cert pool and append ca's cert
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		log.Fatal(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
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

	// GRPC Server
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
	// conn, err := grpc.DialContext(ctx, "0.0.0.0:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	// REST Server
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
	log.Fatalln(http.ListenAndServeTLS(":8080", "cert/server-cert.pem", "cert/server-key.pem", grpcHandlerFunc(servGrpc, gwServer.Handler)))
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			log.Println("GRPC")
			grpcServer.ServeHTTP(w, r)
		} else {
			log.Println("REST")
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

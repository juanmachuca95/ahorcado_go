package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/juanmachuca95/ahorcado_go/generated"
	"github.com/juanmachuca95/ahorcado_go/services/game/models"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:8080", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

var collection *mongo.Collection
var ctx = context.TODO()

func main() {
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	/* Titúlo */
	log.Println("Death By Hanging - Go 👍")

	/* Ahorcado Client gRPC Go */
	client := generated.NewAhorcadoClient(conn)
	stream, err := client.Ahorcado(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Printf("Got message %s)", in.Word)
		}
	}()

	var input string
	for {
		fmt.Scan(&input)
		req := &generated.Word{Word: input}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}

	stream.CloseSend()
	<-waitc

}

func getGame() models.Game {
	return models.Game{}
}

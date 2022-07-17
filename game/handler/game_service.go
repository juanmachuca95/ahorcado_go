package handler

import (
	"context"
	"log"
	"sync"

	gm "github.com/juanmachuca95/ahorcado_go/game/gateway"
	ah "github.com/juanmachuca95/ahorcado_go/protos/ahorcado"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GameService struct {
	ah.UnimplementedAhorcadoServer
	gmGtw       gm.GameGateway
	connections []*Connection
	broadcast   chan *ah.Word
	quit        chan struct{}
	connLock    sync.Mutex
}

func NewGameService(db *mongo.Database) *GameService {
	game := &GameService{
		gmGtw:     gm.NewGameGateway(db),
		broadcast: make(chan *ah.Word),
		quit:      make(chan struct{}),
	}

	go game.Start()
	return game
}

func (g *GameService) Start() {
	running := true
	var msgErr string = ""
	for running {
		select {
		case msg := <-g.broadcast:
			g.connLock.Lock()
			game, err := g.gmGtw.InGame(msg.Word, msg.User, msg.GameId) // GetGame Service
			if err != nil {
				msgErr = err.Error()
			} else {
				msgErr = ""
			}

			gameToSend := ah.Game{
				Id:          game.Id.Hex(),
				Word:        game.Word,
				Winner:      game.Winner,
				Encontrados: game.Encontrados,
				Finalizada:  game.Finalizada,
				Usersend:    msg.User,
				Wordsend:    msg.Word,
				Status:      game.Status,
				Error:       msgErr,
			}

			log.Println("data user ", gameToSend.Usersend, gameToSend.Wordsend)
			log.Println("code ", gameToSend.Status)
			log.Println(gameToSend.Word, gameToSend.Encontrados)
			for _, v := range g.connections {
				log.Println("Enviando . . . ")
				go v.Send(&gameToSend) // Usuario que lo envia
			}
			g.connLock.Unlock()
		case <-g.quit:
			running = false
		}
	}
}

func (g *GameService) Ahorcado(stream ah.Ahorcado_AhorcadoServer) error {
	conn := NewConnectionGame(stream)

	g.connLock.Lock()
	g.connections = append(g.connections, conn)
	g.connLock.Unlock()

	err := conn.GetMessages(g.broadcast)
	if err != nil {
		log.Println(err.Error())
	}

	g.connLock.Lock()
	for i, v := range g.connections {
		if v == conn {
			g.connections = append(g.connections[:i], g.connections[i+1:]...)
		}
	}
	g.connLock.Unlock()

	return err
}

func (g *GameService) GetRandomGame(ctx context.Context, req *emptypb.Empty) (*ah.Game, error) {
	game, err := g.gmGtw.GetGame()
	if err != nil {
		return nil, err
	}

	return &ah.Game{
		Id:          game.Id.Hex(),
		Word:        game.Word,
		Winner:      game.Winner,
		Encontrados: game.Encontrados,
		Finalizada:  game.Finalizada,
		Status:      game.Status,
	}, nil
}

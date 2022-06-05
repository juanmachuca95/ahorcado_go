package game

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
	var errMsg string
	running := true
	for running {
		select {
		case msg := <-g.broadcast:
			g.connLock.Lock()
			game, err := g.gmGtw.InGame(msg.Word, msg.User, msg.GameId) // GetGame Service
			if err != nil {
				errMsg = err.Error()
			}

			gameToSend := ah.Game{
				Id:          game.Id.Hex(),
				Word:        game.Word,
				Winner:      game.Winner,
				Encontrados: game.Encontrados,
				Finalizada:  game.Finalizada,
				UserSend:    msg.User,
				WordSend:    msg.Word,
				Error:       errMsg,
			}
			errMsg = ""
			for _, v := range g.connections {
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
	log.Print("\nAhorcado\n")
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
	game, err := g.gmGtw.GetRandomGame()
	if err != nil {
		return &ah.Game{}, err
	}

	log.Println("Holas", game)
	return &ah.Game{
		Id:          game.Id.Hex(),
		Word:        game.Word,
		Winner:      game.Winner,
		Encontrados: game.Encontrados,
		Finalizada:  game.Finalizada,
	}, nil
}

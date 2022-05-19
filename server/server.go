package server

import (
	"context"
	"io"
	"log"
	"sync"

	"github.com/juanmachuca95/ahorcado_go/generated"
	game "github.com/juanmachuca95/ahorcado_go/services/game/gateway"
)

type Connection struct {
	conn generated.Ahorcado_AhorcadoServer
	send chan *generated.Game
	quit chan struct{}
}

func NewConnection(conn generated.Ahorcado_AhorcadoServer) *Connection {
	c := &Connection{
		conn: conn,
		send: make(chan *generated.Game),
		quit: make(chan struct{}),
	}
	go c.start()
	return c
}

func (c *Connection) Close() error {
	close(c.quit)
	close(c.send)
	return nil
}

func (c *Connection) Send(msg *generated.Game) {
	defer func() {
		// Ignore any errors about sending on a closed channel
		recover()
	}()
	c.send <- msg
}

func (c *Connection) start() {
	running := true
	for running {
		select {
		case msg := <-c.send:
			c.conn.Send(msg) // Ignoring the error, they just don't get this message.
		case <-c.quit:
			running = false
		}
	}
}

func (c *Connection) GetMessages(broadcast chan<- *generated.Word) error {
	for {
		msg, err := c.conn.Recv()
		if err == io.EOF {
			c.Close()
			return nil
		} else if err != nil {
			c.Close()
			return err
		}
		go func(msg *generated.Word) {
			select {
			case broadcast <- msg:
			case <-c.quit:
			}
		}(msg)
	}
}

type AhorcadoServer struct {
	generated.AhorcadoServer
	broadcast   chan *generated.Word
	quit        chan struct{}
	connections []*Connection
	connLock    sync.Mutex
	gameService game.GameGateway
}

func NewAhorcadoServer() *AhorcadoServer {
	srv := &AhorcadoServer{
		broadcast:   make(chan *generated.Word),
		quit:        make(chan struct{}),
		gameService: game.NewGameGateway(),
	}

	go srv.start()
	return srv
}

func (c *AhorcadoServer) Close() error {
	close(c.quit)
	return nil
}

func (c *AhorcadoServer) start() {
	running := true
	for running {
		select {
		case msg := <-c.broadcast:
			c.connLock.Lock()

			/* GetGame Service*/
			game, err := c.gameService.MyGame(msg)
			if err != nil {
				log.Fatal(err)
			}

			/* Usuario que lo envia*/
			game.UserSend = msg.User
			game.WordSend = msg.Word
			for _, v := range c.connections {
				go v.Send(&game)
			}
			c.connLock.Unlock()
		case <-c.quit:
			running = false
		}
	}
}

func (c *AhorcadoServer) Ahorcado(stream generated.Ahorcado_AhorcadoServer) error {
	conn := NewConnection(stream)

	c.connLock.Lock()
	c.connections = append(c.connections, conn)
	c.connLock.Unlock()

	err := conn.GetMessages(c.broadcast)

	c.connLock.Lock()
	for i, v := range c.connections {
		if v == conn {
			c.connections = append(c.connections[:i], c.connections[i+1:]...)
		}
	}
	c.connLock.Unlock()

	return err
}

func (c *AhorcadoServer) GetRandomGame(context context.Context, req *generated.Empty) (*generated.Game, error) {
	game, err := c.gameService.GetRandomGame()
	if err != nil {
		log.Fatal(err.Error())
	}

	return &game, nil
}

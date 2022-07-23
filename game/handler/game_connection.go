package handler

import (
	"io"
	"log"

	ah "github.com/juanmachuca95/ahorcado_go/protos/ahorcado"
)

type Connection struct {
	conn ah.Ahorcado_AhorcadoServer
	send chan *ah.Game
	quit chan struct{}
}

func NewConnectionGame(conn ah.Ahorcado_AhorcadoServer) *Connection {
	c := Connection{
		conn: conn,
		send: make(chan *ah.Game),
		quit: make(chan struct{}),
	}

	go c.start()
	return &c
}

func (c *Connection) Close() error {
	close(c.quit)
	close(c.send)
	return nil
}

func (c *Connection) Send(msg *ah.Game) {
	defer func() {
		// Ignore any errors about sending on a closed channel
		err := recover()
		if err != nil {
			log.Println("Cannot send message to connection closed")
		}
	}()
	c.send <- msg
}

func (c *Connection) start() {
	running := true
	for running {
		select {
		case msg := <-c.send:
			err := c.conn.Send(msg) // Ignoring the error, they just don't get this message.
			if err != nil {
				log.Println(err)
			}
		case <-c.quit:
			running = false
		}
	}
}

func (c *Connection) GetMessages(broadcast chan<- *ah.Word) error {
	for {
		msg, err := c.conn.Recv()
		if err == io.EOF {
			c.Close()
			return nil
		} else if err != nil {
			c.Close()
			return err
		}
		go func(msg *ah.Word) {
			select {
			case broadcast <- msg:
			case <-c.quit:
			}
		}(msg)
	}
}

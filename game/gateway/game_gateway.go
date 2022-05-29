package game

import (
	m "github.com/juanmachuca95/ahorcado_go/game/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameGateway interface {
	GetRandomGame() (*m.Game, error)
	InGame(word, user, id string) (*m.Game, error)
	/* GetRandomGameToSet() (generated.Game, error)
	CreateGame() (models.Game, error)
	MyGame(*generated.Word) (generated.Game, error)
	SeedWords() bool
	UpdateGame() (bool, error) */
}

type GameInDB struct {
	GameStorage
}

func NewGameGateway(db *mongo.Client) GameGateway {
	return &GameInDB{NewGameStorageGateway(db)}
}

func (g *GameInDB) GetRandomGame() (*m.Game, error) {
	return g.getRandomGame()
}

func (g *GameInDB) InGame(word, user, id string) (*m.Game, error) {
	return g.inGame(word, user, id)
}

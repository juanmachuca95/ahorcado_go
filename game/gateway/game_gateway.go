package gateway

import (
	m "github.com/juanmachuca95/ahorcado_go/game/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameGateway interface {
	CreateGames() error
	GetGame() (*m.Game, error)
	InGame(word, user, id string) (*m.Game, error)
}

type GameInDB struct {
	GameStorage
}

func NewGameGateway(db *mongo.Database) GameGateway {
	return &GameInDB{NewGameStorageGateway(db)}
}

func (g *GameInDB) CreateGames() error {
	return g.createGames()
}

func (g *GameInDB) GetGame() (*m.Game, error) {
	game, err := g.getGame()
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	if err == mongo.ErrNoDocuments {
		game, err = g.getRandomGame()
		if err != nil {
			game, err = g.createGame()
			if err != nil {
				return nil, err
			}

			return game, nil
		}

		if err = g.setGame(game.Id.Hex()); err != nil {
			return nil, err
		}
	}

	return game, err
}

func (g *GameInDB) InGame(word, user, id string) (*m.Game, error) {
	return g.inGame(word, user, id)
}

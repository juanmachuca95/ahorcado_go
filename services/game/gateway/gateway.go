package game

import (
	"context"
	"log"
	"time"

	database "github.com/juanmachuca95/ahorcado_go/internal/database/mongo"
	"github.com/juanmachuca95/ahorcado_go/services/game/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameGateway interface {
	CreateGame() (models.Game, error)
	GetGame() (models.Game, error)
}

type GameService struct {
	*mongo.Client
}

func NewGameGateway() GameGateway {
	return &GameService{
		database.Connect(),
	}
}

func (s *GameService) CreateGame() (models.Game, error) {
	collection := s.Client.Database("ahorcado").Collection("game")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	game := models.Game{
		ID:   primitive.NewObjectID(),
		Word: "Developer",
		Encontrados: []string{
			"a",
		},
		Winner: "Juan",
	}

	res, err := collection.InsertOne(ctx, game)
	if err != nil {
		log.Fatal("Error in CreateGame - error: ", err)
	}

	id := res.InsertedID
	log.Println(id)

	return game, nil
}

func (s *GameService) GetGame() (models.Game, error) {
	collection := s.Client.Database("ahorcado").Collection("game")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var games []models.Game
	for cur.Next(ctx) {
		var game models.Game
		err := cur.Decode(&game)
		if err != nil {
			log.Fatal(err)
		}

		games = append(games, game)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	log.Println(games)

	game := games[len(games)-1]
	return game, nil
}

package game

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/juanmachuca95/ahorcado_go/generated"
	database "github.com/juanmachuca95/ahorcado_go/internal/database/mongo"
	"github.com/juanmachuca95/ahorcado_go/services/game/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GameGateway interface {
	CreateGame() (models.Game, error)
	GetGame(*generated.Word) (generated.Game, error)
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
		Id:   primitive.NewObjectID(),
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

func (s *GameService) GetGame(word *generated.Word) (generated.Game, error) {
	collection := s.Client.Database("ahorcado").Collection("game")
	var game models.Game
	objID, _ := primitive.ObjectIDFromHex(word.GameId)
	err := collection.FindOne(context.TODO(), bson.D{{"_id", objID}}).Decode(&game)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return generated.Game{}, nil
		}
		panic(err)
	}

	if game.Finalizada {
		return generated.Game{
			Word:        game.Word,
			Encontrados: game.Encontrados,
			Winner:      game.Winner,
			Finalizada:  game.Finalizada,
		}, nil
	}

	if game.Word == word.Word {
		ok, err := s.UpdateGame(word, game.Id)
		if !ok {
			return generated.Game{}, err
		}

		game.Encontrados = append(game.Encontrados, game.Word)
		return generated.Game{
			Word:        game.Word,
			Winner:      word.User,
			Encontrados: game.Encontrados,
			Finalizada:  true,
		}, nil
	}

	if AlreadyFound(word.Word, game.Encontrados) {
		return generated.Game{}, nil
	}

	if strings.Contains(game.Word, word.Word) {
		game.Encontrados = append(game.Encontrados, game.Word)
		if Win(game.Word, game.Encontrados) {
			ok, err := s.UpdateGame(word, game.Id)
			if !ok {
				return generated.Game{}, err
			}

			game.Encontrados = append(game.Encontrados, game.Word)
			return generated.Game{
				Word:        game.Word,
				Winner:      word.User,
				Encontrados: game.Encontrados,
				Finalizada:  true,
			}, nil
		}
	}
	return generated.Game{}, nil
}

func (s GameService) UpdateGame(word *generated.Word, gameId primitive.ObjectID) (bool, error) {
	collection := s.Client.Database("ahorcado").Collection("game")
	id, _ := primitive.ObjectIDFromHex(gameId.Hex())

	game := models.Game{
		Word:       word.Word,
		Winner:     word.User,
		Finalizada: true,
	}

	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{{"$set", game}})
	if err != nil {
		return false, err
	}

	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
	return true, err
}

func Win(clave string, encontrados []string) bool {
	var fin bool = false
	var lengthClave int = len(clave)
	var lengthEncontrados int = 0
	for _, encontrado := range encontrados {
		count := strings.Count(clave, encontrado)
		lengthEncontrados += count
	}

	log.Println("Cantidad de encontrados = ", lengthEncontrados, " Cantidad total de la clave = ", lengthClave)
	if lengthEncontrados == lengthClave {
		fin = true
	}
	return fin
}

func AlreadyFound(character string, encontrados []string) bool {
	var result bool = false
	for _, encontrado := range encontrados {
		if character == encontrado {
			result = true
		}
	}
	return result
}

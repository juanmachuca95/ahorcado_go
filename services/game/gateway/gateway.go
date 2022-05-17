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
	GetRandomGame() (generated.Game, error)
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

func (s *GameService) GetRandomGame() (generated.Game, error) {
	//maxNum := 5 // The number of response attributes. Returns only one random car.
	var games []models.Game
	log.Println("Holanga")
	//db.mycoll.aggregate([{ $sample: { size: 1 } }])
	pipeline := []bson.D{bson.D{{"$sample", bson.D{{"size", 10}}}}}
	cursor, err := s.Client.Database("ahorcado").Collection("game").Aggregate(context.Background(), pipeline)
	if err != nil {
		return generated.Game{
			Id:          "",
			Word:        "",
			Winner:      "",
			Encontrados: []string{},
			Finalizada:  false,
			Error:       "No fue posible obtener un random game",
		}, err
	}

	for cursor.Next(context.Background()) {
		var game models.Game
		err := cursor.Decode(&game)
		log.Println(game)
		if err != nil {
			log.Fatal("Error decode - error: ", err.Error())
		} else {
			games = append(games, game)
		}
	}

	log.Println("GAMES ---> algo llego", games)
	return generated.Game{}, nil
}

func (s *GameService) GetGame(word *generated.Word) (generated.Game, error) {
	collection := s.Client.Database("ahorcado").Collection("game")
	var game models.Game

	if word.GameId == "" {
		word.GameId = "62811ff88d918fb8157de332"
	}

	objID, _ := primitive.ObjectIDFromHex(word.GameId)
	err := collection.FindOne(context.TODO(), bson.D{{"_id", objID}}).Decode(&game)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return generated.Game{
				Error: "Game no encontrado",
			}, nil
		}
		panic(err)
	}

	log.Println("JUGANDO CON LA PALABRA: ", game.Word)

	if game.Finalizada == true {
		return generated.Game{
			Word:        game.Word,
			Encontrados: game.Encontrados,
			Winner:      game.Winner,
			Finalizada:  game.Finalizada,
			Error:       "",
		}, nil
	}

	log.Println("Comparar ", game.Word, " con ", word.Word)
	if game.Word == word.Word {
		log.Println("COINCIDENCIA TOTAL!")
		ok, err := s.UpdateWinner(word, game.Id)
		if !ok {
			return generated.Game{
				Error: err.Error(),
			}, err
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
		return generated.Game{
			Error: fmt.Sprintf("La letra %v ya figura en la lista de encontrados", word.Word),
		}, nil
	}

	if strings.Contains(game.Word, word.Word) {
		log.Println("ESTA CONTENIDA")
		game.Encontrados = append(game.Encontrados, word.Word)

		log.Println("Encontrados hasta ahora: ", game.Encontrados)
		if Win(game.Word, game.Encontrados) { // si es la última letra para encontrar
			ok, err := s.UpdateWinner(word, game.Id)
			if !ok {
				return generated.Game{
					Error: fmt.Sprintf("No fue posible actualizar el Game - error: %v", err.Error()),
				}, err
			}

			game.Encontrados = append(game.Encontrados, game.Word)
			return generated.Game{
				Id:          game.Id.String(),
				Word:        game.Word,
				Winner:      word.User,
				Encontrados: game.Encontrados,
				Finalizada:  true,
			}, nil
		} else {
			_, err := s.UpdateEncontrados(game.Encontrados, game.Id)
			if err != nil {
				log.Fatal("No fue posible actualizar las letras encontradas - error: ", err)
			}

			return generated.Game{
				Id:          game.Id.String(),
				Word:        game.Word,
				Winner:      "",
				Encontrados: game.Encontrados,
				Error:       "",
			}, nil
		}
	}

	return generated.Game{
		Id:          game.Id.String(),
		Word:        game.Word,
		Winner:      "",
		Encontrados: game.Encontrados,
		Error:       fmt.Sprintf("La letra-palabra %s no es valida", word.Word),
	}, nil
}

func (s GameService) UpdateEncontrados(encontrados []string, gameId primitive.ObjectID) (bool, error) {
	collection := s.Client.Database("ahorcado").Collection("game")
	id, _ := primitive.ObjectIDFromHex(gameId.Hex())

	game := models.Game{
		Encontrados: encontrados,
	}

	result, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{{"$set", game}})
	if err != nil {
		return false, err
	}

	fmt.Printf("Updated %v Document - Encontrados\n", result.ModifiedCount)
	return true, err

}

func (s GameService) UpdateWinner(word *generated.Word, gameId primitive.ObjectID) (bool, error) {
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

	log.Println("Es el fin? ", fin)
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

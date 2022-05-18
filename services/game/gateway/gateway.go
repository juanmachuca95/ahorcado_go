package game

import (
	"context"
	"fmt"
	"log"
	"os"
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
	MyGame(*generated.Word) (generated.Game, error)
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
	pipeline := []bson.D{
		{{"$match", bson.D{{"finalizada", false}}}},
		{{"$sample", bson.D{{"size", 1}}}},
	}
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

	var game models.Game
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&game)
		if err != nil {
			log.Fatal("Error decode - error: ", err.Error())
		}
	}

	if game.Word != "" {
		return generated.Game{
			Id:          game.Id.Hex(),
			Word:        game.Word,
			Winner:      game.Winner,
			Encontrados: game.Encontrados,
			Finalizada:  game.Finalizada,
			Error:       "",
		}, nil
	}

	return generated.Game{
		Error: "No hay juegos dispnibles",
	}, nil
}

func (s *GameService) MyGame(word *generated.Word) (generated.Game, error) {
	game, _ := s.GetGame(word.GameId)

	/* Juego terminado */
	response := generated.Game{}

	if game.Finalizada == true {
		log.Println("1. Finalizada")
		errorMessage := fmt.Sprint("Este juego ha sido finalizado")
		response = generated.Game{
			Id:          game.Id.Hex(),
			Word:        game.Word,
			Winner:      game.Winner,
			Encontrados: game.Encontrados,
			Finalizada:  game.Finalizada,
			Error:       errorMessage,
		}
		return response, nil
	}

	/* El usuario ingresa una palabra y la palabra coincide */
	if game.Word == word.Word {
		log.Println("2. Ha encontrado la palabra.")
		game.Encontrados = append(game.Encontrados, game.Word)

		ok, err := s.UpdateWinner(word, game)
		if !ok {
			response = generated.Game{
				Id:          game.Id.Hex(),
				Word:        game.Word,
				Winner:      game.Winner,
				Encontrados: game.Encontrados,
				Finalizada:  game.Finalizada,
				Error:       err.Error(),
			}
			return response, err
		}

		response = generated.Game{
			Id:          game.Id.Hex(),
			Word:        game.Word,
			Winner:      word.User,
			Encontrados: game.Encontrados,
			Finalizada:  true,
			Error:       "",
		}
		return response, nil
	}

	/* letra ya encontrada */
	if AlreadyFound(word.Word, game.Encontrados) {
		log.Println("3. Ya ha sido encontrada la letra")
		messageError := fmt.Sprintf("La letra %v ya figura en la lista de encontrados", word.Word)
		return generated.Game{
			Id:          game.Id.Hex(),
			Word:        game.Word,
			Winner:      game.Winner,
			Encontrados: game.Encontrados,
			Finalizada:  game.Finalizada,
			Error:       messageError,
		}, nil
	}

	/* La letra ingresada por el usuario coincide con una letra en la palabra del juego */
	if strings.Contains(game.Word, word.Word) {
		log.Println("4. La palabra ingresada coincide con una letra de la palabra del juego")
		game.Encontrados = append(game.Encontrados, word.Word)
		if Win(game.Word, game.Encontrados) { // si es la última letra para encontrar
			ok, err := s.UpdateWinner(word, game)
			if !ok {
				messageError := fmt.Sprintf("No fue posible actualizar el Game - error: %v", err.Error())
				response = generated.Game{
					Id:          game.Id.Hex(),
					Word:        game.Word,
					Winner:      game.Winner,
					Encontrados: game.Encontrados,
					Finalizada:  game.Finalizada,
					Error:       messageError,
				}
				return response, err
			}

			response = generated.Game{
				Id:          game.Id.Hex(),
				Word:        game.Word,
				Winner:      word.User,
				Encontrados: game.Encontrados,
				Finalizada:  true,
				Error:       "",
			}
			return response, nil
		} else { // si no es la última letra del juego actualizamos los encontrados
			log.Println("5. Actualización de letra encontrada")
			_, err := s.UpdateEncontrados(game.Encontrados, game.Id.Hex())
			if err != nil {
				log.Fatal("No fue posible actualizar las letras encontradas - error: ", err)
			}

			response = generated.Game{
				Id:          game.Id.Hex(),
				Word:        game.Word,
				Winner:      game.Winner,
				Encontrados: game.Encontrados,
				Finalizada:  game.Finalizada,
				Error:       "",
			}
			return response, nil
		}
	}

	response = generated.Game{
		Id:          game.Id.Hex(),
		Word:        game.Word,
		Winner:      game.Winner,
		Encontrados: game.Encontrados,
		Error:       fmt.Sprintf("La letra o palabra ingresas %s no existe.", word.Word),
	}
	return response, nil
}

func (s GameService) UpdateEncontrados(encontrados []string, gameId string) (bool, error) {
	collection := s.Client.Database("ahorcado").Collection("game")
	objID, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		fmt.Println("ObjectIDFromHex ERROR", err)
	} else {
		fmt.Println("ObjectIDFromHex:", objID)
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}
	update := bson.M{"$set": bson.M{"encontrados": encontrados}}
	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	// Check for error, else print the UpdateOne() API call results
	if err != nil {
		fmt.Println("UpdateOne() result ERROR:", err)
		os.Exit(1)
	}

	log.Printf("Actualización de letras encontradas en el juego id - %s\n", gameId)
	return true, err
}

func (s GameService) UpdateWinner(word *generated.Word, game models.Game) (bool, error) {
	collection := s.Client.Database("ahorcado").Collection("game")
	objID, err := primitive.ObjectIDFromHex(game.Id.Hex())
	if err != nil {
		fmt.Println("ObjectIDFromHex ERROR", err)
	} else {
		fmt.Println("ObjectIDFromHex:", objID)
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}
	update := bson.M{"$set": bson.M{"encontrados": game.Encontrados, "finalizada": true, "winner": word.User}}
	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		update)
	if err != nil {
		return false, err
	}

	fmt.Printf("Se ha actualizado el gamador del juego - id %s\n", game.Id.Hex())
	return true, err
}

func (s GameService) GetGame(gameId string) (models.Game, error) {
	collection := s.Client.Database("ahorcado").Collection("game")

	objID, _ := primitive.ObjectIDFromHex(gameId)
	var game models.Game
	err := collection.FindOne(context.TODO(), bson.D{{"_id", objID}}).Decode(&game)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return game, err
		}
		panic(err)
	}

	return game, nil
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

package game

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	m "github.com/juanmachuca95/ahorcado_go/game/models"
	q "github.com/juanmachuca95/ahorcado_go/game/querys"
	help "github.com/juanmachuca95/ahorcado_go/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection = "game"

type GameStorage interface {
	getRandomGame() (*m.Game, error)
	getGame(string) (*m.Game, error)
	inGame(word, user, id string) (*m.Game, error)
}

type GameService struct {
	*mongo.Database
}

func NewGameStorageGateway(db *mongo.Database) GameStorage {
	return &GameService{db}
}

func (s *GameService) GetRandomGameToSet() (*m.Game, error) {
	pipeline := q.GetRandomGameToSet()
	cursor, err := s.Collection(collection).Aggregate(context.Background(), pipeline)
	if err != nil {
		return &m.Game{}, err
	}

	var game m.Game
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&game)
		if err != nil {
			log.Fatal("Error decode - error: ", err.Error())
		}
	}

	if game.Word != "" {
		return &game, nil
	}

	return &m.Game{}, errors.New("No hemos obtenido un getRandomGameToSet")
}

func (s *GameService) getRandomGame() (*m.Game, error) {
	pipeline := q.GetRandomGame()
	cursor, err := s.Collection(collection).Aggregate(context.Background(), pipeline)
	if err != nil {
		return &m.Game{
			Encontrados: []string{},
			Finalizada:  false,
		}, err
	}

	var game m.Game
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&game)
		if err != nil {
			log.Fatal("Error decode - error: ", err.Error())
		}
	}

	if game.Word != "" {
		return &game, nil
	}

	return &game, errors.New("No hay juegos disponibles.")
}

func (s *GameService) getGame(gameId string) (*m.Game, error) {
	collection := s.Collection(collection)
	objID, _ := primitive.ObjectIDFromHex(gameId)

	var game m.Game
	var query = q.GetGame(objID, false)
	err := collection.FindOne(context.TODO(), query).Decode(&game)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &game, err
		}
		panic(err)
	}

	return &game, nil
}

func (s *GameService) inGame(word, user, id string) (*m.Game, error) {
	log.Println("HOLA CHE --> word", word, " --> user ", user, " -- game id", id)
	game, err := s.getGame(id)
	if err != nil {
		return &m.Game{}, errors.New("El juego ha finalizado o no está disponible.")
	}

	if help.AlreadyFound(word, game.Encontrados) { // letra ya encontrada
		messageError := fmt.Sprintf("La letra %v ya figura en la lista de encontrados 👎", word)
		return game, errors.New(messageError)
	}

	if game.Word == word { // El usuario ingresa una palabra y la palabra coincide
		log.Println("2. Ha encontrado la palabra.")
		ok, err := s.UpdateWinner(word, user, *game)
		if !ok {
			return game, err
		}

		game.Finalizada = true
		return game, nil
	}

	if !strings.Contains(game.Word, word) { // La letra ingresada por el usuario coincide con una letra en la palabra del juego
		errorMessage := fmt.Sprintf("La letra o palabra ingresada (%s) no existe 👎, -1 intentos.", word)
		return game, errors.New(errorMessage)
	}

	log.Println("4. La palabra ingresada coincide con una letra de la palabra del juego")
	game.Encontrados = append(game.Encontrados, word)
	if help.Win(game.Word, game.Encontrados) { // si es la última letra para encontrar
		ok, err := s.UpdateWinner(word, user, *game)
		if !ok {
			messageError := fmt.Sprintf("No fue posible actualizar el Game - error: %v", err.Error())
			return game, errors.New(messageError)
		}

		game.Finalizada = true
		return game, nil
	}

	log.Println("5. Actualización de letra encontrada") // si no es la última letra del juego actualizamos los encontrados
	_, err = s.UpdateEncontrados(game.Encontrados, game.Id.Hex())
	if err != nil {
		log.Fatal("No fue posible actualizar las letras encontradas - error: ", err)
	}

	return game, nil
}

func (s *GameService) UpdateWinner(word, user string, game m.Game) (bool, error) {
	collection := s.Collection(collection)
	objID, err := primitive.ObjectIDFromHex(game.Id.Hex())
	if err != nil {
		fmt.Println("ObjectIDFromHex ERROR", err)
	} else {
		fmt.Println("ObjectIDFromHex:", objID)
	}

	filter, update := q.UpdateWinner(objID, game.Encontrados, user)
	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		return false, err
	}

	fmt.Printf("Se ha actualizado el ganador del juego - id %s\n", game.Id.Hex())
	_, err = s.UpdateGame()
	if err != nil {
		return false, err
	}

	return true, err
}

func (s *GameService) UpdateEncontrados(encontrados []string, gameId string) (bool, error) {
	collection := s.Collection(collection)
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

	if err != nil {
		return false, err
	}

	log.Printf("Actualización de letras encontradas en el juego id - %s\n", gameId)
	return true, nil
}

func (s *GameService) UpdateGame() (bool, error) {
	collection := s.Collection(collection)
	game, err := s.GetRandomGameToSet() // Obtengo un random game
	if err != nil {
		return false, errors.New("No se ha podido establecer un nuevo game")
	}

	objID, err := primitive.ObjectIDFromHex(game.Id.Hex())
	if err != nil {
		fmt.Println("ObjectIDFromHex ERROR", err)
		return false, err
	} else {
		fmt.Println("ObjectIDFromHex:", objID)
	}

	filter, update := q.UpdateGame(objID)
	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		return false, err
	}

	fmt.Printf("Se actualizo el proximo juego id: %s\n", game.Id)
	return true, err
}

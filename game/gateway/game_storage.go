package gateway

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	m "github.com/juanmachuca95/ahorcado_go/game/models"
	q "github.com/juanmachuca95/ahorcado_go/game/querys"
	helpers "github.com/juanmachuca95/ahorcado_go/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const _collectionGame = "game"
const (
	_codeFound        = 1
	_codeNotFound     = 2
	_codeAlreadyFound = 3
	_codeWinner       = 4
	_codeUnexpected   = 5
)

type GameStorage interface {
	getRandomGame() (*m.Game, error)
	getGame(string) (*m.Game, error)
	inGame(word, user, id string) (*m.Game, error)
	createGames() (bool, error)
}

type GameService struct {
	*mongo.Database
}

func NewGameStorageGateway(db *mongo.Database) GameStorage {
	return &GameService{db}
}

func (s *GameService) createGames() (bool, error) {
	coll := s.Collection(_collectionGame)
	docs := q.CreateGames()

	_, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		panic(err)
	}

	return true, nil
}

func (s *GameService) GetRandomGameToSet() (*m.Game, error) {
	pipeline := q.GetRandomGameToSet()
	cursor, err := s.Collection(_collectionGame).Aggregate(context.Background(), pipeline)
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

	return &m.Game{}, errors.New("no hemos obtenido un getRandomGameToSet")
}

func (s *GameService) getRandomGame() (*m.Game, error) {
	pipeline := q.GetRandomGame()
	cursor, err := s.Collection(_collectionGame).Aggregate(context.Background(), pipeline)
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

	return &game, errors.New("no hay juegos disponibles")
}

func (s *GameService) getGame(gameId string) (*m.Game, error) {
	collection := s.Collection(_collectionGame)
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
	game, err := s.getGame(id)
	if err != nil {
		return nil, errors.New("el juego ha finalizado o no está disponible")
	}

	if helpers.AlreadyFound(word, game.Encontrados) {
		game.Status = _codeAlreadyFound
		return game, nil
	}

	if game.Word == word {
		err := s.UpdateWinner(word, user, *game)
		if err != nil {
			game.Status = _codeUnexpected
			return game, err
		}

		game.Status = _codeWinner
		game.Finalizada = true
		return game, nil
	}

	if !strings.Contains(game.Word, word) {
		game.Status = _codeNotFound
		return game, nil
	}

	game.Encontrados = append(game.Encontrados, word)
	if helpers.Win(game.Word, game.Encontrados) {
		err := s.UpdateWinner(word, user, *game)
		if err != nil {
			game.Status = _codeUnexpected
			return game, nil
		}

		game.Status = _codeWinner
		game.Finalizada = true
		log.Println(game.Status)
		return game, nil
	}

	_, err = s.UpdateEncontrados(game.Encontrados, game.Id.Hex())
	if err != nil {
		log.Fatal("No fue posible actualizar las letras encontradas - error: ", err)
		return nil, err
	}

	game.Status = _codeFound
	return game, nil
}

func (s *GameService) UpdateWinner(word, user string, game m.Game) error {
	collection := s.Collection(_collectionGame)
	objID, err := primitive.ObjectIDFromHex(game.Id.Hex())
	if err != nil {
		fmt.Println("ObjectIDFromHex ERROR", err)
	}

	filter, update := q.UpdateWinner(objID, game.Encontrados, user)
	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		return err
	}

	/* fmt.Printf("Se ha actualizado el ganador del juego - id %s\n", game.Id.Hex())
	_, err = s.UpdateGame()
	if err != nil {
		return err
	}
	*/
	return nil
}

func (s *GameService) UpdateEncontrados(encontrados []string, gameId string) (bool, error) {
	collection := s.Collection(_collectionGame)
	objID, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		fmt.Println("ObjectIDFromHex ERROR", err)
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

	return true, nil
}

func (s *GameService) UpdateGame() (bool, error) {
	collection := s.Collection(_collectionGame)
	game, err := s.GetRandomGameToSet() // Obtengo un random game
	if err != nil {
		return false, errors.New("no se ha podido establecer un nuevo game")
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

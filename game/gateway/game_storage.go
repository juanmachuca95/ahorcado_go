package gateway

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/bxcodec/faker/v3"
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
	getGame() (*m.Game, error)
	setGame(id string) error

	getRandomGame() (*m.Game, error)
	getGameById(string) (*m.Game, error)
	inGame(word, user, id string) (*m.Game, error)
	createGame() (*m.Game, error)
	createGames() error
}

type GameService struct {
	*mongo.Database
}

func NewGameStorageGateway(db *mongo.Database) GameStorage {
	return &GameService{db}
}

func (s *GameService) createGames() error {
	coll := s.Collection(_collectionGame)
	docs := q.CreateGames()

	_, err := coll.InsertMany(context.TODO(), docs)
	if err != nil {
		return err
	}

	return nil
}

func (s *GameService) createGame() (*m.Game, error) {
	coll := s.Collection(_collectionGame)
	word := faker.Word()
	game := m.Game{
		Word:        word,
		Winner:      "",
		Encontrados: []string{},
		Finalizada:  false,
	}

	docs := q.CreateGame(word)
	_, err := coll.InsertOne(context.TODO(), docs)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (s *GameService) getGame() (*m.Game, error) {
	coll := s.Collection(_collectionGame)

	var game m.Game
	var query = q.GetGame()
	err := coll.FindOne(context.TODO(), query).Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (s *GameService) setGame(id string) error {
	coll := s.Collection(_collectionGame)
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("ObjectIDFromHex ERROR", err)
	}

	filter, update := q.SetGame(objID, true)
	_, err = coll.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		return err
	}

	return nil
}

// getRandomGame
// Obtiene un game con ingame en true y finalizado en false
func (s *GameService) getRandomGame() (*m.Game, error) {
	pipeline := q.GetRandomGame()
	cursor, err := s.Collection(_collectionGame).Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var game m.Game
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&game)
		if err != nil {
			log.Fatal("Error decode - error: ", err.Error())
		}
	}

	if game.Word == "" {
		return nil, errors.New("no hay games disponibles")
	}

	return &game, nil
}

func (s *GameService) getGameById(gameId string) (*m.Game, error) {
	coll := s.Collection(_collectionGame)
	objID, _ := primitive.ObjectIDFromHex(gameId)

	var game m.Game
	var query = q.GetGameById(objID, false)
	err := coll.FindOne(context.TODO(), query).Decode(&game)
	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (s *GameService) inGame(word, user, id string) (*m.Game, error) {
	game, err := s.getGame()
	if err != nil {
		return nil, err
	}

	if helpers.AlreadyFound(word, game.Encontrados) {
		game.Status = _codeAlreadyFound
		return game, nil
	}

	if strings.EqualFold(game.Word, word) {
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
	coll := s.Collection(_collectionGame)
	objID, err := primitive.ObjectIDFromHex(game.Id.Hex())
	if err != nil {
		fmt.Println("ObjectIDFromHex ERROR", err)
	}

	filter, update := q.UpdateWinner(objID, game.Encontrados, user)
	_, err = coll.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *GameService) UpdateEncontrados(encontrados []string, gameId string) (bool, error) {
	coll := s.Collection(_collectionGame)
	objID, err := primitive.ObjectIDFromHex(gameId)
	if err != nil {
		fmt.Println("ObjectIDFromHex ERROR", err)
	}

	filter := bson.M{"_id": bson.M{"$eq": objID}}
	update := bson.M{"$set": bson.M{"encontrados": encontrados}}
	_, err = coll.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Delete All games
// Only admin Role
func (s *GameService) DeleteAll() error {
	coll := s.Collection(_collectionGame)

	deleted, err := coll.DeleteMany(context.TODO(), nil)
	if err != nil {
		return err
	}

	log.Println("Se han eliminado ", deleted.DeletedCount, " games")
	return nil
}

package gateway

import (
	"context"
	"errors"
	"fmt"

	"github.com/juanmachuca95/ahorcado_go/game/models"
	"github.com/juanmachuca95/ahorcado_go/game/querys"
	"github.com/juanmachuca95/ahorcado_go/pkg/servicejwt"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	collection = "users"
	coll       *mongo.Collection
)

type AuthStorage interface {
	login(*models.Login) (string, error)
	register(l *models.Login) (string, error)
}

type AuthService struct {
	*mongo.Database
}

func NewAuthStorageGateway(db *mongo.Database) AuthStorage {
	return &AuthService{db}
}

func (s *AuthService) login(l *models.Login) (string, error) {
	coll = s.Collection(collection)
	filter := querys.Login(l.Username)

	var result models.Login
	err := coll.FindOne(context.TODO(), filter, nil).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("el username %s no existe", l.Username)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(l.Password)); err != nil {
		return "", errors.New("credenciales incorrectas")
	}

	jwt := servicejwt.NewServiceJWT()
	token, err := jwt.GenerateToken(*models.NewUser(l.Username), "player")
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) register(l *models.Login) (string, error) {
	coll = s.Collection(collection)
	filter := querys.Login(l.Username)

	var user models.Login
	err := coll.FindOne(context.TODO(), filter, nil).Decode(&user)
	if err == nil {
		return "", errors.New("el username ingresado ya existe")
	}

	passHashed, err := bcrypt.GenerateFromPassword([]byte(l.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("falló al hashear password")
	}

	l.Password = string(passHashed)
	doc := querys.Register(l)
	_, err = coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return "", err
	}

	jwt := servicejwt.NewServiceJWT()
	token, err := jwt.GenerateToken(*models.NewUser(l.Username), "player")
	if err != nil {
		return "", err
	}

	return token, nil
}

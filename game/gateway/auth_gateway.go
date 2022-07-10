package gateway

import (
	"github.com/juanmachuca95/ahorcado_go/game/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthGateway interface {
	Login(*models.Login) (string, error)
	Register(*models.Login) (string, error)
}

type AuthInDB struct {
	AuthStorage
}

func NewAuthGateway(db *mongo.Database) AuthGateway {
	return &AuthInDB{NewAuthStorageGateway(db)}
}

func (a *AuthInDB) Login(l *models.Login) (string, error) {
	return a.AuthStorage.login(l)
}

func (a *AuthInDB) Register(l *models.Login) (string, error) {
	return a.AuthStorage.register(l)
}

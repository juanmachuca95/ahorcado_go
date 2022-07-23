package querys

import (
	"github.com/juanmachuca95/ahorcado_go/game/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Register = func(l *models.Login) bson.D {
	return bson.D{primitive.E{Key: "username", Value: l.Username}, primitive.E{Key: "password", Value: l.Password}}
}

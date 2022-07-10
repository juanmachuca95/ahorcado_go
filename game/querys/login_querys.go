package querys

import (
	"github.com/juanmachuca95/ahorcado_go/game/models"
	"go.mongodb.org/mongo-driver/bson"
)

var Register = func(l *models.Login) interface{} {
	return bson.D{{"username", l.Username}, {"password", l.Password}}
}

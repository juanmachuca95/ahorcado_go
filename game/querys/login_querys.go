package querys

import (
	"github.com/juanmachuca95/ahorcado_go/game/models"
	"go.mongodb.org/mongo-driver/bson"
)

var Register = func(l *models.Login) bson.D {
	return bson.D{{"username", l.Username}, {"password", l.Password}}
}

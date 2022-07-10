package querys

import (
	"go.mongodb.org/mongo-driver/bson"
)

var Login = func(username string) interface{} {
	filter := bson.D{{"username", username}}
	return filter
}

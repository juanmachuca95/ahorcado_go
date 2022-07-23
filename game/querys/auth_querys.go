package querys

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Login = func(username string) bson.D {
	filter := bson.D{primitive.E{Key: "username", Value: username}}
	return filter
}

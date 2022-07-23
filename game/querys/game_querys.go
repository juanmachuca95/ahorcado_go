package querys

import (
	"strings"

	"github.com/bxcodec/faker/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CreateGames = func() []interface{} {
	var docs []interface{}

	for i := 0; i < 100; i++ {
		l := bson.D{
			primitive.E{Key: "word", Value: strings.ToUpper(faker.Word())},
			primitive.E{Key: "winner", Value: nil},
			primitive.E{Key: "finalizada", Value: false},
			primitive.E{Key: "ingame", Value: false},
			primitive.E{Key: "encontrados", Value: nil},
		}
		docs = append(docs, l)
	}

	return docs
}

var CreateGame = func(word string) primitive.D {
	game := bson.D{
		primitive.E{Key: "word", Value: word},
		primitive.E{Key: "winner", Value: nil},
		primitive.E{Key: "finalizada", Value: false},
		primitive.E{Key: "ingame", Value: true},
		primitive.E{Key: "encontrados", Value: nil}}
	return game
}

var GetGameById = func(id primitive.ObjectID, finalizada bool) primitive.M {
	return bson.M{
		"finalizada": finalizada,
		"_id":        id,
	}
}

var GetGame = func() primitive.M {
	return bson.M{
		"finalizada": false,
		"ingame":     true,
	}
}

var SetGame = func(objID primitive.ObjectID, status bool) (filter, update primitive.M) {
	filter = bson.M{"_id": bson.M{"$eq": objID}}
	update = bson.M{"$set": bson.M{"ingame": status}}
	return
}

var GetRandomGame = func() []primitive.D {
	return []bson.D{
		{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "finalizada", Value: false}}}},
		{primitive.E{Key: "$sample", Value: bson.D{primitive.E{Key: "size", Value: 1}}}},
	}
}

var UpdateWinner = func(objID primitive.ObjectID, encontrados []string, user string) (filter, update primitive.M) {
	filter = bson.M{"_id": bson.M{"$eq": objID}}
	update = bson.M{"$set": bson.M{"encontrados": encontrados, "finalizada": true, "winner": user, "ingame": false}}
	return
}

var UpdateGame = func(objID primitive.ObjectID) (filter, update primitive.M) {
	filter = bson.M{"_id": bson.M{"$eq": objID}}
	update = bson.M{"$set": bson.M{"ingame": true}}
	return
}

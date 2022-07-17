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
		l := bson.D{{"word", strings.ToUpper(faker.Word())}, {"winner", nil}, {"finalizada", false}, {"ingame", false}, {"encontrados", nil}}
		docs = append(docs, l)
	}

	return docs
}

var CreateGame = func(word string) primitive.D {
	game := bson.D{{"word", word}, {"winner", nil}, {"finalizada", false}, {"ingame", true}, {"encontrados", nil}}
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
		{{"$match", bson.D{{"finalizada", false}}}},
		{{"$sample", bson.D{{"size", 1}}}},
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

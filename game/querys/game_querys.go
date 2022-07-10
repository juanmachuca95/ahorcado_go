package querys

import (
	"github.com/bxcodec/faker/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CreateGames = func() []interface{} {
	var docs []interface{}

	for i := 0; i < 100; i++ {
		l := bson.D{{"word", faker.Word}, {"winner", ""}, {"finalizada", false}, {"ingame", false}, {"encontrados", nil}}
		docs = append(docs, l)
	}

	return docs
}

var GetGame = func(id primitive.ObjectID, finalizada bool) interface{} {
	return &bson.M{
		"finalizada": finalizada,
		"_id":        id,
	}
}

var GetRandomGame = func() interface{} {
	return []bson.D{
		{{"$match", bson.D{{"finalizada", false}, {"ingame", true}}}},
		{{"$sample", bson.D{{"size", 1}}}},
	}
}

var GetRandomGameToSet = func() interface{} {
	return []bson.D{
		{{"$match", bson.D{{"finalizada", false}}}},
		{{"$sample", bson.D{{"size", 1}}}},
	}
}

var UpdateWinner = func(objID primitive.ObjectID, encontrados []string, user string) (filter, update interface{}) {
	filter = bson.M{"_id": bson.M{"$eq": objID}}
	update = bson.M{"$set": bson.M{"encontrados": encontrados, "finalizada": true, "winner": user, "ingame": false}}
	return
}

var UpdateGame = func(objID primitive.ObjectID) (filter, update interface{}) {
	filter = bson.M{"_id": bson.M{"$eq": objID}}
	update = bson.M{"$set": bson.M{"ingame": true}}
	return
}

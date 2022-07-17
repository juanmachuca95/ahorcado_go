package querys

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var GetTop = func() (matchStage, limitStage, groupStage, sortStage primitive.D) {
	matchStage = bson.D{{"$match", bson.D{{"finalizada", true}}}}
	limitStage = bson.D{{"$limit", 10}}
	groupStage = bson.D{
		{"$group", bson.D{
			{"_id", "$winner"},
			{"winner", bson.D{{"$first", "$winner"}}},
			{"won", bson.D{
				{"$sum", 1},
			}},
		}}}
	sortStage = bson.D{{"$sort", bson.D{{"won", -1}}}}
	return
}

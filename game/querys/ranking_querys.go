package querys

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var GetTop = func() (matchStage, limitStage, groupStage, sortStage primitive.D) {
	matchStage = bson.D{primitive.E{Key: "$match", Value: bson.D{primitive.E{Key: "finalizada", Value: true}}}}
	limitStage = bson.D{primitive.E{Key: "$limit", Value: 10}}
	groupStage = bson.D{
		primitive.E{Key: "$group", Value: bson.D{
			primitive.E{Key: "_id", Value: "$winner"},
			primitive.E{Key: "winner", Value: bson.D{primitive.E{Key: "$first", Value: "$winner"}}},
			primitive.E{Key: "won", Value: bson.D{
				primitive.E{Key: "$sum", Value: 1},
			}},
		}}}
	sortStage = bson.D{primitive.E{Key: "$sort", Value: bson.D{primitive.E{Key: "won", Value: -1}}}}
	return
}

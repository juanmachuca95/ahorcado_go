package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ranking struct {
	Id     primitive.ObjectID `bson:"id" json:"id"`
	Winner string             `bson:"winner" json:"winner"`
	Won    int32              `bson:"won" json:"won"`
}

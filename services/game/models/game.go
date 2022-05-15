package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Game struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Word        string             `bson:"word,omitempty"`
	Encontrados []string           `bson:"encontrados,omitempty"`
	Winner      string             `bson:"winner,omitempty"`
}

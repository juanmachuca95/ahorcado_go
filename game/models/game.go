package game

import "go.mongodb.org/mongo-driver/bson/primitive"

type Game struct {
	Id          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Word        string             `bson:"word" json:"word,omitempty"`
	Encontrados []string           `bson:"encontrados" json:"econtrados,omitempty"`
	Winner      string             `bson:"winner" json:"winner,omitempty"`
	Finalizada  bool               `bson:"finalizada" json:"finalizada,omitempty"`
}

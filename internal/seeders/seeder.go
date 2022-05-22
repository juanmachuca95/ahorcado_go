package seeders

import (
	"context"
	"fmt"
	"log"

	"github.com/bxcodec/faker"
	clientMongoDB "github.com/juanmachuca95/ahorcado_go/internal/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Seeder struct {
	*clientMongoDB.MongoDBClient
}

func NewSeeder() *Seeder {
	return &Seeder{clientMongoDB.Connect()}
}

func (s *Seeder) SeedWords() {
	collection := s.Client.Database("ahorcado").Collection("game")

	var docs []interface{}
	for i := 0; i <= 500; i++ {
		var doc interface{}
		doc = bson.D{{"word", faker.WORD}, {"winner", ""}, {"finalzado", false}, {"encontrados", []string{}}}
		docs = append(docs, doc)
	}

	result, err := collection.InsertMany(context.TODO(), docs)
	if err != nil {
		log.Fatal(err.Error())
	}
	list_ids := result.InsertedIDs
	fmt.Printf("Documents inserted: %v\n", len(list_ids))
	for _, id := range list_ids {
		fmt.Printf("Inserted document with _id: %v\n", id)
	}
}

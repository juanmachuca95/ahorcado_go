package gateway

import (
	"context"
	"log"

	"github.com/juanmachuca95/ahorcado_go/game/models"
	"github.com/juanmachuca95/ahorcado_go/game/querys"
	"go.mongodb.org/mongo-driver/mongo"
)

type RankingStorage interface {
	getTop() ([]*models.Ranking, error)
}

type RankingService struct {
	*mongo.Database
}

func NewRankingStorage(db *mongo.Database) RankingStorage {
	return &RankingService{db}
}

func (r *RankingService) getTop() ([]*models.Ranking, error) {
	coll := r.Database.Collection(_collectionGame)

	matchStage, limitStage, groupStage, sortStage := querys.GetTop()
	cursor, err := coll.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage, sortStage, limitStage})
	if err != nil {
		log.Fatal(err)
	}

	//var rankings []models.Ranking
	var rankings []*models.Ranking
	if err = cursor.All(context.TODO(), &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}

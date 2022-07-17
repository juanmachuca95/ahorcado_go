package gateway

import (
	"github.com/juanmachuca95/ahorcado_go/game/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type RankingGateway interface {
	GetTop() ([]*models.Ranking, error)
}

type RankingInDB struct {
	RankingStorage
}

func NewRankingGateway(db *mongo.Database) *RankingInDB {
	return &RankingInDB{NewRankingStorage(db)}
}

func (r *RankingInDB) GetTop() ([]*models.Ranking, error) {
	return r.getTop()
}

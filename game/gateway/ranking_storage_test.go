package gateway

import (
	"testing"

	"github.com/joho/godotenv"
	database "github.com/juanmachuca95/ahorcado_go/pkg/database/mongo"
	"github.com/stretchr/testify/assert"
)

func TestGetTop(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}
	db := database.Connect()
	storage := NewRankingGateway(db)
	rankings, err := storage.getTop()

	t.Log(err)
	assert.Nil(t, err)
	assert.True(t, len(rankings) > 0)
}

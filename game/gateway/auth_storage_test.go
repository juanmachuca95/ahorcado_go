package gateway

import (
	"testing"

	"github.com/juanmachuca95/ahorcado_go/game/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthStorageLogin(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("success_login", func(mt *mtest.T) {
		coll = mt.Coll
		expectedUser := models.Login{
			Username: "john",
			Password: "Machuca12",
		}

		passHashed, err := bcrypt.GenerateFromPassword([]byte(expectedUser.Password), bcrypt.DefaultCost)
		require.NoError(t, err)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			primitive.E{Key: "username", Value: expectedUser.Username},
			primitive.E{Key: "password", Value: passHashed},
		}))

		authGtw := NewAuthStorageGateway(mt.DB)
		token, err := authGtw.login(&expectedUser)
		assert.Nil(t, err)
		assert.NotNil(t, token)
	})

}

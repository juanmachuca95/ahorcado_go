package servicejwt

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/juanmachuca95/ahorcado_go/game/models"
)

type AuthJWT interface {
	GenerateToken(models.User, string) (string, error) // Generar el token
	ValidateToken(string) (*claims, error)             // Validar el token
	// Decodificar el token
}

type ServiceJWT struct {
	secretKey string
}

func NewServiceJWT() AuthJWT {
	return &ServiceJWT{
		secretKey: getSecretKey(),
	}
}

type claims struct {
	User models.User `json:"user"`
	Role string      `json:"role"`
	jwt.RegisteredClaims
}

func getSecretKey() string {
	return os.Getenv("TOKEN_KEY")
}

func (s *ServiceJWT) GenerateToken(user models.User, role string) (string, error) {
	claims := &claims{
		User:             user,
		Role:             role,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Hour * 48)}},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	_token, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		log.Println("can not generate token")
		panic(err)
	}

	return _token, nil
}

func (s *ServiceJWT) ValidateToken(accessToken string) (*claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(s.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

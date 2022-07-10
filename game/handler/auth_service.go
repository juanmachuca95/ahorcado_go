package handler

import (
	"context"
	"errors"

	"github.com/juanmachuca95/ahorcado_go/game/gateway"
	"github.com/juanmachuca95/ahorcado_go/game/models"
	au "github.com/juanmachuca95/ahorcado_go/protos/auth"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	au.UnimplementedAuthServer
	authGtw gateway.AuthGateway
}

func NewAuthService(db *mongo.Database) *AuthService {
	return &AuthService{
		authGtw: gateway.NewAuthGateway(db),
	}
}

func (a *AuthService) Login(ctx context.Context, req *au.RequestLogin) (*au.ResponseLogin, error) {
	if req.Username == "" {
		return nil, errors.New("el campo username es requerido")
	}

	if req.Password == "" {
		return nil, errors.New("el campo password es requerido")
	}

	l := models.NewLogin(req.Username, req.Password)
	token, err := a.authGtw.Login(l)
	if err != nil {
		return nil, err
	}

	return &au.ResponseLogin{
		Token: token,
	}, nil
}

func (a *AuthService) Register(ctx context.Context, req *au.RequestLogin) (*au.ResponseLogin, error) {
	if req.Username == "" {
		return nil, errors.New("el campo username es requerido")
	}

	if req.Password == "" {
		return nil, errors.New("el campo password es requerido")
	}

	l := models.NewLogin(req.Username, req.Password)
	token, err := a.authGtw.Register(l)
	if err != nil {
		return nil, err
	}

	return &au.ResponseLogin{
		Token: token,
	}, nil
}

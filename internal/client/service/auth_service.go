package service

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
)

type authService struct {
	client   api.Invoker
	security interfaces.SecuritySource
}

func NewAuthService(client api.Invoker, security interfaces.SecuritySource) interfaces.AuthService {
	return &authService{
		client:   client,
		security: security,
	}
}

func (s *authService) Register(ctx context.Context, login, password string) (string, error) {
	res, err := s.client.RegisterPost(ctx, &api.UserCredentials{Login: login, Password: password})
	if err != nil {
		return "", err
	}
	switch res := res.(type) {
	case *api.AuthToken:
		return s.handleAuth(res, login, password)
	case *api.RegisterPostBadRequest:
		return "", interfaces.ErrBadRequest
	case *api.RegisterPostConflict:
		return "", interfaces.ErrLoginTaken
	default:
		return "", interfaces.ErrUnexpected
	}
}

func (s *authService) Login(ctx context.Context, login, password string) (string, error) {
	res, err := s.client.LoginPost(ctx, &api.UserCredentials{Login: login, Password: password})
	if err != nil {
		return "", err
	}
	switch res := res.(type) {
	case *api.AuthToken:
		return s.handleAuth(res, login, password)
	case *api.LoginPostBadRequest:
		return "", interfaces.ErrBadRequest
	case *api.Unauthorized:
		return "", interfaces.ErrUnauthorized
	default:
		return "", interfaces.ErrUnexpected
	}
}

func (s *authService) handleAuth(res *api.AuthToken, login, password string) (string, error) {
	s.security.SetToken(res.Token)
	return s.getUserID(res)
}

func (s *authService) getUserID(res *api.AuthToken) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(res.Token, &jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	return token.Claims.GetSubject()
}

package handlers

import (
	"context"
	"errors"

	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
)

type AuthHandler struct {
	service interfaces.Service
}

func NewAuthHandler(s interfaces.Service) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterPost(ctx context.Context, req *api.UserCredentials) (api.RegisterPostRes, error) {
	token, err := h.service.Register(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, interfaces.ErrLoginTaken) {
			return &api.RegisterPostConflict{}, nil
		}
		return nil, err
	}
	return &api.AuthToken{Token: token}, nil
}

func (h *AuthHandler) LoginPost(ctx context.Context, req *api.UserCredentials) (api.LoginPostRes, error) {
	token, err := h.service.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, interfaces.ErrUnauthorized) {
			return &api.Unauthorized{}, nil
		}
		return nil, err
	}
	return &api.AuthToken{Token: token}, nil
}

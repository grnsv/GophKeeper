package handlers

import (
	"context"

	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
	"github.com/ogen-go/ogen/ogenerrors"
)

type contextKey string

const userIDContextKey contextKey = "userID"

type SecurityHandler struct {
	jwts interfaces.JWTService
}

func NewSecurityHandler(jwts interfaces.JWTService) api.SecurityHandler {
	return &SecurityHandler{jwts: jwts}
}

func (h *SecurityHandler) HandleBearerAuth(ctx context.Context, operationName api.OperationName, t api.BearerAuth) (context.Context, error) {
	userID, err := h.jwts.ParseJWT(t.GetToken())
	if err != nil {
		return ctx, &ogenerrors.SecurityError{
			OperationContext: ogenerrors.OperationContext{Name: operationName},
			Security:         "BearerAuth",
			Err:              err,
		}
	}

	return context.WithValue(ctx, userIDContextKey, userID), nil
}

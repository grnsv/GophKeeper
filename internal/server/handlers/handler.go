package handlers

import (
	"context"
	"errors"

	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
)

var ErrUserIDNotFound = errors.New("user ID not found in context")

type Handler struct {
	*AuthHandler
	*RecordHandler
	*InfoHandler
}

func NewHandler(s interfaces.Service) api.Invoker {
	return &Handler{
		AuthHandler:   NewAuthHandler(s),
		RecordHandler: NewRecordHandler(s),
		InfoHandler:   NewInfoHandler(s),
	}
}

func getUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	if !ok {
		return "", ErrUserIDNotFound
	}
	return userID, nil
}

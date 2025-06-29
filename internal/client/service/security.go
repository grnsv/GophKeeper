package service

import (
	"context"
	"errors"
	"sync"

	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
)

type securitySource struct {
	mu    sync.RWMutex
	token string
}

func NewSecuritySource() interfaces.SecuritySource {
	return &securitySource{}
}

func (s *securitySource) BearerAuth(ctx context.Context, operationName string) (api.BearerAuth, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.token == "" {
		return api.BearerAuth{}, errors.New("token not set")
	}

	return api.BearerAuth{Token: s.token}, nil
}

func (s *securitySource) SetToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.token = token
}

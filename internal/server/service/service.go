package service

import (
	"context"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
	"github.com/grnsv/GophKeeper/internal/server/models"
)

type Service struct {
	storage      interfaces.Storage
	jwts         interfaces.JWTService
	buildVersion string
	buildDate    time.Time
}

func New(storage interfaces.Storage, jwts interfaces.JWTService, buildVersion, buildDate string) (interfaces.Service, error) {
	s := &Service{
		storage:      storage,
		jwts:         jwts,
		buildVersion: buildVersion,
	}
	if buildDate != "" {
		date, err := time.Parse("2006-01-02", buildDate)
		if err != nil {
			return nil, err
		}
		s.buildDate = date
	}

	return s, nil
}

func (s *Service) Register(ctx context.Context, login, password string) (string, error) {
	exists, err := s.storage.IsLoginExists(ctx, login)
	if err != nil {
		return "", err
	}
	if exists {
		return "", interfaces.ErrLoginTaken
	}

	user := &models.User{
		Login: login,
	}
	if user.PasswordHash, err = argon2id.CreateHash(password, argon2id.DefaultParams); err != nil {
		return "", err
	}
	if err = s.storage.CreateUser(ctx, user); err != nil {
		return "", err
	}

	return s.jwts.BuildJWT(user.ID)
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.storage.FindUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, interfaces.ErrNotFound) {
			return "", interfaces.ErrUnauthorized
		}
		return "", err
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil {
		return "", err
	}
	if !match {
		return "", interfaces.ErrUnauthorized
	}

	return s.jwts.BuildJWT(user.ID)
}

func (s *Service) GetRecords(ctx context.Context, userID string) ([]*models.Record, error) {
	return s.storage.GetRecords(ctx, userID)
}

func (s *Service) SaveRecord(ctx context.Context, rec *models.Record) error {
	if rec.Version > 1 {
		return s.storage.UpdateRecord(ctx, rec)
	}
	return s.storage.UpdateOrCreateRecord(ctx, rec)
}

func (s *Service) GetRecord(ctx context.Context, userID string, id uuid.UUID) (*models.Record, error) {
	return s.storage.GetRecord(ctx, userID, id)
}

func (s *Service) DeleteRecord(ctx context.Context, userID string, id uuid.UUID) error {
	return s.storage.DeleteRecord(ctx, userID, id)
}

func (s *Service) GetVersion(ctx context.Context) (buildVersion string, buildDate time.Time) {
	return s.buildVersion, s.buildDate
}

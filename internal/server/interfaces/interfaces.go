package interfaces

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/grnsv/GophKeeper/internal/server/models"
)

var (
	ErrLoginTaken      = errors.New("login already exists")
	ErrNotFound        = errors.New("not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrVersionConflict = errors.New("version conflict")
)

type Service interface {
	Register(ctx context.Context, login, password string) (token string, err error)
	Login(ctx context.Context, login, password string) (token string, err error)
	GetRecords(ctx context.Context, userID string) ([]*models.Record, error)
	SaveRecord(ctx context.Context, rec *models.Record) error
	GetRecord(ctx context.Context, userID string, id uuid.UUID) (*models.Record, error)
	DeleteRecord(ctx context.Context, userID string, id uuid.UUID) error
	GetVersion(ctx context.Context) (buildVersion string, buildDate time.Time)
}

type JWTService interface {
	BuildJWT(userID string) (token string, err error)
	ParseJWT(token string) (userID string, err error)
}

type Storage interface {
	UserRepository
	RecordRepository
}

type UserRepository interface {
	Close() error
	IsLoginExists(ctx context.Context, login string) (bool, error)
	CreateUser(ctx context.Context, user *models.User) error
	FindUserByLogin(ctx context.Context, login string) (*models.User, error)
}

type RecordRepository interface {
	Close() error
	GetRecords(ctx context.Context, userID string) ([]*models.Record, error)
	CreateRecord(ctx context.Context, rec *models.Record) error
	UpdateOrCreateRecord(ctx context.Context, rec *models.Record) error
	UpdateRecord(ctx context.Context, rec *models.Record) error
	GetRecord(ctx context.Context, userID string, id uuid.UUID) (*models.Record, error)
	DeleteRecord(ctx context.Context, userID string, id uuid.UUID) error
}

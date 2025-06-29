package interfaces

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

var (
	ErrLoginTaken      = errors.New("login already exists")
	ErrNotFound        = errors.New("not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrVersionConflict = errors.New("version conflict")
	ErrBadRequest      = errors.New("bad request")
	ErrUnexpected      = errors.New("unexpected response")
)

type Service interface {
	Close() error
	FetchServerVersion(ctx context.Context) (versionInfo models.VersionInfo, err error)
	Register(ctx context.Context, login, password string) error
	Login(ctx context.Context, login, password string) error
	GetRecords() ([]*models.Record, error)
	SaveRecord(ctx context.Context, record *models.Record) error
	GetRecord(id uuid.UUID) (*models.Record, error)
	FetchRecord(ctx context.Context, id uuid.UUID) (*models.Record, error)
	DeleteRecord(ctx context.Context, record *models.Record) error
	Sync(ctx context.Context) error
}

type Storage interface {
	Close() error
	GetRecords() ([]*models.Record, error)
	SaveRecord(record *models.Record) error
	GetRecord(id uuid.UUID) (*models.Record, error)
	IsRecordExists(id uuid.UUID) (exists bool, err error)
	DeleteRecord(id uuid.UUID) error
}

type NewUserStorage func(userID string, encryptionKey []byte) (Storage, error)

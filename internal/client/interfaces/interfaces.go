package interfaces

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/grnsv/GophKeeper/internal/api"
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

type SecuritySource interface {
	api.SecuritySource
	SetToken(token string)
}

type Service interface {
	AuthService
	CryptoService
	SyncService
	Storage
	FetchServerVersion(ctx context.Context) (versionInfo models.VersionInfo, err error)
}

type NewAuthService func(client api.Invoker, security SecuritySource) AuthService
type AuthService interface {
	Register(ctx context.Context, login, password string) (userID string, err error)
	Login(ctx context.Context, login, password string) (userID string, err error)
}

type NewCryptoService func(newCryptoStorage NewCryptoStorage) CryptoService
type CryptoService interface {
	InitCrypto(userID, login, password string) (Storage, error)
	EncryptRecord(record *models.Record) error
	DecryptRecord(record *models.Record) error
}

type NewSyncService func(client api.Invoker, storage Storage, crypto CryptoService) SyncService
type SyncService interface {
	PushRecord(ctx context.Context, record *models.Record) (*models.Record, error)
	PullRecord(ctx context.Context, id uuid.UUID) (*models.Record, error)
	ForgetRecord(ctx context.Context, record *models.Record) error
	Sync(ctx context.Context) error
}

type NewCryptoStorage func(userID string, encryptionKey []byte) (Storage, error)
type Storage interface {
	Close() error
	GetRecords() ([]*models.Record, error)
	SaveRecord(record *models.Record) error
	GetRecord(id uuid.UUID) (*models.Record, error)
	IsRecordExists(id uuid.UUID) (exists bool, err error)
	DeleteRecord(id uuid.UUID) error
}

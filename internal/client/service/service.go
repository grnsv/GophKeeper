package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

const maxUniqueIDAttempts = 5

type service struct {
	interfaces.AuthService
	interfaces.CryptoService
	interfaces.SyncService
	interfaces.Storage
	client         api.Invoker
	newSyncService interfaces.NewSyncService
}

func New(client api.Invoker, security interfaces.SecuritySource,
	newAuthService interfaces.NewAuthService,
	newCryptoService interfaces.NewCryptoService,
	newSyncService interfaces.NewSyncService,
	newCryptoStorage interfaces.NewCryptoStorage,
) interfaces.Service {
	return &service{
		AuthService:    newAuthService(client, security),
		CryptoService:  newCryptoService(newCryptoStorage),
		client:         client,
		newSyncService: newSyncService,
	}
}

func (s *service) FetchServerVersion(ctx context.Context) (versionInfo models.VersionInfo, err error) {
	res, err := s.client.VersionGet(ctx)
	if err != nil {
		return
	}
	if serverVersion, ok := res.BuildVersion.Get(); ok {
		versionInfo.BuildVersion.SetTo(serverVersion)
	}
	if serverBuildDate, ok := res.BuildDate.Get(); ok {
		versionInfo.BuildDate.SetTo(serverBuildDate.Format("2006-01-02"))
	}
	return
}

func (s *service) Register(ctx context.Context, login, password string) (string, error) {
	userID, err := s.AuthService.Register(ctx, login, password)
	return s.handleAuth(userID, login, password, err)
}

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	userID, err := s.AuthService.Login(ctx, login, password)
	return s.handleAuth(userID, login, password, err)
}

func (s *service) handleAuth(userID, login, password string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	s.Storage, err = s.CryptoService.InitCrypto(userID, login, password)
	if err != nil {
		return "", err
	}
	s.SyncService = s.newSyncService(s.client, s.Storage, s.CryptoService)

	return userID, nil
}

func (s *service) PushRecord(ctx context.Context, record *models.Record) (*models.Record, error) {
	if record.ID == uuid.Nil {
		id, err := s.newUniqueID()
		if err != nil {
			return record, err
		}
		record.ID = id
	}
	record.Version++
	record.Status = models.RecordStatusPending

	err := s.Storage.SaveRecord(record)
	if err != nil {
		return record, err
	}

	record, err = s.SyncService.PushRecord(ctx, record)
	if err != nil {
		return record, err
	}

	return record, s.Storage.SaveRecord(record)
}

func (s *service) newUniqueID() (uuid.UUID, error) {
	for range maxUniqueIDAttempts {
		id, err := uuid.NewRandom()
		if err != nil {
			return uuid.Nil, err
		}
		alreadyExists, err := s.Storage.IsRecordExists(id)
		if err != nil {
			return uuid.Nil, err
		}
		if !alreadyExists {
			return id, err
		}
	}
	return uuid.Nil, errors.New("failed to create unique id")
}

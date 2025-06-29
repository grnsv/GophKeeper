package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
	"github.com/grnsv/GophKeeper/internal/client/models"
	"golang.org/x/crypto/argon2"
)

type Service struct {
	client         *api.Client
	newUserStorage interfaces.NewUserStorage
	storage        interfaces.Storage
	encryptionKey  []byte
}

func New(client *api.Client, newUserStorage interfaces.NewUserStorage) interfaces.Service {
	return &Service{
		client:         client,
		newUserStorage: newUserStorage,
	}
}

func (s *Service) Close() error {
	if s.storage == nil {
		return nil
	}
	return s.storage.Close()
}

func (s *Service) FetchServerVersion(ctx context.Context) (versionInfo models.VersionInfo, err error) {
	res, err := s.client.VersionGet(ctx)
	if err != nil {
		return
	}
	if serverVersion, ok := res.BuildVersion.Get(); ok {
		versionInfo.BuildVersion.Set(serverVersion)
	}
	if serverBuildDate, ok := res.BuildDate.Get(); ok {
		versionInfo.BuildDate.Set(serverBuildDate.Format("2006-01-02"))
	}
	return
}

func (s *Service) Register(ctx context.Context, login, password string) error {
	res, err := s.client.RegisterPost(ctx, &api.UserCredentials{Login: login, Password: password})
	if err != nil {
		return err
	}
	switch v := res.(type) {
	case *api.AuthToken:
		return s.initUserStorage(v, login, password)
	case *api.RegisterPostBadRequest:
		return interfaces.ErrBadRequest
	case *api.RegisterPostConflict:
		return interfaces.ErrLoginTaken
	default:
		return interfaces.ErrUnexpected
	}
}

func (s *Service) Login(ctx context.Context, login, password string) error {
	res, err := s.client.LoginPost(ctx, &api.UserCredentials{Login: login, Password: password})
	if err != nil {
		return err
	}
	switch v := res.(type) {
	case *api.AuthToken:
		return s.initUserStorage(v, login, password)
	case *api.LoginPostBadRequest:
		return interfaces.ErrBadRequest
	case *api.Unauthorized:
		return interfaces.ErrUnauthorized
	default:
		return interfaces.ErrUnexpected
	}
}

func (s *Service) initUserStorage(res *api.AuthToken, login, password string) error {
	userID, err := s.getUserID(res)
	if err != nil {
		return fmt.Errorf("token error: %w", err)
	}

	s.encryptionKey = s.generateKey(userID, login, password)
	s.storage, err = s.newUserStorage(userID, s.encryptionKey)
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}

	return nil
}

func (s *Service) getUserID(res *api.AuthToken) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(res.Token, &jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	return token.Claims.GetSubject()
}

func (s *Service) generateKey(userID, login, password string) []byte {
	salt := []byte(login + userID)
	return argon2.IDKey([]byte(password), salt, 2, 128*1024, 4, 32)
}

func (s *Service) encryptRecord(record *models.Record) error {
	aesGCM, err := s.newGCM()
	if err != nil {
		return err
	}
	record.Nonce = make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, record.Nonce); err != nil {
		return err
	}
	record.Data = aesGCM.Seal(nil, record.Nonce, record.Data, nil)
	return nil
}

func (s *Service) decryptRecord(record *models.Record) error {
	aesGCM, err := s.newGCM()
	if err != nil {
		return err
	}
	record.Data, err = aesGCM.Open(nil, record.Nonce, record.Data, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) newGCM() (cipher.AEAD, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aesGCM, nil
}

func (s *Service) GetRecords() ([]*models.Record, error) {
	return s.storage.GetRecords()
}

func (s *Service) SaveRecord(ctx context.Context, record *models.Record) error {
	if record.ID == uuid.Nil {
		id, err := s.newUniqueID()
		if err != nil {
			return err
		}
		record.ID = id
	}
	record.Version++
	record.Status = models.RecordStatusPending

	err := s.storage.SaveRecord(record)
	if err != nil {
		return err
	}

	return s.syncPendingRecord(ctx, record)
}

func (s *Service) syncPendingRecord(ctx context.Context, record *models.Record) error {
	if err := s.encryptRecord(record); err != nil {
		return err
	}
	res, err := s.client.RecordsIDPut(ctx, &api.Record{
		ID:       api.NewOptUUID(record.ID),
		Type:     api.RecordType(record.Type),
		Data:     record.Data,
		Nonce:    record.Nonce,
		Metadata: record.Metadata,
		Version:  record.Version,
	}, api.RecordsIDPutParams{
		ID: record.ID,
	})
	if err != nil {
		return err
	}
	switch res.(type) {
	case *api.RecordsIDPutNoContent:
		record.Status = models.RecordStatusSynced
		return s.storage.SaveRecord(record)
	case *api.RecordsIDPutBadRequest:
		return interfaces.ErrBadRequest
	case *api.Unauthorized:
		return interfaces.ErrUnauthorized
	case *api.RecordsIDPutConflict:
		record.Status = models.RecordStatusConflict
		return s.storage.SaveRecord(record)
	default:
		return interfaces.ErrUnexpected
	}
}

func (s *Service) newUniqueID() (uuid.UUID, error) {
	for range 5 {
		id, err := uuid.NewRandom()
		if err != nil {
			return uuid.Nil, err
		}
		alreadyExists, err := s.storage.IsRecordExists(id)
		if err != nil {
			return uuid.Nil, err
		}
		if !alreadyExists {
			return id, err
		}
	}
	return uuid.Nil, errors.New("failed to create unique id")
}

func (s *Service) GetRecord(id uuid.UUID) (*models.Record, error) {
	return s.storage.GetRecord(id)
}

func (s *Service) FetchRecord(ctx context.Context, id uuid.UUID) (*models.Record, error) {
	res, err := s.client.RecordsIDGet(ctx, api.RecordsIDGetParams{ID: id})
	if err != nil {
		return nil, err
	}
	switch rec := res.(type) {
	case *api.RecordWithId:
		record := &models.Record{
			ID:       id,
			Type:     string(rec.Type),
			Data:     rec.Data,
			Nonce:    rec.Nonce,
			Metadata: rec.Metadata,
			Version:  rec.Version,
			Status:   models.RecordStatusSynced,
		}
		if err := s.decryptRecord(record); err != nil {
			return nil, err
		}
		return record, nil
	case *api.RecordsIDGetNotFound:
		return nil, interfaces.ErrNotFound
	case *api.Unauthorized:
		return nil, interfaces.ErrUnauthorized
	default:
		return nil, interfaces.ErrUnexpected
	}
}

func (s *Service) DeleteRecord(ctx context.Context, record *models.Record) error {
	if record.Status != models.RecordStatusDeleted {
		record.Status = models.RecordStatusDeleted
		err := s.storage.SaveRecord(record)
		if err != nil {
			return err
		}
	}
	res, err := s.client.RecordsIDDelete(ctx, api.RecordsIDDeleteParams{ID: record.ID})
	if err != nil {
		return err
	}
	switch res.(type) {
	case *api.RecordsIDDeleteNoContent:
		return s.storage.DeleteRecord(record.ID)
	case *api.Unauthorized:
		return interfaces.ErrUnauthorized
	default:
		return interfaces.ErrUnexpected
	}
}

func (s *Service) Sync(ctx context.Context) error {
	serverRecords, err := s.fetchRecords(ctx)
	if err != nil {
		return err
	}
	if err = s.syncIn(serverRecords); err != nil {
		return err
	}
	if err = s.syncOut(ctx, serverRecords); err != nil {
		return err
	}

	return nil
}

func (s *Service) fetchRecords(ctx context.Context) (map[uuid.UUID]*models.Record, error) {
	res, err := s.client.RecordsGet(ctx)
	if err != nil {
		return nil, err
	}
	switch res := res.(type) {
	case *api.RecordsGetOKApplicationJSON:
		return s.decryptFetchedRecords(res)
	case *api.Unauthorized:
		return nil, interfaces.ErrUnauthorized
	default:
		return nil, interfaces.ErrUnexpected
	}
}

func (s *Service) decryptFetchedRecords(res *api.RecordsGetOKApplicationJSON) (map[uuid.UUID]*models.Record, error) {
	records := make(map[uuid.UUID]*models.Record, len(*res)*8/7+1)
	for _, rec := range *res {
		record := &models.Record{
			ID:       rec.ID,
			Type:     string(rec.Type),
			Data:     rec.Data,
			Nonce:    rec.Nonce,
			Metadata: rec.Metadata,
			Version:  rec.Version,
		}
		if err := s.decryptRecord(record); err != nil {
			return nil, err
		}
		records[record.ID] = record
	}
	return records, nil
}

func (s *Service) syncIn(serverRecords map[uuid.UUID]*models.Record) error {
	for _, serverRecord := range serverRecords {
		localRecord, err := s.storage.GetRecord(serverRecord.ID)
		if err != nil {
			return err
		}
		if localRecord.Status != models.RecordStatusSynced {
			continue
		}
		if localRecord.Version-serverRecord.Version > 0 {
			localRecord.Status = models.RecordStatusConflict
			err = s.storage.SaveRecord(localRecord)
			if err != nil {
				return err
			}
		}
		err = s.storage.SaveRecord(serverRecord)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) syncOut(ctx context.Context, serverRecords map[uuid.UUID]*models.Record) error {
	localRecords, err := s.GetRecords()
	if err != nil {
		return err
	}
	for _, localRecord := range localRecords {
		switch localRecord.Status {
		case models.RecordStatusPending:
			s.syncPendingRecord(ctx, localRecord)
		case models.RecordStatusDeleted:
			s.DeleteRecord(ctx, localRecord)
		case models.RecordStatusConflict, models.RecordStatusSynced:
			_, exists := serverRecords[localRecord.ID]
			if !exists {
				s.DeleteRecord(ctx, localRecord)
			}
		}
	}

	return nil
}

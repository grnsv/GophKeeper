package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type syncService struct {
	client  api.Invoker
	storage interfaces.Storage
	crypto  interfaces.CryptoService
}

func NewSyncService(client api.Invoker, storage interfaces.Storage, crypto interfaces.CryptoService) interfaces.SyncService {
	return &syncService{
		client:  client,
		storage: storage,
		crypto:  crypto,
	}
}

func (s *syncService) PushRecord(ctx context.Context, record *models.Record) (*models.Record, error) {
	encrypted := *record
	if err := s.crypto.EncryptRecord(&encrypted); err != nil {
		return record, err
	}

	res, err := s.client.RecordsIDPut(ctx, &api.Record{
		ID:      api.NewOptUUID(encrypted.ID),
		Type:    api.RecordType(encrypted.Type),
		Data:    encrypted.Data,
		Nonce:   encrypted.Nonce,
		Version: encrypted.Version,
	}, api.RecordsIDPutParams{
		ID: encrypted.ID,
	})
	if err != nil {
		return record, err
	}

	switch res.(type) {
	case *api.RecordsIDPutNoContent:
		record.Status = models.RecordStatusSynced
	case *api.RecordsIDPutBadRequest:
		return record, interfaces.ErrBadRequest
	case *api.Unauthorized:
		return record, interfaces.ErrUnauthorized
	case *api.RecordsIDPutConflict:
		record.Status = models.RecordStatusConflict
	default:
		return record, interfaces.ErrUnexpected
	}
	return record, nil
}

func (s *syncService) PullRecord(ctx context.Context, id uuid.UUID) (*models.Record, error) {
	res, err := s.client.RecordsIDGet(ctx, api.RecordsIDGetParams{ID: id})
	if err != nil {
		return nil, err
	}
	switch rec := res.(type) {
	case *api.RecordWithId:
		record := &models.Record{
			ID:      id,
			Type:    models.RecordType(rec.Type),
			Data:    rec.Data,
			Nonce:   rec.Nonce,
			Version: rec.Version,
			Status:  models.RecordStatusSynced,
		}
		if err := s.crypto.DecryptRecord(record); err != nil {
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

func (s *syncService) ForgetRecord(ctx context.Context, record *models.Record) error {
	if record.Status != models.RecordStatusDeleted {
		record.Status = models.RecordStatusDeleted
		if err := s.storage.SaveRecord(record); err != nil {
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

func (s *syncService) Sync(ctx context.Context) error {
	serverRecords, err := s.fetchRecords(ctx)
	if err != nil {
		return err
	}
	if err = s.pull(serverRecords); err != nil {
		return err
	}
	if err = s.push(ctx, serverRecords); err != nil {
		return err
	}

	return nil
}

func (s *syncService) fetchRecords(ctx context.Context) (map[uuid.UUID]*models.Record, error) {
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

func (s *syncService) decryptFetchedRecords(res *api.RecordsGetOKApplicationJSON) (map[uuid.UUID]*models.Record, error) {
	records := make(map[uuid.UUID]*models.Record, len(*res)*8/7+1)
	for _, rec := range *res {
		record := &models.Record{
			ID:      rec.ID,
			Type:    models.RecordType(rec.Type),
			Data:    rec.Data,
			Nonce:   rec.Nonce,
			Version: rec.Version,
		}
		if err := s.crypto.DecryptRecord(record); err != nil {
			return nil, err
		}
		records[record.ID] = record
	}
	return records, nil
}

func (s *syncService) pull(serverRecords map[uuid.UUID]*models.Record) error {
	for _, serverRecord := range serverRecords {
		localRecord, err := s.storage.GetRecord(serverRecord.ID)
		if errors.Is(err, interfaces.ErrNotFound) {
			serverRecord.Status = models.RecordStatusSynced
			if err := s.storage.SaveRecord(serverRecord); err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if localRecord.Status != models.RecordStatusSynced {
			continue
		}
		if localRecord.Version > serverRecord.Version {
			localRecord.Status = models.RecordStatusConflict
			if err := s.storage.SaveRecord(localRecord); err != nil {
				return err
			}
		} else {
			serverRecord.Status = models.RecordStatusSynced
			if err := s.storage.SaveRecord(serverRecord); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *syncService) push(ctx context.Context, serverRecords map[uuid.UUID]*models.Record) error {
	localRecords, err := s.storage.GetRecords()
	if err != nil {
		return err
	}
	for _, localRecord := range localRecords {
		switch localRecord.Status {
		case models.RecordStatusPending:
			s.PushRecord(ctx, localRecord)
		case models.RecordStatusDeleted:
			s.ForgetRecord(ctx, localRecord)
		case models.RecordStatusConflict, models.RecordStatusSynced:
			_, exists := serverRecords[localRecord.ID]
			if !exists {
				s.ForgetRecord(ctx, localRecord)
			}
		default:
			s.PushRecord(ctx, localRecord)
		}
	}

	return nil
}

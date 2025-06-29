package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

type Storage struct {
	db *badger.DB
}

func New(userID string, encryptionKey []byte) (interfaces.Storage, error) {
	path, err := getDBPath(userID)
	if err != nil {
		return nil, err
	}
	db, err := badger.Open(badger.DefaultOptions(path).WithLogger(nil).WithEncryptionKey(encryptionKey))
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func getDBPath(userID string) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	appDir := filepath.Join(cacheDir, "GophKeeper")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(appDir, userID), nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) GetRecords() ([]*models.Record, error) {
	var records []*models.Record
	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var record models.Record
				if err := json.Unmarshal(v, &record); err != nil {
					return err
				}
				records = append(records, &record)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (s *Storage) SaveRecord(record *models.Record) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(record.ID[:], data)
	})
}

func (s *Storage) GetRecord(id uuid.UUID) (*models.Record, error) {
	var record models.Record
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(id[:])
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &record)
		})
	})
	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}
	return &record, err
}

func (s *Storage) IsRecordExists(id uuid.UUID) (exists bool, err error) {
	err = s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(id[:])
		return err
	})
	if err == nil {
		exists = true
	} else if errors.Is(err, badger.ErrKeyNotFound) {
		err = nil
	}

	return
}

func (s *Storage) DeleteRecord(id uuid.UUID) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(id[:])
	})
}

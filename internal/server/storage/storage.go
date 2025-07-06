package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
	_ "github.com/lib/pq"
)

type Storage struct {
	interfaces.UserRepository
	interfaces.RecordRepository
	db *sql.DB
}

func New(ctx context.Context, dsn, migrationsPath string) (interfaces.Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		return nil, err
	}
	defer m.Close()
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	storage := &Storage{db: db}
	storage.UserRepository, err = NewUserRepository(ctx, db)
	if err != nil {
		return nil, err
	}
	storage.RecordRepository, err = NewRecordRepository(ctx, db)
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) Close() error {
	var errs []error
	if err := s.UserRepository.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := s.RecordRepository.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := s.db.Close(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

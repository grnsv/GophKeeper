package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
	"github.com/grnsv/GophKeeper/internal/server/models"
)

type RecordRepository struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

func NewRecordRepository(ctx context.Context, db *sql.DB) (interfaces.RecordRepository, error) {
	r := &RecordRepository{
		db:    db,
		stmts: make(map[string]*sql.Stmt, 5),
	}
	if err := r.initStatements(ctx); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *RecordRepository) initStatements(ctx context.Context) error {
	queries := map[string]string{
		"GetRecords":   `SELECT * FROM records WHERE user_id = $1`,
		"CreateRecord": `INSERT INTO records (id, user_id, type, data, nonce, version) VALUES ($1, $2, $3, $4, $5, $6)`,
		"ExistsRecord": `SELECT EXISTS (SELECT 1 FROM records WHERE id = $1 AND user_id = $2) as exists`,
		"GetRecord":    `SELECT * FROM records WHERE id = $1 AND user_id = $2 LIMIT 1`,
		"DeleteRecord": `DELETE FROM records WHERE id = $1 AND user_id = $2`,
	}
	for key, query := range queries {
		stmt, err := r.db.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		r.stmts[key] = stmt
	}

	return nil
}

func (r *RecordRepository) Close() error {
	var errs []error
	for _, stmt := range r.stmts {
		if err := stmt.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (r *RecordRepository) GetRecords(ctx context.Context, userID string) ([]*models.Record, error) {
	var records []*models.Record
	rows, err := r.stmts["GetRecords"].QueryContext(ctx, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record models.Record
		if err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.Type,
			&record.Data,
			&record.Nonce,
			&record.Version,
		); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (r *RecordRepository) CreateRecord(ctx context.Context, rec *models.Record) error {
	if _, err := r.stmts["CreateRecord"].ExecContext(ctx, rec.ID, rec.UserID, rec.Type, rec.Data, rec.Nonce, rec.Version); err != nil {
		return err
	}
	return nil
}

func (r *RecordRepository) UpdateOrCreateRecord(ctx context.Context, rec *models.Record) error {
	var exists bool
	if err := r.stmts["ExistsRecord"].QueryRowContext(ctx, rec.ID, rec.UserID).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return r.UpdateRecord(ctx, rec)
	}
	return r.CreateRecord(ctx, rec)
}

func (r *RecordRepository) UpdateRecord(ctx context.Context, rec *models.Record) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentVersion int
	err = tx.QueryRowContext(ctx,
		"SELECT version FROM records WHERE id = $1 AND user_id = $2 FOR UPDATE",
		rec.ID, rec.UserID,
	).Scan(&currentVersion)
	if err != nil {
		return err
	}
	if rec.Version-currentVersion != 1 {
		return interfaces.ErrVersionConflict
	}
	_, err = tx.ExecContext(ctx,
		"UPDATE records SET type = $1, data = $2, nonce = $3, version = $4 WHERE id = $5 AND user_id = $6",
		rec.Type, rec.Data, rec.Nonce, rec.Version, rec.ID, rec.UserID,
	)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *RecordRepository) GetRecord(ctx context.Context, userID string, id uuid.UUID) (*models.Record, error) {
	var rec models.Record
	if err := r.stmts["GetRecord"].QueryRowContext(ctx, id, userID).Scan(
		&rec.ID,
		&rec.UserID,
		&rec.Type,
		&rec.Data,
		&rec.Nonce,
		&rec.Version,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrNotFound
		}
		return nil, err
	}
	return &rec, nil
}

func (r *RecordRepository) DeleteRecord(ctx context.Context, userID string, id uuid.UUID) error {
	if _, err := r.stmts["DeleteRecord"].ExecContext(ctx, id, userID); err != nil {
		return err
	}
	return nil
}

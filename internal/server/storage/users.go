package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/grnsv/GophKeeper/internal/server/interfaces"
	"github.com/grnsv/GophKeeper/internal/server/model"
)

type UserRepository struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

func NewUserRepository(ctx context.Context, db *sql.DB) (interfaces.UserRepository, error) {
	r := &UserRepository{
		db:    db,
		stmts: make(map[string]*sql.Stmt, 3),
	}
	if err := r.initStatements(ctx); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *UserRepository) initStatements(ctx context.Context) error {
	queries := map[string]string{
		"IsLoginExists":   `SELECT EXISTS(SELECT * FROM users WHERE login = $1) AS exists`,
		"CreateUser":      `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`,
		"FindUserByLogin": `SELECT * FROM users WHERE login = $1 LIMIT 1`,
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

func (r *UserRepository) Close() error {
	var errs []error
	for _, stmt := range r.stmts {
		if err := stmt.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (r *UserRepository) IsLoginExists(ctx context.Context, login string) (bool, error) {
	var exists bool
	if err := r.stmts["IsLoginExists"].QueryRowContext(ctx, login).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	if err := r.stmts["CreateUser"].QueryRowContext(ctx, user.Login, user.PasswordHash).Scan(&user.ID); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindUserByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	if err := r.stmts["FindUserByLogin"].QueryRowContext(ctx, login).Scan(
		&user.ID,
		&user.Login,
		&user.PasswordHash,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interfaces.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

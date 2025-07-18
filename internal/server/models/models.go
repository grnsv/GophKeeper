package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string
	Login        string
	PasswordHash string
	CreatedAt    time.Time
}

type Record struct {
	ID      uuid.UUID
	UserID  string
	Type    string
	Data    []byte
	Nonce   []byte
	Version int
}

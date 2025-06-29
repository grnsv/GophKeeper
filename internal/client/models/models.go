package models

import (
	"github.com/go-faster/jx"
	"github.com/google/uuid"
	"github.com/grnsv/GophKeeper/internal/api"
)

type OptString string

func (s OptString) String() string {
	if s == "" {
		return "N/A"
	}
	return string(s)
}

func (s *OptString) Set(v string) {
	*s = OptString(v)
}

type VersionInfo struct {
	BuildVersion OptString
	BuildDate    OptString
}

type Versions struct {
	Client VersionInfo
	Server VersionInfo
}

type RecordStatus string

const (
	RecordStatusPending  RecordStatus = "pending"
	RecordStatusSynced   RecordStatus = "synced"
	RecordStatusConflict RecordStatus = "conflict"
	RecordStatusDeleted  RecordStatus = "deleted"
)

type RecordType api.RecordType

const (
	RecordTypeCredentials RecordType = RecordType(api.RecordTypeCredentials)
	RecordTypeText        RecordType = RecordType(api.RecordTypeText)
	RecordTypeBinary      RecordType = RecordType(api.RecordTypeBinary)
	RecordTypeCard        RecordType = RecordType(api.RecordTypeCard)
)

type Record struct {
	ID       uuid.UUID
	Type     RecordType
	Data     []byte
	Nonce    []byte
	Metadata map[string]jx.Raw
	Version  int
	Status   RecordStatus
}

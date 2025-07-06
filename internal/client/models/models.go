package models

import (
	"github.com/google/uuid"
	"github.com/grnsv/GophKeeper/internal/api"
)

func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

type OptString struct {
	Value string
	Set   bool
}

func (o OptString) IsSet() bool { return o.Set }

func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

func (o OptString) String() string {
	if o.Set {
		return o.Value
	}
	return "N/A"
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
	ID      uuid.UUID
	Type    RecordType
	Data    []byte
	Nonce   []byte
	Version int
	Status  RecordStatus
}

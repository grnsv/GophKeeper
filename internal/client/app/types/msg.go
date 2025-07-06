package types

import "github.com/grnsv/GophKeeper/internal/client/models"

type ErrMsg struct {
	Err error
}

type ErrClearedMsg struct{}

type FetchVersionsMsg struct {
	ServerVersion models.VersionInfo
	Err           error
}

type MenuSelectedMsg struct {
	Item string
}

type RecordTypeSelectedMsg struct {
	RecordType models.RecordType
}

type BackToMenuMsg struct{}

type AuthMsg ErrMsg

type RecordsMsg struct {
	Records []*models.Record
	Err     error
}

type RecordSelectedMsg struct {
	Record *models.Record
}

type SyncTickMsg struct{}

type SyncMsg struct {
	HasConflicts bool
	Err          error
}

type DataMsg struct {
	Data []byte
}

type MetadataMsg struct {
	Metadata Metadata
}

type ConflictMsg struct {
	HasConflicts bool
}

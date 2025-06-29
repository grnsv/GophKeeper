package types

import "github.com/grnsv/GophKeeper/internal/client/models"

type FetchVersionsMsg struct {
	ServerVersion models.VersionInfo
	Offline       bool
	Err           error
}

type MenuSelectedMsg struct {
	Item string
}

type RecordTypeSelectedMsg struct {
	RecordType models.RecordType
}

type BackToMenuMsg struct{}

type AuthMsg struct {
	Err error
}

type RecordsMsg struct {
	Records []*models.Record
	Err     error
}

type SyncTickMsg struct{}

type SyncMsg struct {
	Err error
}

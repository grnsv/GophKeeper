package commands

import (
	"context"
	"encoding/json"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grnsv/GophKeeper/internal/client/app/types"
	"github.com/grnsv/GophKeeper/internal/client/interfaces"
	"github.com/grnsv/GophKeeper/internal/client/models"
)

const timeout = 5 * time.Second

func ClearErrorAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return types.ErrClearedMsg{}
	})
}

func FetchVersions(svc interfaces.Service) tea.Cmd {
	return func() tea.Msg {
		return fetchVersions(svc)
	}
}

func fetchVersions(svc interfaces.Service) types.FetchVersionsMsg {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	msg := types.FetchVersionsMsg{}
	msg.ServerVersion, msg.Err = svc.FetchServerVersion(ctx)

	return msg
}

func Select(item string) tea.Cmd {
	return func() tea.Msg {
		return types.MenuSelectedMsg{Item: item}
	}
}

func SelectType(recordType models.RecordType) tea.Cmd {
	return func() tea.Msg {
		return types.RecordTypeSelectedMsg{RecordType: recordType}
	}
}

func BackToMenu() tea.Msg {
	return types.BackToMenuMsg{}
}

func Register(svc interfaces.Service, login, password string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err := svc.Register(ctx, login, password)
		return types.AuthMsg{Err: err}
	}
}

func Login(svc interfaces.Service, login, password string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err := svc.Login(ctx, login, password)
		return types.AuthMsg{Err: err}
	}
}

func Show(svc interfaces.Service) tea.Cmd {
	return func() tea.Msg {
		records, err := svc.GetRecords()
		return types.RecordsMsg{Records: records, Err: err}
	}
}

func SyncTick() tea.Cmd {
	return tea.Tick(10*time.Minute, func(_ time.Time) tea.Msg {
		return types.SyncTickMsg{}
	})
}

func TrySync(svc interfaces.Service) tea.Cmd {
	return func() tea.Msg {
		msg := fetchVersions(svc)
		if msg.Err != nil {
			return msg
		}
		return Sync(svc)
	}
}

func Sync(svc interfaces.Service) tea.Cmd {
	return func() tea.Msg {
		return types.SyncMsg{Err: svc.Sync(context.Background())}
	}
}

func Error(err error) tea.Cmd {
	return func() tea.Msg {
		return types.ErrMsg{Err: err}
	}
}

func SubmitMetadata(metadata types.Metadata) tea.Cmd {
	return func() tea.Msg {
		return types.MetadataMsg{Metadata: metadata}
	}
}

func SubmitData[T types.Data](data T) tea.Cmd {
	return func() tea.Msg {
		data, err := json.Marshal(data)
		if err != nil {
			return Error(err)
		}
		return types.DataMsg{Data: data}
	}
}

func SaveRecord(svc interfaces.Service, record *models.Record) tea.Cmd {
	return func() tea.Msg {
		return types.ErrMsg{Err: svc.SaveRecord(context.Background(), record)}
	}
}

package handlers

import (
	"context"
	"errors"

	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
	"github.com/grnsv/GophKeeper/internal/server/models"
)

var ErrUserIDNotFound = errors.New("user ID not found in context")

type Handler struct {
	api.UnimplementedHandler
	service interfaces.Service
}

func NewHandler(s interfaces.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterPost(ctx context.Context, req *api.UserCredentials) (api.RegisterPostRes, error) {
	token, err := h.service.Register(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, interfaces.ErrLoginTaken) {
			return &api.RegisterPostConflict{}, nil
		}
		return nil, err
	}
	return &api.AuthToken{Token: token}, nil
}

func (h *Handler) LoginPost(ctx context.Context, req *api.UserCredentials) (api.LoginPostRes, error) {
	token, err := h.service.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, interfaces.ErrUnauthorized) {
			return &api.Unauthorized{}, nil
		}
		return nil, err
	}
	return &api.AuthToken{Token: token}, nil
}

func (h *Handler) VersionGet(ctx context.Context) (*api.VersionInfo, error) {
	buildVersion, buildDate := h.service.GetVersion(ctx)
	res := &api.VersionInfo{}
	if buildVersion != "" {
		res.BuildVersion = api.NewOptString(buildVersion)
	}
	if !buildDate.IsZero() {
		res.BuildDate = api.NewOptDate(buildDate)
	}

	return res, nil
}

func (h *Handler) getUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	if !ok {
		return "", ErrUserIDNotFound
	}
	return userID, nil
}

func (h *Handler) convertRecordToApiRecord(rec *models.Record) *api.RecordWithId {
	return &api.RecordWithId{
		ID:      rec.ID,
		Type:    api.RecordType(rec.Type),
		Data:    rec.Data,
		Nonce:   rec.Nonce,
		Version: rec.Version,
	}
}

func (h *Handler) RecordsGet(ctx context.Context) (api.RecordsGetRes, error) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	records, err := h.service.GetRecords(ctx, userID)
	if err != nil {
		return nil, err
	}

	length := len(records)
	out := make(api.RecordsGetOKApplicationJSON, length)
	for k, rec := range records {
		out[k] = *h.convertRecordToApiRecord(rec)
	}
	return &out, nil
}

func (h *Handler) RecordsIDPut(ctx context.Context, req *api.Record, params api.RecordsIDPutParams) (api.RecordsIDPutRes, error) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	rec := &models.Record{
		ID:      params.ID,
		UserID:  userID,
		Type:    string(req.Type),
		Data:    req.Data,
		Nonce:   req.Nonce,
		Version: req.Version,
	}
	if err = h.service.SaveRecord(ctx, rec); err != nil {
		if errors.Is(err, interfaces.ErrVersionConflict) {
			return &api.RecordsIDPutConflict{}, nil
		}
		return nil, err
	}
	return &api.RecordsIDPutNoContent{}, nil
}

func (h *Handler) RecordsIDGet(ctx context.Context, params api.RecordsIDGetParams) (api.RecordsIDGetRes, error) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	rec, err := h.service.GetRecord(ctx, userID, params.ID)
	if err != nil {
		if errors.Is(err, interfaces.ErrNotFound) {
			return &api.RecordsIDGetNotFound{}, nil
		}
		return nil, err
	}
	return h.convertRecordToApiRecord(rec), nil
}

func (h *Handler) RecordsIDDelete(ctx context.Context, params api.RecordsIDDeleteParams) (api.RecordsIDDeleteRes, error) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	if err = h.service.DeleteRecord(ctx, userID, params.ID); err != nil {
		return nil, err
	}
	return &api.RecordsIDDeleteNoContent{}, nil
}

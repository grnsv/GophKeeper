package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-faster/jx"
	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
	"github.com/grnsv/GophKeeper/internal/server/model"
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
	return &api.AuthToken{Token: api.NewOptString(token)}, nil
}

func (h *Handler) LoginPost(ctx context.Context, req *api.UserCredentials) (api.LoginPostRes, error) {
	token, err := h.service.Login(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, interfaces.ErrUnauthorized) {
			return &api.Unauthorized{}, nil
		}
		return nil, err
	}
	return &api.AuthToken{Token: api.NewOptString(token)}, nil
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

func (h *Handler) convertRecordToApiRecord(rec *model.Record) (*api.Record, error) {
	record := api.Record{
		ID:       api.NewOptUUID(rec.ID),
		Type:     api.RecordType(rec.Type),
		Data:     rec.Data,
		Metadata: make(api.RecordMetadata),
		Version:  rec.Version,
	}
	var metadata map[string]json.RawMessage
	if err := json.Unmarshal(rec.Metadata, &metadata); err != nil {
		return nil, fmt.Errorf("unmarshal metadata: %w", err)
	}
	for k, v := range metadata {
		record.Metadata[k] = jx.Raw(v)
	}

	return &record, nil
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
		record, err := h.convertRecordToApiRecord(rec)
		if err != nil {
			return nil, err
		}
		out[k] = *record
	}
	return &out, nil
}

func (h *Handler) RecordsIDPut(ctx context.Context, req *api.Record, params api.RecordsIDPutParams) (api.RecordsIDPutRes, error) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		return nil, err
	}
	rec := &model.Record{
		ID:      params.ID,
		UserID:  userID,
		Type:    string(req.Type),
		Data:    req.Data,
		Version: req.Version,
	}
	if req.Metadata != nil {
		metadataBytes, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, fmt.Errorf("marshal metadata: %w", err)
		}
		rec.Metadata = json.RawMessage(metadataBytes)
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
	return h.convertRecordToApiRecord(rec)
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

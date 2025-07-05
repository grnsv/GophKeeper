package handlers

import (
	"context"
	"errors"

	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
	"github.com/grnsv/GophKeeper/internal/server/models"
)

type RecordHandler struct {
	service interfaces.Service
}

func NewRecordHandler(s interfaces.Service) *RecordHandler {
	return &RecordHandler{service: s}
}

func (h *RecordHandler) convertRecordToApiRecord(rec *models.Record) *api.RecordWithId {
	return &api.RecordWithId{
		ID:      rec.ID,
		Type:    api.RecordType(rec.Type),
		Data:    rec.Data,
		Nonce:   rec.Nonce,
		Version: rec.Version,
	}
}

func (h *RecordHandler) RecordsGet(ctx context.Context) (api.RecordsGetRes, error) {
	userID, err := getUserID(ctx)
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

func (h *RecordHandler) RecordsIDPut(ctx context.Context, req *api.Record, params api.RecordsIDPutParams) (api.RecordsIDPutRes, error) {
	userID, err := getUserID(ctx)
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

func (h *RecordHandler) RecordsIDGet(ctx context.Context, params api.RecordsIDGetParams) (api.RecordsIDGetRes, error) {
	userID, err := getUserID(ctx)
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

func (h *RecordHandler) RecordsIDDelete(ctx context.Context, params api.RecordsIDDeleteParams) (api.RecordsIDDeleteRes, error) {
	userID, err := getUserID(ctx)
	if err != nil {
		return nil, err
	}
	if err = h.service.DeleteRecord(ctx, userID, params.ID); err != nil {
		return nil, err
	}
	return &api.RecordsIDDeleteNoContent{}, nil
}

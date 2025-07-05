package handlers

import (
	"context"

	"github.com/grnsv/GophKeeper/internal/api"
	"github.com/grnsv/GophKeeper/internal/server/interfaces"
)

type InfoHandler struct {
	service interfaces.Service
}

func NewInfoHandler(s interfaces.Service) *InfoHandler {
	return &InfoHandler{service: s}
}

func (h *InfoHandler) VersionGet(ctx context.Context) (*api.VersionInfo, error) {
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

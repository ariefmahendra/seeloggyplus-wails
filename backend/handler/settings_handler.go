package handler

import (
	"context"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/usecase"
)

type SettingsHandler struct {
	uc usecase.SettingUC
}

func NewSettingsHandler(uc usecase.SettingUC) *SettingsHandler {
	return &SettingsHandler{uc: uc}
}

func (h *SettingsHandler) Migrate(ctx context.Context) error {
	return h.uc.Migrate(ctx)
}

func (h *SettingsHandler) GetListSettings(ctx context.Context) ([]*dto.SettingsResponseDto, error) {
	return h.uc.GetListSettings(ctx)
}

func (h *SettingsHandler) FindByKey(ctx context.Context, key string) (*dto.SettingsResponseDto, error) {
	return h.uc.FindByKey(ctx, key)
}

func (h *SettingsHandler) UpdateSettings(ctx context.Context, settings *dto.SettingsRequestDto) (*dto.SettingsResponseDto, error) {
	return h.uc.UpdateSettings(ctx, settings)
}

func (h *SettingsHandler) InsertDefaultSettings(ctx context.Context) error {
	return h.uc.InsertDefaultSettings(ctx)
}

func (h *SettingsHandler) UpdateAllSettings(ctx context.Context, settings []*dto.SettingsRequestDto) error {
	return h.uc.UpdateAllSettings(ctx, settings)
}

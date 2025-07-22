package usecase

import (
	"context"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/entity"
	"seeloggyplus/backend/repository"
)

type SettingUC interface {
	Migrate(ctx context.Context) error
	GetListSettings(ctx context.Context) ([]*dto.SettingsResponseDto, error)
	FindByKey(ctx context.Context, key string) (*dto.SettingsResponseDto, error)
	UpdateSettings(ctx context.Context, settings *dto.SettingsRequestDto) (*dto.SettingsResponseDto, error)
	InsertDefaultSettings(ctx context.Context) error
	UpdateAllSettings(ctx context.Context, settings []*dto.SettingsRequestDto) error
}

type settingsUCImpl struct {
	repository repository.SettingsRepository
}

func (s *settingsUCImpl) Migrate(ctx context.Context) error {
	err := s.repository.Migrate(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *settingsUCImpl) GetListSettings(ctx context.Context) ([]*dto.SettingsResponseDto, error) {
	settingsList, err := s.repository.GetListSettings(ctx)
	if err != nil {
		return nil, err
	}

	var settingsListDto = make([]*dto.SettingsResponseDto, 0, len(settingsList))
	for _, setting := range settingsList {
		settingsListDto = append(settingsListDto, &dto.SettingsResponseDto{
			ID:          setting.ID,
			Name:        setting.Name,
			Key:         setting.Key,
			Value:       setting.Value,
			Description: setting.Description,
			Type:        setting.Type,
			CreatedAt:   setting.CreatedAt,
			UpdatedAt:   setting.UpdatedAt,
		})
	}

	return settingsListDto, nil
}

func (s *settingsUCImpl) FindByKey(ctx context.Context, key string) (*dto.SettingsResponseDto, error) {
	settings, err := s.repository.FindByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	return &dto.SettingsResponseDto{
		ID:          settings.ID,
		Name:        settings.Name,
		Key:         settings.Key,
		Value:       settings.Value,
		Description: settings.Description,
		Type:        settings.Type,
		CreatedAt:   settings.CreatedAt,
		UpdatedAt:   settings.UpdatedAt,
	}, nil
}

func (s *settingsUCImpl) UpdateSettings(ctx context.Context, settings *dto.SettingsRequestDto) (*dto.SettingsResponseDto, error) {
	var updatedSettingsPayload = &entity.Settings{
		Key:   settings.Key,
		Value: settings.Value,
	}

	updatedSettings, err := s.repository.UpdateSettings(ctx, updatedSettingsPayload)
	if err != nil {
		return nil, err
	}

	return &dto.SettingsResponseDto{
		ID:          updatedSettings.ID,
		Name:        updatedSettings.Name,
		Key:         updatedSettings.Key,
		Value:       updatedSettings.Value,
		Description: updatedSettings.Description,
		Type:        updatedSettings.Type,
		CreatedAt:   updatedSettings.CreatedAt,
		UpdatedAt:   updatedSettings.UpdatedAt,
	}, nil
}

func (s *settingsUCImpl) InsertDefaultSettings(ctx context.Context) error {
	defaultSettings := []*entity.Settings{
		{
			Key:         "light_theme",
			Name:        "Light Theme",
			Value:       "true",
			Description: "Default theme for the application",
			Type:        "boolean",
		},
		{
			Key:         "buffer_size",
			Name:        "Buffer Size",
			Value:       "100",
			Description: "Size of the buffer for log processing",
			Type:        "number",
		},
	}

	for _, setting := range defaultSettings {
		if err := s.repository.InsertDefaultSettings(ctx, setting); err != nil {
			return err
		}
	}
	return nil
}

func (s *settingsUCImpl) UpdateAllSettings(ctx context.Context, settings []*dto.SettingsRequestDto) error {
	var settingsEntities = make([]*entity.Settings, 0, len(settings))
	for _, setting := range settings {
		settingsEntities = append(settingsEntities, &entity.Settings{
			Key:   setting.Key,
			Value: setting.Value,
		})
	}

	return s.repository.UpdateAllSettings(ctx, settingsEntities)
}

func NewSettingsUC(repository repository.SettingsRepository) SettingUC {
	return &settingsUCImpl{repository: repository}
}

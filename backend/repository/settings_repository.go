package repository

import (
	"context"
	"database/sql"
	"errors"
	"seeloggyplus/backend/config"
	"seeloggyplus/backend/entity"
	"seeloggyplus/backend/logger"
)

type SettingsRepository interface {
	Migrate(ctx context.Context) error
	GetListSettings(ctx context.Context) ([]*entity.Settings, error)
	FindByKey(ctx context.Context, key string) (*entity.Settings, error)
	UpdateSettings(ctx context.Context, settings *entity.Settings) (*entity.Settings, error)
	InsertDefaultSettings(ctx context.Context, settings *entity.Settings) error
	UpdateAllSettings(ctx context.Context, settings []*entity.Settings) error
}

type settingsRepositoryImpl struct {
	db *sql.DB
}

func (s *settingsRepositoryImpl) Migrate(ctx context.Context) error {
	log := logger.Get()
	log.Info().Msg("Migrating settings table")

	_, err := s.db.ExecContext(ctx, config.MigrateSettingsTable)
	if err != nil {
		log.Error().Err(err).Msg("Failed to migrate settings table")
		return err
	}

	return nil
}

func (s *settingsRepositoryImpl) GetListSettings(ctx context.Context) ([]*entity.Settings, error) {
	log := logger.Get()
	log.Info().Msg("Retrieving list of settings")

	rows, err := s.db.QueryContext(ctx, config.GetListSettings)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve settings")
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close settings rows")
		}
	}(rows)

	var settingsList []*entity.Settings
	for rows.Next() {
		var settings entity.Settings
		if err := rows.Scan(
			&settings.ID,
			&settings.Name,
			&settings.Key,
			&settings.Value,
			&settings.Description,
			&settings.Type,
			&settings.CreatedAt,
			&settings.UpdatedAt,
		); err != nil {
			log.Error().Err(err).Msg("Failed to scan settings row")
			return nil, err
		}
		settingsList = append(settingsList, &settings)
	}

	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error occurred while iterating over settings rows")
		return nil, err
	}

	return settingsList, nil
}

func (s *settingsRepositoryImpl) FindByKey(ctx context.Context, key string) (*entity.Settings, error) {
	log := logger.Get()
	log.Info().Str("key", key).Msg("Finding settings by key")

	var settings entity.Settings
	err := s.db.QueryRowContext(ctx, config.FindSettingsByKey, key).Scan(
		&settings.ID,
		&settings.Name,
		&settings.Key,
		&settings.Value,
		&settings.Description,
		&settings.Type,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Str("key", key).Msg("No settings found for the given key")
			return nil, nil
		}
		log.Error().Err(err).Str("key", key).Msg("Failed to find settings by key")
		return nil, err
	}

	return &settings, nil
}

func (s *settingsRepositoryImpl) UpdateSettings(ctx context.Context, settings *entity.Settings) (*entity.Settings, error) {
	log := logger.Get()
	log.Info().Str("key", settings.Key).Msg("Updating settings")

	var updatedSettings entity.Settings
	err := s.db.QueryRowContext(
		ctx,
		config.UpdateSettings,
		settings.Value,
		settings.Key,
	).Scan(
		&updatedSettings.ID,
		&updatedSettings.Name,
		&updatedSettings.Key,
		&updatedSettings.Value,
		&updatedSettings.Description,
		&updatedSettings.Type,
		&updatedSettings.CreatedAt,
		&updatedSettings.UpdatedAt,
	)

	if err != nil {
		log.Error().Err(err).Str("key", settings.Key).Msg("Failed to update settings")
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Str("key", settings.Key).Msg("No settings found for update")
			return nil, nil
		}
		return nil, err
	}

	return &updatedSettings, nil
}

func (s *settingsRepositoryImpl) InsertDefaultSettings(ctx context.Context, settings *entity.Settings) error {
	log := logger.Get()
	log.Info().Str("key", settings.Key).Msg("Inserting default settings if not exists")

	_, err := s.db.Exec(
		config.InsertDefaultSettings,
		settings.Key,
		settings.Name,
		settings.Value,
		settings.Description,
		settings.Type,
	)

	if err != nil {
		log.Error().Err(err).Str("key", settings.Key).Msg("Failed to insert default settings")
		return err
	}

	return nil
}

func (s *settingsRepositoryImpl) UpdateAllSettings(ctx context.Context, settings []*entity.Settings) error {
	log := logger.Get()
	log.Info().Msg("Updating all settings")

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction for updating all settings")
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("Failed to rollback transaction")
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				log.Error().Err(commitErr).Msg("Failed to commit transaction")
			}
		}
	}()

	for _, setting := range settings {
		if _, err := tx.ExecContext(ctx, config.UpdateSettings, setting.Value, setting.Key); err != nil {
			log.Error().Err(err).Str("key", setting.Key).Msg("Failed to update setting in transaction")
			return err
		}
	}

	return nil
}

func NewSettingsRepository(db *sql.DB) SettingsRepository {
	return &settingsRepositoryImpl{
		db: db,
	}
}

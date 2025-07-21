package repository

import (
	"database/sql"
	"seeloggyplus/backend/entity"
)

type SettingsRepository interface {
	Migrate() error
	InsertSettings()
	GetListSettings() ([]*entity.Settings, error)
	FindByKey(key string) (*entity.Settings, error)
	UpdateSettings(settings *entity.Settings) error
}

type settingsRepositoryImpl struct {
	db *sql.DB
}

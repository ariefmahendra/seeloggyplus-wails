package config

const (
	MigrateSettingsTable = `CREATE TABLE IF NOT EXISTS mst_settings (
    		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    		name TEXT NOT NULL,
    		key TEXT NOT NULL UNIQUE,
    		value TEXT NOT NULL,
    		description TEXT,
    		type TEXT NOT NULL,
    		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`

	GetListSettings       = `SELECT id, name, key, value, description, type, created_at, updated_at FROM mst_settings`
	FindSettingsByKey     = `SELECT id, name, key, value, description, type, created_at, updated_at FROM mst_settings WHERE key = $1`
	UpdateSettings        = `UPDATE mst_settings SET value = $1, updated_at = CURRENT_TIMESTAMP WHERE key = $2 RETURNING id, name, key, value, description, type, created_at, updated_at`
	InsertDefaultSettings = `INSERT INTO mst_settings (key, name, value, description, type) VALUES ($1, $2, $3, $4, $5) ON CONFLICT(key) DO NOTHING`
)

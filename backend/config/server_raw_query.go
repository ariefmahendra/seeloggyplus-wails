package config

const (
	// CreateServer Server Query
	CreateServer   = `INSERT INTO mst_server(id, name, address, port, username, password, type) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id, name, address, port, username, password, type, created_at, updated_at`
	GetServerById  = `SELECT id, name, address, port, username, password, type, created_at, updated_at FROM mst_server WHERE id = $1`
	GetListServers = `SELECT id, name, address, port, username, password, type, created_at, updated_at FROM mst_server`
	UpdateServer   = `UPDATE mst_server SET name = $1, address = $2, port = $3, username = $4, password = $5, type = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $7`
	DeleteServer   = `DELETE FROM mst_server WHERE id = $1`

	ServerMigrateTable = `CREATE TABLE IF NOT EXISTS mst_server (
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		address TEXT NOT NULL,
		port TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		type TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`
)

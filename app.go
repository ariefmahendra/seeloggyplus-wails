package main

import (
	"context"
	"database/sql"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/handler"
	"seeloggyplus/backend/logger"
)

// App struct
type App struct {
	ctx    context.Context
	db     *sql.DB
	Server *handler.ServerHandler `json:"-"`
}

// NewApp creates a new App application struct
func NewApp(db *sql.DB, server *handler.ServerHandler) *App {
	return &App{
		db:     db,
		Server: server,
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	logger.InitLogger()

	appLog := logger.Get()
	if a.Server != nil {
		err := a.Server.Migrate(ctx)
		if err != nil {
			appLog.Error().Err(err).Msg("Error when migrate table")
		}
		appLog.Debug().Msg("Table migrated successfully")
	}
}

func (a *App) Shutdown(ctx context.Context) {
	logger.CloseLogfile()
	if a.db != nil {
		a.db.Close()
	}
}

func (a *App) ListServers() ([]*dto.ServerResponse, error) {
	return a.Server.ListServers(a.ctx)
}

func (a *App) AddServer(payload dto.ServerCreateRequest) (*dto.ServerResponse, error) {
	return a.Server.AddServer(a.ctx, &payload)
}

func (a *App) UpdateServer(payload dto.ServerUpdateRequest) error {
	return a.Server.UpdateServer(a.ctx, &payload)
}

func (a *App) DeleteServer(id string) error {
	return a.Server.DeleteServer(a.ctx, id)
}

func (a *App) GetServerById(id string) (*dto.ServerResponse, error) {
	return a.Server.GetServerById(a.ctx, id)
}

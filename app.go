package main

import (
	"context"
	"database/sql"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/handler"
	"seeloggyplus/backend/logger"
	"sync"
)

// App struct
type App struct {
	ctx        context.Context
	db         *sql.DB
	Server     *handler.ServerHandler     `json:"-"`
	Connection *handler.ConnectionHandler `json:"-"`

	cancellableRequests   map[string]context.CancelFunc
	cancellableRequestsMu sync.Mutex
}

// NewApp creates a new App application struct
func NewApp(db *sql.DB, server *handler.ServerHandler, connection *handler.ConnectionHandler) *App {
	return &App{
		db:                  db,
		Server:              server,
		Connection:          connection,
		cancellableRequests: make(map[string]context.CancelFunc),
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

func (a *App) TestConnection(requestID string, payload dto.ServerCreateRequest) error {
	ctx, cancel := context.WithCancel(a.ctx)

	a.cancellableRequestsMu.Lock()
	a.cancellableRequests[requestID] = cancel
	a.cancellableRequestsMu.Unlock()

	defer func() {
		a.cancellableRequestsMu.Lock()
		delete(a.cancellableRequests, requestID)
		a.cancellableRequestsMu.Unlock()
		log := logger.Get()
		log.Info().Str("requestID", requestID).Msg("Cleaned up cancellable request")
	}()

	return a.Connection.TestConnection(ctx, &payload)
}

func (a *App) CancelRequest(requestID string) {
	a.cancellableRequestsMu.Lock()
	defer a.cancellableRequestsMu.Unlock()

	log := logger.Get()
	if cancel, ok := a.cancellableRequests[requestID]; ok {
		log.Info().Str("requestID", requestID).Msg("Cancelling request")
		cancel()
		delete(a.cancellableRequests, requestID)
	} else {
		log.Warn().Str("requestID", requestID).Msg("Attempted to cancel a request that does not exist or is already completed")
	}
}

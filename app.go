package main

import (
	"context"
	"database/sql"
	"seeloggyplus/backend/dto"
	"seeloggyplus/backend/entity"
	"seeloggyplus/backend/handler"
	"seeloggyplus/backend/logger"
	"sync"
)

// App struct
type App struct {
	ctx            context.Context
	db             *sql.DB
	Server         *handler.ServerHandler         `json:"-"`
	Connection     *handler.ConnectionHandler     `json:"-"`
	RemoteFile     *handler.RemoteFileHandler     `json:"-"`
	LocalFile      *handler.LocalFileHandler      `json:"-"`
	SessionManager *handler.SessionManagerHandler `json:"-"`

	cancellableRequests   map[string]context.CancelFunc
	cancellableRequestsMu sync.Mutex
}

// NewApp creates a new App application struct
func NewApp(db *sql.DB, server *handler.ServerHandler, connection *handler.ConnectionHandler, remoteFile *handler.RemoteFileHandler, localFile *handler.LocalFileHandler, sessionManager *handler.SessionManagerHandler) *App {
	return &App{
		db:                  db,
		Server:              server,
		Connection:          connection,
		RemoteFile:          remoteFile,
		LocalFile:           localFile,
		SessionManager:      sessionManager,
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

// Server Management handler

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

// Connection handler

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

// Remote File Handler

func (a *App) GetListFiles(sessionID string, path string) ([]dto.FileInfo, error) {
	return a.RemoteFile.GetListFiles(sessionID, path)
}

// Local File Handler

func (a *App) ListFiles(path string) ([]dto.FileInfo, error) {
	return a.LocalFile.ListFiles(path)
}

func (a *App) GetUserHomeDir() (string, error) {
	return a.LocalFile.GetUserHomeDir()
}

// Session Manager Handler

func (a *App) ConnectSession(server *dto.ServerSessionManagement) (string, error) {
	return a.SessionManager.ConnectSession(a.ctx, server)
}

func (a *App) GetSession(sessionID string) (*entity.Session, bool) {
	return a.SessionManager.GetSession(sessionID)
}

func (a *App) CloseSession(sessionID string) error {
	return a.SessionManager.CloseSession(sessionID)
}

func (a *App) CloseAllSession() {
	a.SessionManager.CloseAllSession()
}

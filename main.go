package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"seeloggyplus/backend/handler"
	"seeloggyplus/backend/logger"
	"seeloggyplus/backend/repository"
	"seeloggyplus/backend/usecase"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	_ "modernc.org/sqlite"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	logZero := logger.Get()
	db, err := initializeDatabase()
	if err != nil {
		logZero.Fatal().Err(err).Msg("Error when initialize database")
	}

	serverRepo := repository.NewServerRepository(db)
	settingsRepo := repository.NewSettingsRepository(db)

	serverUC := usecase.NewServerUseCase(serverRepo)
	connectionUC := usecase.NewConnectionUseCase()
	sessionManagerUC := usecase.NewSessionManagerUC(serverRepo)
	remoteFileUC := usecase.NewRemoteFileUC(sessionManagerUC)
	localFileUC := usecase.NewLocalFileUC()
	settingsUC := usecase.NewSettingsUC(settingsRepo)

	serverHandler := handler.NewServerHandler(serverUC)
	connectionHandler := handler.NewConnectionHandler(connectionUC)
	remoteFileHandler := handler.NewRemoteFileHandler(remoteFileUC)
	sessionManagerHandler := handler.NewSessionManagerHandler(sessionManagerUC)
	localFileHandler := handler.NewLocalFileHandler(localFileUC)
	settingsHandler := handler.NewSettingsHandler(settingsUC)

	app := NewApp(
		db,
		serverHandler,
		connectionHandler,
		remoteFileHandler,
		localFileHandler,
		sessionManagerHandler,
		settingsHandler,
	)

	err = wails.Run(&options.App{
		Title:  "seeloggyplus",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		logZero.Fatal().Err(err).Msg("Error when run wails")
	}
}

func initializeDatabase() (*sql.DB, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan direktori config: %w", err)
	}
	appConfigDir := filepath.Join(configDir, "seeloggyplus")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori app: %w", err)
	}
	dbPath := filepath.Join(appConfigDir, "seeloggyplus.db")
	log.Println("Database path:", dbPath)

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka database: %w", err)
	}

	logZero := logger.Get()
	logZero.Debug().Msg("Database opened successfully")

	return db, nil
}

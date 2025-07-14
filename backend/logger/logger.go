package logger

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var appLogger zerolog.Logger
var logFile *os.File

func InitLogger() {
	var err error

	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting user config dir: %v", err)
	}
	appConfigDir := filepath.Join(configDir, "seeloggyplus/logs")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		log.Fatalf("Error creating config dir: %v", err)
	}
	logPath := filepath.Join(appConfigDir, "seeloggyplus.log")

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		PartsOrder: []string{
			zerolog.TimestampFieldName,
			zerolog.LevelFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
	}

	multiWriter := io.MultiWriter(consoleWriter, &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,   // max size in MB
		MaxBackups: 3,    // Max number of old log files to keep
		MaxAge:     28,   // Max age in days to keep a log file
		Compress:   true, // Compress old log files
	})

	appLogger = zerolog.New(multiWriter).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Logger()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	appLogger.Info().Msg("Logger initialized")
}

func Get() *zerolog.Logger {
	return &appLogger
}

func CloseLogfile() {
	appLogger.Info().Msg("Closing log file")
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			appLogger.Error().Err(err).Msg("Error closing log file")
		}
	}
}

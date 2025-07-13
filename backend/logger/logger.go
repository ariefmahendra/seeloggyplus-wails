package logger

import (
	"github.com/rs/zerolog"
	"io"
	"log"
	"os"
	"time"
)

var appLogger zerolog.Logger
var logFile *os.File

func InitLogger() {
	var err error

	file, err := os.OpenFile("seeloggyplus.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

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

	multiWriter := io.MultiWriter(consoleWriter, file)

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

package pkg

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

func NewPrettyLogger() *Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := &Logger{zerolog.New(output).With().Caller().Timestamp().Logger()}

	return logger
}

func NewLogger() *Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	logger := &Logger{zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()}

	return logger
}

func NewEmptyLogger() *Logger {
	logger := &Logger{zerolog.New(zerolog.Nop())}

	return logger
}

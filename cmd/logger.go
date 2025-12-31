package cmd

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log *LoggerService

type LoggerService struct {
	log zerolog.Logger
}

func NewLogger(file *os.File) {
	var output io.Writer

	if App.Environment == "development" {
		// Logging to both file and stdout during development
		fileOut := zerolog.ConsoleWriter{
			Out:        file,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		}
		consoleOut := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		output = zerolog.MultiLevelWriter(fileOut, consoleOut)

	} else if App.Environment == "production" {
		// Logging only to file during production
		output = zerolog.ConsoleWriter{
			Out:        file,
			TimeFormat: time.RFC3339,
			NoColor:    true,
		}

	} else {
		panic(errors.New("could not identify environment"))
	}

	Log = &LoggerService{
		log: zerolog.New(output).With().Timestamp().Logger(),
	}
}

func (l *LoggerService) Info(msg string) {
	l.log.WithLevel(zerolog.InfoLevel).Msgf("%s", msg)
}

func (l *LoggerService) Debug(msg string) {
	l.log.WithLevel(zerolog.DebugLevel).Msgf("%s", msg)
}

func (l *LoggerService) Warn(msg string) {
	l.log.WithLevel(zerolog.WarnLevel).Msgf("%s", msg)
}

func (l *LoggerService) Error(msg string, err error) {
	l.log.WithLevel(zerolog.ErrorLevel).Err(err).Msgf("%s", msg)
}

func (l *LoggerService) Fatal(msg string, err error) {
	l.log.WithLevel(zerolog.FatalLevel).Err(err).Msgf("%s", msg)
}

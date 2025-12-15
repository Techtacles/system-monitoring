package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

const tagName string = "service"

var logger zerolog.Logger

func init() {

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	logger = zerolog.New(output).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
}

func Info(tag, msg string) {
	logger := logger.With().Str(tagName, tag).Logger()
	logger.Info().Msg(msg)

}

func Error(tag, msg string, err error) {
	logger := logger.With().Str(tagName, tag).Logger()
	logger.Error().Err(err).Msg(msg)
}

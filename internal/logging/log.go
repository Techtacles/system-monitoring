package logging

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
)

const tagName string = "service"

var logger zerolog.Logger

func init() {
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				file = file[i+1:]
				break
			}
		}
		return file + ":" + strconv.Itoa(line)
	}

	logger = zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Caller().
		Logger()
}

func Info(tag, msg string) {
	logger := logger.With().Str(tagName, tag).Logger()
	logger.Info().Msg(msg)

}

func Error(tag, msg string, err error) {
	logger := logger.With().Str(tagName, tag).Logger()
	logger.Error().Msg(msg)
}

package log

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	log zerolog.Logger
}

func (l *Logger) SetOutput(out io.Writer) {
	l.log = log.Output(zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: time.RFC3339,
		NoColor:    !(out == os.Stdout || out == os.Stderr),
	})
}

func (l *Logger) Debug(msg string) {
	l.log.Debug().Msg(msg)
}

func (l *Logger) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *Logger) Warn(msg string) {
	l.log.Warn().Msg(msg)
}

func (l *Logger) Error(err error) {
	l.log.Error().Err(err).Msg("")
}

var Log = &Logger{}

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	Log.SetOutput(os.Stdout)
}

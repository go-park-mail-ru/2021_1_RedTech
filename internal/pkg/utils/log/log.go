package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(err error)
	SetOutput(out io.Writer)
}

type myLogger struct {
	log zerolog.Logger
}

func (l *myLogger) SetOutput(out io.Writer) {
	l.log = log.Output(zerolog.ConsoleWriter{Out: out})
}

func (l *myLogger) Debug(msg string) {
	l.log.Debug().Msg(msg)
}

func (l *myLogger) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *myLogger) Warn(msg string) {
	l.log.Warn().Msg(msg)
}

func (l *myLogger) Error(err error) {
	l.log.Error().Err(err).Msg("")
}

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	var Log = &myLogger{}
	Log.log = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

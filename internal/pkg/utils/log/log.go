package log

import (
	"io"
	"os"
	"time"

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
	l.log = log.Output(zerolog.ConsoleWriter{
		Out:        out,
		TimeFormat: time.RFC3339,
		NoColor:    out != os.Stdout || out != os.Stderr,
	})
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

var Log = &myLogger{}

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	Log.SetOutput(os.Stdout)
}

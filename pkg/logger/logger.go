package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	log zerolog.Logger
}

type Config struct {
	Level   string
	Format  string // "json" | "console"
	Service string
	Env     string
}

func New(cfg Config) *Logger {
	level, _ := zerolog.ParseLevel(cfg.Level)
	zerolog.SetGlobalLevel(level)

	var output zerolog.Logger

	if cfg.Format == "console" {
		output = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		output = zerolog.New(os.Stdout)
	}

	logger := output.
		With().
		Timestamp().
		Caller().
		Str("service", cfg.Service).
		Str("env", cfg.Env).
		Logger()

	return &Logger{log: logger}
}

func (l *Logger) Info() *zerolog.Event {
	return l.log.Info()
}

func (l *Logger) Error() *zerolog.Event {
	return l.log.Error()
}

func (l *Logger) Debug() *zerolog.Event {
	return l.log.Debug()
}

func (l *Logger) Msg(event *zerolog.Event, msg string) {
	event.Msg(msg)
}

func (l *Logger) Msgf(event *zerolog.Event, format string, args ...interface{}) {
	event.Msg(fmt.Sprintf(format, args...))
}

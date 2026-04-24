package logger

import (
	"context"

	"github.com/rs/zerolog"
)

type contextKey string

const (
	loggerKey        contextKey = "logger"
	correlationIDKey contextKey = "correlation_id"
)

func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationIDKey, id)
}

func WithLogger(ctx context.Context, log zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func FromContext(ctx context.Context) zerolog.Logger {
	if l, ok := ctx.Value(loggerKey).(zerolog.Logger); ok {
		return l
	}
	return zerolog.Nop()
}

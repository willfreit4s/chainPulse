package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Middleware(base *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			correlationID := r.Header.Get("X-Correlation-ID")
			if correlationID == "" {
				correlationID = uuid.New().String()
			}

			// datadog trace
			span, _ := tracer.SpanFromContext(r.Context())
			traceID := ""
			if span != nil {
				traceID = fmt.Sprintf("%d", span.Context().TraceID())
			}

			reqLogger := base.log.With().
				Str("correlation_id", correlationID).
				Str("trace_id", traceID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Dur("duration", time.Since(start)).
				Logger()

			reqLogger.Info().Msg("request started")

			ctx := r.Context()
			ctx = WithLogger(ctx, reqLogger)
			ctx = WithCorrelationID(ctx, correlationID)

			r = r.WithContext(ctx)

			rr := &responseRecorder{ResponseWriter: w, statusCode: 200}

			next.ServeHTTP(rr, r)

			reqLogger.Info().Msg("request completed")
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

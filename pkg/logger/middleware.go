package logger

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

			reqZerolog := base.With().
				Str("correlation_id", correlationID).
				Str("trace_id", traceID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Logger()

			reqLogger := &Logger{log: reqZerolog}

			reqLogger.Info().Msg("request started")

			ctx := r.Context()
			ctx = WithLogger(ctx, reqLogger)
			ctx = WithCorrelationID(ctx, correlationID)

			r = r.WithContext(ctx)

			rr := &responseRecorder{ResponseWriter: w, statusCode: 200}

			next.ServeHTTP(rr, r)

			reqLogger.Info().
				Int("status", rr.statusCode).
				Dur("duration", time.Since(start)).
				Msg("request completed")
		})
	}
}

// Create a file interceptor for gRPC
func LoggerInterceptor(base *Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		start := time.Now()

		var correlationID string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if values := md.Get("x-correlation-id"); len(values) > 0 {
				correlationID = values[0]
			}
		}

		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// datadog trace
		span, _ := tracer.SpanFromContext(ctx)
		traceID := ""
		if span != nil {
			traceID = fmt.Sprintf("%d", span.Context().TraceID())
		}

		reqZerolog := base.With().
			Str("correlation_id", correlationID).
			Str("trace_id", traceID).
			Str("grpc_method", info.FullMethod).
			Logger()

		reqLogger := &Logger{log: reqZerolog}

		ctx = WithLogger(ctx, reqLogger)
		ctx = WithCorrelationID(ctx, correlationID)

		reqLogger.Info().Msg("grpc request started")

		resp, err := handler(ctx, req)

		if err != nil {
			reqLogger.Error().
				Err(err).
				Dur("duration", time.Since(start)).
				Msg("grpc request failed")
		} else {
			reqLogger.Info().
				Dur("duration", time.Since(start)).
				Msg("grpc request completed")
		}

		return resp, err
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

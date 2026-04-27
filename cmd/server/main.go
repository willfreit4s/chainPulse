package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/willfreit4s/chainPulse/configs"
	"github.com/willfreit4s/chainPulse/internal/app"
	tokenv1 "github.com/willfreit4s/chainPulse/internal/infra/grpc/pb/proto/token/v1"
	"github.com/willfreit4s/chainPulse/pkg/database"
	"github.com/willfreit4s/chainPulse/pkg/logger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := initLogger(cfg)

	ctx := context.Background()
	ctx = logger.WithLogger(ctx, log)

	db, err := database.InitDatabase(ctx, cfg)
	if err != nil {
		log.Error().Err(err).Msg("database initialization failed")
		panic(err)
	}

	container := app.NewContainer(cfg, db)
	tokenSvc := container.TokenService
	grpcServer := initGRPCServer(tokenSvc, log)

	gateway, err := initGateway(ctx, tokenSvc)
	if err != nil {
		log.Error().Err(err).Msg("gateway initialization failed")
		panic(err)
	}

	httpServer := initHTTPServer(cfg, log, grpcServer, gateway)

	errCh := make(chan error, 1)
	go func() {
		log.Info().Int("port", cfg.ServerPort).Msg("server listening")
		errCh <- httpServer.ListenAndServe()
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-sigCtx.Done():
		log.Info().Msg("shutdown signal received")
	case serveErr := <-errCh:
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			log.Error().Err(serveErr).Msg("server stopped unexpectedly")
			panic(serveErr)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("http shutdown failed")
	} else {
		log.Info().Msg("http server stopped")
	}

	grpcServer.GracefulStop()
	log.Info().Msg("grpc server stopped")

	db.Close()
	log.Info().Msg("database pool closed")
	log.Info().Msg("shutdown completed")
}

func initGRPCServer(tokenSvc tokenv1.TokenServiceServer, log *logger.Logger) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.LoggerInterceptor(log)),
	)
	tokenv1.RegisterTokenServiceServer(grpcServer, tokenSvc)
	reflection.Register(grpcServer)
	return grpcServer
}

func SetupCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: false,
	})
}

func initLogger(cfg *configs.Config) *logger.Logger {
	log := logger.New(logger.Config{
		Level:   "info",
		Format:  "console",
		Service: cfg.ServiceName,
		Env:     cfg.Environment,
	})

	log.Info().
		Str("service_name", cfg.ServiceName).
		Str("environment", cfg.Environment).
		Int("port", cfg.ServerPort).
		Msg("starting application")

	return log
}

func initGateway(ctx context.Context, tokenSvc tokenv1.TokenServiceServer) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()
	if err := tokenv1.RegisterTokenServiceHandlerServer(ctx, mux, tokenSvc); err != nil {
		return nil, fmt.Errorf("register token gateway: %w", err)
	}
	return mux, nil
}

func initHTTPServer(
	cfg *configs.Config,
	log *logger.Logger,
	grpcServer *grpc.Server,
	handler http.Handler,
) *http.Server {

	handler = logger.Middleware(log)(handler)
	handler = applyCORS(handler)

	return &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.ServerPort),
		Handler:           multiplexHandler(grpcServer, handler),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func applyCORS(handler http.Handler) http.Handler {
	c := SetupCORS()
	return c.Handler(handler)
}

func multiplexHandler(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	mixedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
			return
		}
		httpHandler.ServeHTTP(w, r)
	})

	return h2c.NewHandler(mixedHandler, &http2.Server{})
}

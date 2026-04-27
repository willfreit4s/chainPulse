package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/willfreit4s/chainPulse/configs"
	tokenv1 "github.com/willfreit4s/chainPulse/internal/infra/grpc/pb/proto/token/v1"
	"github.com/willfreit4s/chainPulse/internal/infra/grpc/service"
	"github.com/willfreit4s/chainPulse/internal/usecase"
)

type Container struct {
	TokenService tokenv1.TokenServiceServer
}

func NewContainer(cfg *configs.Config, db *pgxpool.Pool) *Container {

	// repositories
	// tokenRepo := repository.NewTokenRepository(db)

	// usecases
	healthUC := usecase.NewHealthCheckUseCase(cfg)

	// services
	tokenService := service.NewTokenService(healthUC)

	return &Container{
		TokenService: tokenService,
	}
}

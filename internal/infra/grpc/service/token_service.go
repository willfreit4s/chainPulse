package service

import (
	"context"

	tokenv1 "github.com/willfreit4s/chainPulse/internal/infra/grpc/pb/proto/token/v1"
	"github.com/willfreit4s/chainPulse/internal/usecase"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TokenService struct {
	tokenv1.UnimplementedTokenServiceServer
	healthUC usecase.HealthCheckUseCase
}

func NewTokenService(healthUC usecase.HealthCheckUseCase) *TokenService {
	return &TokenService{
		healthUC: healthUC,
	}
}

func (s *TokenService) HealthCheck(ctx context.Context, _ *emptypb.Empty) (*tokenv1.HealthCheckResponse, error) {
	out, err := s.healthUC.Execute(ctx)
	if err != nil {
		return nil, err
	}

	return &tokenv1.HealthCheckResponse{
		Status:      out.Status,
		Version:     out.Version,
		ServiceName: out.ServiceName,
		Env:         out.Env,
	}, nil
}

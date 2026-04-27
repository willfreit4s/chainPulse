package usecase

import (
	"context"

	"github.com/willfreit4s/chainPulse/configs"
)

type HealthCheckOutput struct {
	Status      string
	Version     string
	ServiceName string
	Env         string
}

type HealthCheckUseCase interface {
	Execute(ctx context.Context) (*HealthCheckOutput, error)
}

type healthCheckUseCase struct {
	cfg *configs.Config
}

func NewHealthCheckUseCase(cfg *configs.Config) HealthCheckUseCase {
	return &healthCheckUseCase{cfg: cfg}
}

func (uc *healthCheckUseCase) Execute(ctx context.Context) (*HealthCheckOutput, error) {
	return &HealthCheckOutput{
		Status:      "ok",
		Version:     uc.cfg.ServiceVersion,
		ServiceName: uc.cfg.ServiceName,
		Env:         uc.cfg.Environment,
	}, nil
}

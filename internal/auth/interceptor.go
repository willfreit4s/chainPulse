package auth

import (
	"context"

	"github.com/willfreit4s/chainPulse/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor(v *Validator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log := logger.FromContext(ctx)

		if info.FullMethod == "/token.v1.TokenService/HealthCheck" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization")
		}

		claims, err := AuthenticateBearerToken(v, authHeaders[0])
		if err != nil {
			log.Error().Err(err).Msg("token validation failed")
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		ctx = WithClaims(ctx, claims)

		return handler(ctx, req)
	}
}

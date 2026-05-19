package auth

import (
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AuthenticateBearerToken(v *Validator, authHeader string) (*Claims, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
	}

	return v.Validate(token)
}

func Middleware(v *Validator, publicPaths ...string) func(http.Handler) http.Handler {
	public := make(map[string]struct{}, len(publicPaths))
	for _, path := range publicPaths {
		public[path] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			if _, ok := public[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := AuthenticateBearerToken(v, r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims)))
		})
	}
}

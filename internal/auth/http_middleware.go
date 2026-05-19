package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errInvalidAuthorizationHeader = errors.New("invalid authorization header")

func AuthenticateBearerToken(v *Validator, authHeader string) (*Claims, error) {
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, errInvalidAuthorizationHeader
	}

	token := parts[1]
	if token == "" {
		return nil, errInvalidAuthorizationHeader
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
				writeUnauthorizedResponse(w, r, next)
				return
			}

			next.ServeHTTP(w, r.WithContext(WithClaims(r.Context(), claims)))
		})
	}
}

func writeUnauthorizedResponse(w http.ResponseWriter, r *http.Request, next http.Handler) {
	unauthenticated := status.Error(codes.Unauthenticated, "unauthenticated")

	if mux, ok := next.(*runtime.ServeMux); ok {
		_, marshaler := runtime.MarshalerForRequest(mux, r)
		runtime.HTTPError(r.Context(), mux, marshaler, w, r, unauthenticated)
		return
	}

	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

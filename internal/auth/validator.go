package auth

import (
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Validator struct {
	jwks     *keyfunc.JWKS
	issuer   string
	audience string
}

func NewValidator(domain, audience string) (*Validator, error) {
	jwksURL := fmt.Sprintf(
		"https://%s/.well-known/jwks.json",
		domain,
	)

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return nil, err
	}

	return &Validator{
		jwks:     jwks,
		issuer:   fmt.Sprintf("https://%s/", domain),
		audience: audience,
	}, nil
}

func (v *Validator) Validate(tokenString string) (*Claims, error) {
	token, err := jwt.Parse(
		tokenString,
		v.jwks.Keyfunc,
		jwt.WithAudience(v.audience),
		jwt.WithIssuer(v.issuer),
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(5*time.Second),
	)

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if nbf, err := claimsMap.GetNotBefore(); err == nil {
		if time.Now().Before(nbf.Time) {
			return nil, fmt.Errorf("token not active yet")
		}
	}

	return &Claims{
		Subject: claimsMap["sub"].(string),
		Scope:   claimsMap["scope"].(string),
		Issuer:  claimsMap["iss"].(string),
	}, nil
}

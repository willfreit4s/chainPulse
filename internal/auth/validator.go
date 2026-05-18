package auth

import (
	"fmt"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Validator struct {
	jwks     *keyfunc.JWKS
	issuer   string
	audience string
}

func NewValidator(domain, audience string) (*Validator, error) {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)

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
	token, err := jwt.Parse(tokenString, v.jwks.Keyfunc)

	if err != nil {
		return nil, err
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claimsMap["iss"] != v.issuer {
		return nil, fmt.Errorf("invalid issuer")
	}

	return &Claims{
		Subject: claimsMap["sub"].(string),
		Scope:   claimsMap["scope"].(string),
	}, nil
}

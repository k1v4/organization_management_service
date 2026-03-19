package jwtpkg

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
)

type TokenVerifier struct {
	verifier *oidc.IDTokenVerifier
}

func NewTokenVerifier(ctx context.Context, issuerURL string) (*TokenVerifier, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, fmt.Errorf("jwtpkg: failed to get provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: "account"})
	return &TokenVerifier{verifier: verifier}, nil
}

func ExtractToken(c echo.Context) string {
	bearerToken := c.Request().Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}
	return strings.TrimPrefix(bearerToken, "Bearer ")
}

func (tv *TokenVerifier) ValidateToken(ctx context.Context, tokenString string) (*oidc.IDToken, error) {
	idToken, err := tv.verifier.Verify(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("jwtpkg: invalid token: %w", err)
	}
	return idToken, nil
}

func (tv *TokenVerifier) GetIdentityID(ctx context.Context, tokenString string) (string, error) {
	idToken, err := tv.ValidateToken(ctx, tokenString)
	if err != nil {
		return "", err
	}

	var claims struct {
		Sub string `json:"sub"`
	}
	if err = idToken.Claims(&claims); err != nil {
		return "", fmt.Errorf("jwtpkg: failed to parse claims: %w", err)
	}
	if claims.Sub == "" {
		return "", fmt.Errorf("jwtpkg: sub claim is empty")
	}

	return claims.Sub, nil
}

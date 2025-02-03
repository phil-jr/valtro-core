package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"resources/types"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("331f3866-08f2-4813-a34f-9e22c14e2d5f-e6f2f77c-a96a-4f82-8542-d45b589d7fad")

func ValidateJWT(tokenStr string) (*types.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &types.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is what we expect.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*types.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}

func extractBearerToken(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("authorization header does not start with %q", bearerPrefix)
	}

	// Remove the prefix and trim any extra spaces.
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, bearerPrefix))
	if token == "" {
		return "", fmt.Errorf("token is empty after removing prefix")
	}
	return token, nil
}

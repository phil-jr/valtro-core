package handlers

import (
	"errors"
	"time"
	"users/types"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("331f3866-08f2-4813-a34f-9e22c14e2d5f-e6f2f77c-a96a-4f82-8542-d45b589d7fad")

type Claims struct {
	Email       string `json:"email"`
	UserUuid    string `json:"userUuid"`
	CompanyUuid string `json:"companyUuid"`
	jwt.RegisteredClaims
}

func GenerateJWT(user types.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Email:       user.Email,
		UserUuid:    user.UserId,
		CompanyUuid: user.CompanyId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is what we expect.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Optionally validate standard claims (e.g. Expiration)
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}

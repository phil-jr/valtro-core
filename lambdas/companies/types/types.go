package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email       string `json:"email"`
	UserUuid    string `json:"userUuid"`
	CompanyUuid string `json:"companyUuid"`
	jwt.RegisteredClaims
}

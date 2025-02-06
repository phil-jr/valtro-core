package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	UserId    string    `json:"userId"`
	CompanyId string    `json:"companyId"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Admin     bool      `json:"admin"`
	CreatedAt time.Time `json:"createdAt"`
}

type Metric struct {
	Name      string  `json:"metricName"`
	Value     float64 `json:"metricValue"`
	Unit      string  `json:"metricUnit"`
	Timestamp string  `json:"timestamp"`
}

type Claims struct {
	Email       string `json:"email"`
	UserUuid    string `json:"userUuid"`
	CompanyUuid string `json:"companyUuid"`
	jwt.RegisteredClaims
}

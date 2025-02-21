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
	Aggregate int8    `json:"aggregate"`
}

type Claims struct {
	Email       string `json:"email"`
	UserUuid    string `json:"userUuid"`
	CompanyUuid string `json:"companyUuid"`
	jwt.RegisteredClaims
}

type Cost struct {
	ResourceID     string    `json:"resourceId"`
	Cost           float64   `json:"cost"`
	Aggregate      int       `json:"aggregate"`
	StartTimestamp time.Time `json:"startTimestamp"`
	EndTimestamp   time.Time `json:"endTimestamp"`
	CreatedAt      time.Time `json:"createdAt"`
}

type Resource struct {
	ResourceID   string    `json:"resourceId"`
	ResourceName string    `json:"resourceName"`
	ResourceType string    `json:"resourceType"`
	CreatedAt    time.Time `json:"createdAt"`
}

type EvenMetric struct {
	Name  string    `json:"metricName"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Value float64   `json:"value"`
	Unit  string    `json:"metricUnit"`
}

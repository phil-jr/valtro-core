package types

import "time"

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

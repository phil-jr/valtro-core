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

type SignIn struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type JwtSignInResponse struct {
	Token string `json:"jwtToken"`
	Ttl   int32  `json:"ttl"`
}

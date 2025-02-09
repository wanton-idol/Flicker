package dto

import "github.com/golang-jwt/jwt"

type UserRegistrationDTO struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Code      string `json:"code"`
	Mobile    string `json:"mobile"`
	Password  string `json:"password"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserData struct {
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type VerifyUser struct {
	User *UserData `json:"user,omitempty"`
	OTP  string    `json:"otp,omitempty"`
}

type GoogleLoginResponse struct {
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Sub           string `json:"sub"`
}

type UserEmail struct {
	Email string `json:"email"`
}

type UserIDToken struct {
	IDToken string `json:"id_token"`
}

type TokenInfo struct {
	Iss           string `json:"iss"`
	Sub           string `json:"sub"` // userId
	Azp           string `json:"azp"`
	Aud           string `json:"aud"` // clientId
	Iat           int64  `json:"iat"`
	Exp           int64  `json:"exp"` // expired time
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Local         string `json:"locale"`
	jwt.StandardClaims
}

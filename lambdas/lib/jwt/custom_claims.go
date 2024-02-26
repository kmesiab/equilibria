package jwt

import "github.com/golang-jwt/jwt"

type CustomClaims struct {
	*jwt.StandardClaims
	UserID          int64
	Email           string
	PhoneNumber     string
	Firstname       string
	Lastname        string
	AccountStatus   string
	AccountStatusID int64
	PhoneVerified   bool
}

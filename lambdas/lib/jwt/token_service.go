package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang-jwt/jwt"

	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/models"
)

const EnableAWSSessionDebug = false

const (
	Issuer                       = "equilibria"
	Audience                     = "equilibria"
	ParameterStorePrivateKeyName = "private_rsa_key"
	ParameterStorePublicKeyName  = "public_rsa_key"
)

// TokenServiceInterface defines the methods for handling JWT token operations.
type TokenServiceInterface interface {
	Generate(user *models.User, expirationMinutes int, key *rsa.PrivateKey) (string, error)
	Validate(tokenString string, publicKey *rsa.PublicKey) (*CustomClaims, error)
	Issue(user *models.User, expiration int, keyRotator KeyRotatorInterface) (string, error)
	GetPrivateKey(rotator KeyRotatorInterface) (*rsa.PrivateKey, error)
	GetPublicKey(rotator KeyRotatorInterface) (*rsa.PublicKey, error)
}

type TokenService struct{}

func (t TokenService) Generate(user *models.User, expirationMinutes int, key *rsa.PrivateKey) (string, error) {

	claims := CreateCustomClaims(user, expirationMinutes)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t TokenService) Validate(tokenString string, publicKey *rsa.PublicKey) (*CustomClaims, error) {

	// Parse the token with the provided public key
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {

			return nil, errors.New("unexpected signing method")
		}

		return publicKey, nil
	})

	// Check for token parsing errors
	if err != nil {

		return nil, err
	}

	// Check if the token is valid and has not expired
	if token.Valid {

		// If the token is valid, extract and return the claims
		if claims, ok := token.Claims.(*CustomClaims); ok {
			// Check if the token has expired
			if claims.ExpiresAt < time.Now().Unix() {

				return nil, errors.New("token has expired")
			}

			return claims, nil
		}

		return nil, errors.New("invalid token claims")
	}

	return nil, errors.New("invalid token")
}

func (t TokenService) Issue(user *models.User, expiration int, keyRotator KeyRotatorInterface) (string, error) {

	privateKey, err := t.GetPrivateKey(keyRotator)

	if err != nil {
		return "", err
	}

	log.New("Re-Issuing JWT for %s...", user.PhoneNumber).Log()

	userJWT, err := t.Generate(user, expiration, privateKey)

	if err != nil {

		return "", err
	}

	return userJWT, nil
}

func (t TokenService) GetPrivateKey(rotator KeyRotatorInterface) (*rsa.PrivateKey, error) {
	sess, err := session.NewSession(getAWSSessionConfig())

	if err != nil {

		return nil, fmt.Errorf("couldn't create the AWS session")
	}

	// Get a rotator service for accessing the keys in AWS Parameter Store
	rotatorService := NewKeyRotator(sess)
	log.New("Getting the private key from remote...").Log()

	privateKey, err := rotatorService.GetCurrentRSAPrivateKey(ParameterStorePrivateKeyName)

	if err != nil {

		return nil, fmt.Errorf("couldn't get the private key")
	}
	return privateKey, nil
}

func (t TokenService) GetPublicKey(rotator KeyRotatorInterface) (*rsa.PublicKey, error) {
	sess, err := session.NewSession(getAWSSessionConfig())

	if err != nil {

		return nil, fmt.Errorf("couldn't create the AWS session")
	}

	log.New("Getting the public key from remote...").Log()

	// Get a rotator service for accessing the keys in AWS Parameter Store
	rotatorService := NewKeyRotator(sess)
	publicKey, err := rotatorService.GetCurrentRSAPublicKey(ParameterStorePublicKeyName)

	if err != nil {

		return nil, fmt.Errorf("couldn't get the public key")
	}
	return publicKey, nil
}

func getAWSSessionConfig() *aws.Config {

	cfg := aws.NewConfig()

	if EnableAWSSessionDebug {
		cfg.LogLevel = aws.LogLevel(

			aws.LogDebugWithRequestErrors |
				aws.LogDebugWithRequestRetries |
				aws.LogDebugWithSigning |
				aws.LogDebugWithHTTPBody,
		)
	}

	return cfg
}

func CreateCustomClaims(user *models.User, expirationMinutes int) *CustomClaims {

	expiry := time.Minute * time.Duration(expirationMinutes)

	return &CustomClaims{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiry).Unix(),
			Issuer:    Issuer,
			Audience:  Audience,
		},
		UserID:          user.ID,
		Email:           user.Email,
		PhoneNumber:     user.PhoneNumber,
		Firstname:       user.Firstname,
		Lastname:        user.Lastname,
		AccountStatus:   user.AccountStatus.Name,
		AccountStatusID: user.AccountStatusID,
		PhoneVerified:   user.PhoneVerified,
	}
}

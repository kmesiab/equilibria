package jwt

import (
	"crypto/rsa"

	rotator "github.com/kmesiab/go-key-rotator"

	"github.com/aws/aws-sdk-go/aws/session"
)

type KeyRotatorInterface interface {
	GetPublicKeyParameterStoreKeyName() string
	GetPrivateKeyParameterStoreKeyName() string
	GetCurrentRSAPrivateKey(parameterStoreKey string) (*rsa.PrivateKey, error)
	GetCurrentRSAPublicKey(parameterStoreKey string) (*rsa.PublicKey, error)
}

type KeyRotator struct {
	rotator rotator.KeyRotator
}

func NewKeyRotator(sess *session.Session) *KeyRotator {

	ps := rotator.NewAWSParameterStore(sess)
	kr := rotator.NewKeyRotator(ps)
	kri := *kr

	return &KeyRotator{kri}
}

func (r *KeyRotator) GetPublicKeyParameterStoreKeyName() string {
	return ParameterStorePublicKeyName
}
func (r *KeyRotator) GetPrivateKeyParameterStoreKeyName() string {
	return ParameterStorePrivateKeyName
}

func (r *KeyRotator) GetCurrentRSAPrivateKey(parameterStoreKey string) (*rsa.PrivateKey, error) {

	return r.rotator.GetCurrentRSAPrivateKey(parameterStoreKey)
}

func (r *KeyRotator) GetCurrentRSAPublicKey(parameterStoreKey string) (*rsa.PublicKey, error) {
	return r.rotator.GetCurrentRSAPublicKey(parameterStoreKey)
}

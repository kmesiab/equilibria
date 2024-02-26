package jwt

import (
	"crypto/rand"
	"crypto/rsa"
)

var (
	globalMockKeyRotator *MockKeyRotator = nil
)

type MockKeyRotator struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewMockKeyRotator() *MockKeyRotator {
	if globalMockKeyRotator == nil {

		privateKey, publicKey, _ := generateDummyKeyPair()

		globalMockKeyRotator = &MockKeyRotator{
			PrivateKey: privateKey,
			PublicKey:  publicKey,
		}
	}
	return globalMockKeyRotator
}

func (r *MockKeyRotator) GetPublicKeyParameterStoreKeyName() string {
	return ParameterStorePublicKeyName
}

func (r *MockKeyRotator) GetPrivateKeyParameterStoreKeyName() string {
	return ParameterStorePrivateKeyName
}

func (r *MockKeyRotator) GetCurrentRSAPrivateKey(_ string) (*rsa.PrivateKey, error) {
	return r.PrivateKey, nil
}

func (r *MockKeyRotator) GetCurrentRSAPublicKey(_ string) (*rsa.PublicKey, error) {
	return r.PublicKey, nil
}

func generateDummyKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

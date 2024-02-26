package test

import (
	"crypto/rsa"
)

type MockKeyRotator struct{}

func (r *MockKeyRotator) GetCurrentRSAPrivateKey(_ string) (*rsa.PrivateKey, error) {
	return nil, nil
}

func (r *MockKeyRotator) GetCurrentRSAPublicKey(_ string) (*rsa.PublicKey, error) {
	return nil, nil
}

func (r *MockKeyRotator) GetPublicKeyParameterStoreKeyName() string {
	return ""
}
func (r *MockKeyRotator) GetPrivateKeyParameterStoreKeyName() string {
	return ""
}

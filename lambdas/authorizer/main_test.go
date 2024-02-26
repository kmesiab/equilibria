package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/jwt"
)

func IsApproved(response events.APIGatewayCustomAuthorizerResponse) bool {
	return response.PolicyDocument.Statement[0].Effect == "Allow"
}

func TestMain_HandleRequest_DifferentTokens(t *testing.T) {
	publicKey, err := getPublicKey()

	require.NoError(t, err)
	require.NotNil(t, publicKey, "Public key must not be nil")

	handler := AuthorizerLambdaHandler{
		KeyRotator:   jwt.NewMockKeyRotator(),
		TokenService: jwt.TokenService{},
	}

	for name, token := range tokens {
		t.Run(name, func(t *testing.T) {

			response, err := handler.HandleRequest(events.APIGatewayCustomAuthorizerRequest{
				Type:               TokenTypeString,
				AuthorizationToken: token,
				MethodArn:          "arn:aws:execute-api:/../*/POST/",
			})

			require.NoError(t, err, "HandleRequest should not error")

			switch name {
			case "valid":
				assert.True(t, IsApproved(response), "Valid token should be approved")
			case "expired", "invalid", "empty":
				assert.False(t, IsApproved(response), "Expired, invalid, or empty tokens should not be approved")
			}
		})
	}
}

func getPublicKey() (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing the public key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaKey, nil
}

var tokens = map[string]string{
	"empty":   "",
	"invalid": "invalid",
	"expired": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJnby1mb3J0dW5lLXRlbGxlciIsImV4cCI6MTcwNDg3NzA2OSwiaXNzIjoiZ28tZm9ydHVuZS10ZWxsZXIiLCJVc2VySUQiOjE5LCJFbWFpbCI6ImttZXNpYWJAZ21haWwuY29tIiwiUGhvbmVOdW1iZXIiOiIrMTI1MzMyNDMwNzEiLCJGaXJzdG5hbWUiOiJLZXZpbiIsIkxhc3RuYW1lIjoiTWVzaWFiIiwiQWNjb3VudFN0YXR1cyI6IkFjdGl2ZSIsIkFjY291bnRTdGF0dXNJRCI6MiwiUGhvbmVWZXJpZmllZCI6dHJ1ZX0.eHxDu4VLYqv8YgEsvZH4YiXC5gyBXWNT1VY7DIlsIf1u2F3EfAbWyPtRAfrBGeoSKFeZWHYAa3hDGMOYzl53Eadl-cXT19QFNIokvF9a5uDm_PQjjXaVjREaYE_M0GAnNPySo4NVTTLO0-tOh42xLLRSY1w1WpVwNAejEnbW0D3jXVyoYcy4U4VohsLoio-wNcPt7eIWRUHKDfmuabq_idX-GN5kl99N8yec4RWBvO77-KbP73PtvGPp1iLX9Og2O522m_TxddKyK6eYFAN7mvVNE8l1XJfYSN6UMkaa2_0jUnnEWMNyKex5DZFQ-wopARe74rWlVCy_buaPzl-QMQ",
	// "valid":   "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJnby1mb3J0dW5lLXRlbGxlciIsImV4cCI6MTcwNDkyMDgzOCwiaXNzIjoiZ28tZm9ydHVuZS10ZWxsZXIiLCJVc2VySUQiOjE5LCJFbWFpbCI6ImttZXNpYWJAZ21haWwuY29tIiwiUGhvbmVOdW1iZXIiOiIrMTI1MzMyNDMwNzEiLCJGaXJzdG5hbWUiOiJLZXZpbiIsIkxhc3RuYW1lIjoiTWVzaWFiIiwiQWNjb3VudFN0YXR1cyI6IkFjdGl2ZSIsIkFjY291bnRTdGF0dXNJRCI6MiwiUGhvbmVWZXJpZmllZCI6dHJ1ZX0.YI_zgTyYgkUM0L1zm8u-gaord8sAQnrhvysGpNSX_cxxEysaoqyKDhZfkrK_AtwifayPvTGYb-VAxSxjp3Hd88dzIpVk_Cqptx6AT9x7sdfAEFJN6aeb8JhIzbJllAg9g3IDL3FU0ifs3czeiXn90iJXtdyq47UlfueOgPCMqQ0zOzDTUcsxDzJRp9TSaAFP1ddwpaedjN129zxJRvY_V6DxNveP_gkYkE98KI06EX3KvRm0SLMXDqs6b4m5HzwDD9y1gZh6EISEWmNYbWiIbUOeY-8j_jtUD5OgdHgWhtbCBoi8PvWoK-4BtV0GYe_uqqTX6y3Y0zEZJcy4T5TaAg",
}

const publicKey = `-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0BUNk2pH3sZmaG6hH6Um
NTym1qh0kPJZ93zXftDkFRWj1c8rtISqb9DGYHDVOBh9b6xySAtnlmCtCZd/1ZZg
ZdF2DvTT3SD6nqLbeCmtQDZwmZkZUqK9LY1FylV+oHvOpMBRLb4JzKZaWKiPVbIU
w5sw/HKmQPK+OIrKdAqaKlC18iWFo9or+2ZnBxKgxs0XOBnkY2dMhPNZkWawezlu
y20OMLLsK8P3XaY6da8xIYK9FPee4RESmTivUaBv+HKyejGcyV2j+JpkJlBD60cZ
iPFEttKcqqmc90lJaWNrHmi4btc3RQRc8HWyhvkiWNckXTCmipdJyRs3O4uJxlnW
NwIDAQAB
-----END RSA PUBLIC KEY-----
`

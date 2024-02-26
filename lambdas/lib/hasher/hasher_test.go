package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kmesiab/equilibria/lambdas/lib/test"
)

func TestHashPassword(t *testing.T) {
	hashed, err := HashPassword(test.DefaultTestPassword)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, test.DefaultTestPassword, hashed)
}

func TestCheck(t *testing.T) {
	hashed, _ := HashPassword(test.DefaultTestPassword)

	assert.True(t, CheckPassword(test.DefaultTestPassword, hashed))
	assert.False(t, CheckPassword("wrongpassword", hashed))
}

func TestIsSecureString(t *testing.T) {
	assert.False(t, IsSecureString("short"))
	assert.True(t, IsSecureString("longenoughpassword"))
}

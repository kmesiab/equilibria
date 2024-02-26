package config_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/test"
)

func TestConfig_Get(t *testing.T) {

	test.SetEnvVars()

	want := test.GenerateTestConfig()
	got := config.Get()

	assert.NotNil(t, got, "Config should not be nil")
	assert.Equal(t, want.OpenAIAPIKey, got.OpenAIAPIKey)
	assert.Equal(t, want.DatabaseHost, got.DatabaseHost)
	assert.Equal(t, want.DatabaseUser, got.DatabaseUser)
	assert.Equal(t, want.DatabasePassword, got.DatabasePassword)
	assert.Equal(t, want.DatabaseName, got.DatabaseName)
	assert.Equal(t, want.TwilioAuthToken, got.TwilioAuthToken)
	assert.Equal(t, want.TwilioSID, got.TwilioSID)
	assert.Equal(t, want.TwilioPhoneNumber, got.TwilioPhoneNumber)
	assert.Equal(t, strconv.Itoa(want.LogLevel), strconv.Itoa(got.LogLevel)) // Compare as strings
	assert.Equal(t, want.SMSQueueURL, got.SMSQueueURL)
	assert.Equal(t, want.TwilioStatusCallbackURL, got.TwilioStatusCallbackURL)
	assert.Equal(t, want.TwilioVerifyServiceSID, got.TwilioVerifyServiceSID)

}

func TestConfig_IsSingleton(t *testing.T) {
	test.SetEnvVars()

	config1 := config.Get()
	config2 := config.Get()

	assert.Equal(t, config1, config2)
}

func TestConfig_GetMissingEnvironmentVariable(t *testing.T) {
	test.SetEnvVars()
	err := os.Unsetenv("OPENAI_API_KEY")
	require.NoError(t, err, "Error should be nil when unsetting environment variable")

	cfg := config.Get()
	assert.Nil(t, cfg)

	test.SetEnvVars()

}

package twilio

import (
	"net/url"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	"github.com/kmesiab/equilibria/lambdas/lib/test"
)

func TestIsValidWebhookRequest_ValidSignature(t *testing.T) {
	// Set the correct environment variables if needed
	test.SetEnvVars()

	// Token used for generating the signature (Replace with actual token)
	token := "dummy_twilio_auth_token"

	// Sample request body
	bodyValues := url.Values{"Body": []string{"value"}}
	body := bodyValues.Encode()

	validSignature := "TS/8+YTxhA9bhQvFNottGq0LQFg="

	// Create a mock APIGatewayProxyRequest with the valid signature and URL-encoded body
	request := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"X-Twilio-Signature": validSignature,
			"Host":               "foo.com",
			"Path":               "/dev/sms-receive",
		},
		Body: body,
	}

	// Check if the webhook request is valid
	isValid := IsValidWebhookRequest(request, token, false)

	// Assert that the request is valid
	assert.True(t, isValid)
}

func TestIsValidWebhookRequest_InvalidSignature(t *testing.T) {
	test.SetEnvVars()
	request := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"x-twilio-signature": "invalid-signature",
		},
		Body: "Body=value",
	}

	isValid := IsValidWebhookRequest(request, "token", false)

	assert.False(t, isValid)
}

func TestIsValidWebhookRequest_ParseError(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		Body: "Invalid body",
	}

	isValid := IsValidWebhookRequest(request, "token", true)

	assert.False(t, isValid)
}

func TestBuildRequestURLFromProxyRequest(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Host": "example.com",
		},
		Path: "/path",
	}

	constructedURL := BuildRequestURLFromProxyRequest(request)

	assert.Equal(t, "https://example.com/dev/path", constructedURL)
}

func TestBuildRequestURLFromProxyRequest_WithQueryString(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Host": "example.com",
		},
		Path: "/path",
		QueryStringParameters: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	proxyRequest := BuildRequestURLFromProxyRequest(request)

	assert.Equal(t, "https://example.com/dev/path", proxyRequest,
		"Query string should be stripped")
}

func TestBuildRequestURLFromProxyRequest_EncodesQueryString(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Host": "example.com",
		},
		Path: "/path",
		QueryStringParameters: map[string]string{
			"key": "value 1",
		},
	}

	proxyRequest := BuildRequestURLFromProxyRequest(request)

	assert.Equal(t, "https://example.com/dev/path", proxyRequest, "The querystring should be stripped")
}

func TestConvertURLValuesToMap_Empty(t *testing.T) {
	values := url.Values{}
	result := ConvertURLValuesToMap(values)
	assert.Equal(t, 0, len(result))
}

func TestConvertURLValuesToMap_SingleValue(t *testing.T) {
	values := url.Values{"key": []string{"value"}}
	result := ConvertURLValuesToMap(values)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "value", result["key"])
}

func TestConvertURLValuesToMap_MultiValue(t *testing.T) {
	values := url.Values{"key": []string{"value1", "value2"}}
	result := ConvertURLValuesToMap(values)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "value1", result["key"])
}

func TestConvertURLValuesToMap_MultiKey(t *testing.T) {
	values := url.Values{
		"key1": []string{"value1"},
		"key2": []string{"value2"},
	}
	result := ConvertURLValuesToMap(values)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "value2", result["key2"])
}

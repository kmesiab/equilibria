package nrclex_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kmesiab/equilibria/lambdas/lib/nrclex"
)

/*
func TestNRCLexClient_Integration_Get(t *testing.T) {


	client := utils.NewRestClient()
	nrcClient := nrclex.NewNRCLexClient(client.GetClient())
	resp, err := nrcClient.AnalyzeText("I'm in love, this is an amazing day!")

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Greater(t, 0, resp.EmotionScore.Joy)
	require.Greater(t, 1, resp.VaderEmotionScore.Compound)
}
*/

func TestNRCLexClient_AnalyzeText(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that the request is as expected
		require.Equal(t, "POST", r.Method)

		defer func(Body io.ReadCloser) {

			err := Body.Close()
			if err != nil {

				log.Println("Error closing response body in TestNRCLexClient_AnalyzeText")
			}
		}(r.Body)

		body, err := io.ReadAll(r.Body)

		require.NoError(t, err)
		require.Contains(t, string(body), url.QueryEscape("text"))

		// Respond with a mock response
		mockResponse := nrclex.APIResponse{
			Text: "test",
			EmotionScore: nrclex.EmotionScores{
				Anticipation: 1,
				Joy:          2,
				Positive:     3,
				Surprise:     4,
				Trust:        5,
			},
			VaderEmotionScore: nrclex.VaderEmotionScore{
				Compound: 0.5,
				Neg:      0.1,
				Neu:      0.2,
				Pos:      0.7,
			},
		}

		respBytes, err := json.Marshal(mockResponse)
		require.NoError(t, err)

		_, err = w.Write(respBytes)

		require.NoError(t, err, "Error writing mock response")

	}))
	defer server.Close()

	// Use the mock server's URL as the BaseURL
	client := nrclex.NewNRCLexClient(http.DefaultClient)
	client.BaseURL = server.URL

	// Call AnalyzeText
	text := "This is a test."
	response, err := client.AnalyzeText(text)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, "test", response.Text)
	assert.Equal(t, float64(1), response.EmotionScore.Anticipation)
	assert.Equal(t, float64(2), response.EmotionScore.Joy)
}

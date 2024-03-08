package nrclex

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/kmesiab/equilibria/lambdas/lib/utils"
)

const apiURL = "https://langtool.net/sentiment"

type NRCLexClient struct {
	BaseURL string
	Client  utils.SimpleHttpClientInterface
}

func NewNRCLexClient(client utils.SimpleHttpClientInterface) *NRCLexClient {
	return &NRCLexClient{
		BaseURL: apiURL,
		Client:  client,
	}
}

func (n *NRCLexClient) AnalyzeText(text string) (*APIResponse, error) {

	var (
		resp *http.Response
		body []byte
		err  error
	)

	postData := url.Values{}
	postData.Set("text", text)

	if resp, err = n.Client.PostForm(n.BaseURL, postData); err != nil {

		return nil, fmt.Errorf(
			"error posting to %s: %s", n.BaseURL, err.Error(),
		)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()

		// If we throw an error here all we can do is log it.
		if err != nil {

			log.Println("Error closing response body in NRCLexClient.AnalyzeText")
		}
	}(resp.Body)

	if body, err = io.ReadAll(resp.Body); err != nil {

		return nil, err
	}

	if resp.StatusCode != http.StatusOK {

		return n.handleAPIError(body, resp.StatusCode)
	}

	analysis := &APIResponse{}

	if err = json.Unmarshal(body, analysis); err != nil {

		return nil, err
	}

	return analysis, nil
}

func (n *NRCLexClient) handleAPIError(body []byte, statusCode int) (*APIResponse, error) {
	var errorResponse APIErrorResponse

	err := json.Unmarshal(body, &errorResponse)

	if err != nil {

		return nil, fmt.Errorf(
			"error unmarshaling error response: %s", err.Error(),
		)
	}

	return nil, fmt.Errorf(
		"unexpected status code: %d: %s", statusCode, errorResponse.Error,
	)
}

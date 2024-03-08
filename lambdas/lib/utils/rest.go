package utils

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

const Timeout = 6 * time.Second

const RetryCount = 3
const RetryWaitTime = 2 * time.Second
const RetryMaxWaitTime = 6 * time.Second

var RetryExceededError = errors.New("retry count exceeded")

var (
	once        sync.Once
	restyClient *resty.Client
	retryCodes  = map[int]bool{
		http.StatusTooManyRequests:     true,
		http.StatusInternalServerError: true,
		http.StatusBadGateway:          true,
		http.StatusServiceUnavailable:  true,
		http.StatusGatewayTimeout:      true,
	}
)

type SimpleHttpClientInterface interface {
	PostForm(string, url.Values) (*http.Response, error)
}

// NewRestClient returns a singleton resty.Client with retry/backoff logic
// Note: This function and its client are not thread safe.
func NewRestClient() *resty.Client {
	// The sync.Once package provides a mechanism for initializing a value exactly once
	once.Do(func() {
		restyClient = resty.New().
			SetRetryCount(RetryCount).
			SetRetryWaitTime(RetryWaitTime).
			SetTimeout(Timeout).
			SetRetryMaxWaitTime(RetryMaxWaitTime).
			AddRetryCondition(func(resp *resty.Response, err error) bool {
				return err != nil || retryCodes[resp.StatusCode()]
			}).
			SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
				return 0, RetryExceededError
			})
	})

	return restyClient
}

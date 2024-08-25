package helix

import (
	"errors"
	"time"

	"github.com/kvizyx/twitchkit"
	httpcore "github.com/kvizyx/twitchkit/http-core"
)

const (
	// NoRetryInterval is an indicator for the client to do retries
	// one by one without any interval.
	NoRetryInterval = -1
)

type (
	Client struct {
		authProvider twitchkit.AuthProvider
		httpClient   httpcore.HTTPClient
		retryConfig  RetryConfig
	}

	ClientConfig struct {
		AuthProvider twitchkit.AuthProvider
		HTTPClient   httpcore.HTTPClient
		RetryConfig  RetryConfig
	}

	// RetryConfig is a configuration for client's requests retry behavior.
	//
	// Twitch rate-limits reference: https://dev.twitch.tv/docs/api/guide/#twitch-rate-limits
	RetryConfig struct {
		// RetryAll forces all requests made with current client instance to retry
		// if they are need to (rate-limited, http.StatusServiceUnavailable etc.)
		//
		// If you want to retry only specific requests then Client.WithRetry must be used.
		RetryAll bool

		// RetryOnUnavailable forces all requests made with current client instance to retry
		// on http.StatusServiceUnavailable status code even if RetryAll is set to false.
		RetryOnUnavailable bool

		// RetryOnUnavailableTimes is a number of attempts to retry the request that
		// returned http.StatusServiceUnavailable status code at first.
		//
		// By default, only one attempt will be made.
		RetryOnUnavailableTimes uint8

		// RetryOnUnavailableInterval is an interval at which requests that returned
		// http.StatusServiceUnavailable status code at first will be retried.
		// If you want to set zero interval use NoRetryInterval as value.
		//
		// By default, it's one second.
		RetryOnUnavailableInterval time.Duration

		// MaxRateLimitTimeout is a maximum timeout to wait for retry after request was rate-limited.
		// You should also consider that the minimum timeout division is one second (for the safety
		// of comparison with timeout returned by server) so only seconds will be accounted.
		//
		// By default, any rate-limit timeout will be acceptable.
		MaxRateLimitTimeout time.Duration
	}
)

func NewClient(cfg ClientConfig) (*Client, error) {
	if cfg.AuthProvider == nil {
		return nil, errors.New("authentication provider should not be empty")
	}

	// TODO: do smth with default values settings

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = httpcore.DefaultHTTPClient()
	}

	if cfg.RetryConfig.RetryOnUnavailableTimes == 0 {
		cfg.RetryConfig.RetryOnUnavailableTimes = 1
	}

	switch cfg.RetryConfig.RetryOnUnavailableInterval {
	case 0:
		cfg.RetryConfig.RetryOnUnavailableInterval = 1 * time.Second
	case NoRetryInterval:
		cfg.RetryConfig.RetryOnUnavailableInterval = 0
	}

	return &Client{
		authProvider: cfg.AuthProvider,
		httpClient:   cfg.HTTPClient,
		retryConfig:  cfg.RetryConfig,
	}, nil
}

// WithRetry returns copy of the original Client with retry enabled.
// If retry already enabled then no effect will be taken and original
// client will be returned.
func (c Client) WithRetry() Client {
	if c.retryConfig.RetryAll {
		return c
	}

	clientWithRetry := c
	clientWithRetry.retryConfig.RetryAll = true

	return clientWithRetry
}

// AsUser ...
func (c Client) AsUser(userID string, fn any) {
	// TODO: implement me
}

package helix

import (
	"errors"

	"github.com/kvizyx/twitchkit/auth-provider"
	"github.com/kvizyx/twitchkit/http-core"
)

// UserContext ...
type UserContext struct {
	UserID string
}

type (
	Client struct {
		authProvider authprovider.AuthProvider
		httpClient   httpcore.HTTPClient
		userCtx      UserContext
		retryConfig  RetryConfig
	}

	ClientConfig struct {
		AuthProvider authprovider.AuthProvider
		HTTPClient   httpcore.HTTPClient
		RetryConfig  RetryConfig
	}
)

func NewClient(cfg ClientConfig) (*Client, error) {
	if cfg.AuthProvider == nil {
		return nil, errors.New("authentication provider should not be empty")
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = httpcore.DefaultHTTPClient()
	}

	return &Client{
		authProvider: cfg.AuthProvider,
		httpClient:   cfg.HTTPClient,
		retryConfig:  finalizeRetryConfig(cfg.RetryConfig),
	}, nil
}

// WithRetry returns copy of the original Client with retry enabled for all cases.
// This method can be useful if you want to retry only specific requests.
func (c Client) WithRetry() Client {
	clientWithRetry := c
	clientWithRetry.retryConfig.RetryRateLimit = true
	clientWithRetry.retryConfig.RetryUnavailable = true

	return clientWithRetry
}

// WithRetryConfig returns copy of the original Client with provided retry config.
func (c Client) WithRetryConfig(cfg RetryConfig) Client {
	clientWithRetry := c
	clientWithRetry.retryConfig = cfg

	return clientWithRetry
}

// AsUser ...
func (c Client) AsUser(userID string, fn func(client Client)) {
	clientWithUserCtx := c
	clientWithUserCtx.userCtx = UserContext{
		UserID: userID,
	}

	fn(clientWithUserCtx)
}

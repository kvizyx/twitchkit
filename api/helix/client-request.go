package helix

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kvizyx/twitchkit/api"
	httpcore "github.com/kvizyx/twitchkit/http-core"
)

var (
	ErrRetryTimeout = errors.New("retry timeout is greater than it set to be")
)

func (c Client) setAuthHeaders(_ *http.Request) error {
	// TODO: retrieve and set access token from auth provider
	return nil
}

func (c Client) doRequest(req *http.Request, dest any) (api.ResponseMetadata, error) {
	if err := c.setAuthHeaders(req); err != nil {
		return api.ResponseMetadata{}, fmt.Errorf("set authentication headers: %w", err)
	}

	metadata, err := httpcore.DoAPIRequest(req, dest, c.httpClient)
	if err != nil {
		if !errors.Is(err, httpcore.ErrUnsuccessfulRequest) {
			return metadata, err
		}
	}

	switch metadata.StatusCode {
	case http.StatusUnauthorized:
		// TODO: try to retry with access token refreshing
	case http.StatusTooManyRequests:
		if c.retryConfig.RetryAll {
			return c.retryRateLimitedRequest(req, metadata, dest)
		}
	case http.StatusServiceUnavailable:
		if c.retryConfig.RetryAll || c.retryConfig.RetryOnUnavailable {
			return c.retryRequest(
				req,
				dest,
				c.retryConfig.RetryOnUnavailableTimes,
				0,
				c.retryConfig.RetryOnUnavailableInterval,
			)
		}
	}

	return metadata, err
}

func (c Client) retryRateLimitedRequest(
	req *http.Request,
	metadata api.ResponseMetadata,
	dest any,
) (api.ResponseMetadata, error) {
	var (
		serverLimitTimeout = metadata.RateLimitReset() - time.Now().Unix()
		maxLimitTimeout    = int64(c.retryConfig.MaxRateLimitTimeout.Seconds())
	)

	if (serverLimitTimeout >= maxLimitTimeout) && maxLimitTimeout != 0 {
		return metadata, ErrRetryTimeout
	}

	retryAfter := (time.Duration(serverLimitTimeout)) * time.Second

	return c.retryRequest(req, dest, 1, retryAfter, 0)
}

func (c Client) retryRequest(
	req *http.Request,
	dest any,
	times uint8,
	after, interval time.Duration,
) (api.ResponseMetadata, error) {
	if after > 0 {
		retryAfter := time.NewTimer(after)
		<-retryAfter.C
	}

	var waitInterval *time.Ticker
	if interval > 0 {
		waitInterval = time.NewTicker(interval)
	}

	var (
		err      error
		metadata api.ResponseMetadata
	)

	for range times {
		metadata, err = httpcore.DoAPIRequest(req, dest, c.httpClient)
		if err == nil {
			return metadata, nil
		}

		if waitInterval != nil {
			<-waitInterval.C
		}
	}

	return metadata, err
}

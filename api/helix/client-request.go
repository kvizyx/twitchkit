package helix

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
	case http.StatusTooManyRequests:
		if c.retryConfig.RetryAll {
			limitResetHeader := metadata.Header.Get("RateLimit-Reset")

			limitResetValue, err := strconv.ParseInt(limitResetHeader, 10, 64)
			if err != nil {
				return metadata, fmt.Errorf("parse rate-limit reset header: %w", err)
			}

			limitTimeout := limitResetValue - time.Now().Unix()
			if limitTimeout >= int64(c.retryConfig.MaxRateLimitTimeout.Seconds()) {
				return metadata, ErrRetryTimeout
			}

			retryAfter := (time.Duration(limitTimeout)) * time.Second

			return c.retryRequest(req, dest, 1, retryAfter, 0)
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

func (c Client) retryRateLimitRequest() {

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
		if waitInterval != nil {
			<-waitInterval.C
		}

		metadata, err = httpcore.DoAPIRequest(req, dest, c.httpClient)
		if err == nil {
			return metadata, nil
		}
	}

	return metadata, err
}

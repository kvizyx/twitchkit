package helix

import (
	"time"
)

const (
	// NoRetryInterval is an indicator for the client to do retries
	// one by one without any interval.
	NoRetryInterval = -1

	// RetryUntilOK is an indicator for the client to do retries
	// until successful response status.
	RetryUntilOK = -1
)

// RetryConfig is a global configuration for Client's requests retry behavior. If you want to
// retry only specific requests then Client.WithRetry or Client.WithRetryConfig is your
// choice.
//
// Twitch rate-limits reference: https://dev.twitch.tv/docs/api/guide/#twitch-rate-limits
type RetryConfig struct {
	// RetryRateLimit forces all requests made with current client instance to retry
	// if they are rate-limited.
	RetryRateLimit bool

	// RetryUnavailable forces all requests made with current client instance to retry on
	// http.StatusServiceUnavailable response status code (by the way Twitch advises
	// you to do this).
	RetryUnavailable bool

	// RetryUnavailableTimes is a number of attempts to retry the request with http.StatusServiceUnavailable
	// response status code at first. You can use RetryUntilOK as value for this field to retry until
	// successful response status.
	//
	// By default, only one attempt will be made.
	RetryUnavailableTimes int32

	// RetryUnavailableInterval is an interval at which requests with http.StatusServiceUnavailable
	// response status code at first will be retried. You can use NoRetryInterval as value for
	// this field to set zero interval.
	//
	// By default, it's one second.
	RetryUnavailableInterval time.Duration

	// MaxRateLimitTimeout is a maximum timeout to wait for retry after request was rate-limited.
	// You should also consider that the minimum timeout division is one second (for the safety
	// of comparison with timeout returned by server) so only seconds will be accounted.
	//
	// By default, any rate-limit timeout will be acceptable.
	MaxRateLimitTimeout time.Duration
}

var DefaultRetryConfig = RetryConfig{
	RetryRateLimit:           true,
	RetryUnavailable:         true,
	RetryUnavailableTimes:    1,
	RetryUnavailableInterval: 1 * time.Second,
	MaxRateLimitTimeout:      0, // no limit
}

// finalizeRetryConfig returns copy of the original RetryConfig with values
// from DefaultRetryConfig in place of unset or unacceptable values.
func finalizeRetryConfig(cfg RetryConfig) RetryConfig {
	cleanCfg := cfg

	if cfg.RetryUnavailableTimes == 0 {
		cleanCfg.RetryUnavailableTimes = DefaultRetryConfig.RetryUnavailableTimes
	}

	if cfg.RetryUnavailableInterval <= 0 {
		cleanCfg.RetryUnavailableInterval = DefaultRetryConfig.RetryUnavailableInterval
	}

	if cfg.RetryUnavailableInterval == NoRetryInterval {
		cleanCfg.RetryUnavailableInterval = 0
	}

	if cfg.MaxRateLimitTimeout < 0 {
		cfg.MaxRateLimitTimeout = 0
	}

	return cleanCfg
}

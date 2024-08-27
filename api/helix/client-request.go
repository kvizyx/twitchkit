package helix

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kvizyx/twitchkit/api"
	"github.com/kvizyx/twitchkit/api/oauth"
	"github.com/kvizyx/twitchkit/auth-provider"
	"github.com/kvizyx/twitchkit/http-core"
)

var (
	ErrRetryTimeout = errors.New("retry timeout is greater than it set to be")
	ErrAuthNoUserID = errors.New("no user ID was provided for authorized request")
)

// RequestAuthParams ...
type RequestAuthParams struct {
	UserID string
	Scopes []string
}

func (c Client) doRequest(req *http.Request, dest any, authParams RequestAuthParams) (api.ResponseMetadata, error) {
	// specified scopes means that we are forced to do request with user access token.
	if len(authParams.Scopes) != 0 {
		if len(authParams.UserID) == 0 {
			return api.ResponseMetadata{}, ErrAuthNoUserID
		}

		userToken, err := c.authProvider.UserAccessToken(
			req.Context(),
			authParams.UserID,
			authParams.Scopes,
		)
		if err != nil {
			return api.ResponseMetadata{}, fmt.Errorf("get user access token: %w", err)
		}

		if !oauth.IsTokenExpired(&userToken) {
			return c.doAuthorizedRequest(req, dest, &userToken, authParams.UserID)
		}

		freshToken, err := c.tryRefreshUserAccessToken(req.Context(), authParams.UserID)
		if err != nil {
			return api.ResponseMetadata{}, err
		}

		return c.doAuthorizedRequest(req, dest, &freshToken, authParams.UserID)
	}

	ctxUserID := authParams.UserID
	if len(c.userCtx.UserID) != 0 {
		ctxUserID = c.userCtx.UserID
	}

	// if user context is not empty (AsUser method was used) and it's ID exist in
	// provider, then we will use user access token. Otherwise, app access token
	// will be used.
	accessToken, err := c.authProvider.AnyAccessToken(req.Context(), ctxUserID)
	if err != nil {
		return api.ResponseMetadata{}, fmt.Errorf("get any access token: %w", err)
	}

	if len(accessToken.RefreshToken()) != 0 && oauth.IsTokenExpired(accessToken) {
		freshToken, err := c.tryRefreshUserAccessToken(req.Context(), ctxUserID)
		if err != nil {
			return api.ResponseMetadata{}, err
		}

		return c.doAuthorizedRequest(req, dest, &freshToken, ctxUserID)
	}

	return c.doAuthorizedRequest(req, dest, accessToken, ctxUserID)
}

func (c Client) doAuthorizedRequest(
	req *http.Request,
	dest any,
	accessToken oauth.AccessToken,
	userID string,
) (api.ResponseMetadata, error) {
	authType := c.authProvider.AuthorizationType()

	api.SetClientHeader(req, c.authProvider.ClientID())
	api.SetAuthHeader(req, authType, accessToken.AccessToken())

	metadata, err := httpcore.DoAPIRequest(req, dest, c.httpClient)
	if err != nil {
		if !errors.Is(err, httpcore.ErrUnsuccessfulRequest) {
			return metadata, err
		}
	}

	switch metadata.StatusCode {
	case http.StatusUnauthorized:
		if len(accessToken.RefreshToken()) == 0 {
			appToken, err := c.authProvider.AppAccessToken(req.Context(), true)
			if err != nil {
				return api.ResponseMetadata{}, err
			}

			api.SetAuthHeader(req, authType, appToken.AccessToken())

			return c.retryRequest(req, dest, 1, 0, 0)
		}

		freshToken, err := c.tryRefreshUserAccessToken(req.Context(), userID)
		if err != nil {
			return metadata, err
		}

		api.SetAuthHeader(req, authType, freshToken.AccessToken())

		return c.retryRequest(req, dest, 1, 0, 0)
	case http.StatusTooManyRequests:
		if c.retryConfig.RetryRateLimit {
			return c.retryRateLimitedRequest(req, metadata, dest)
		}
	case http.StatusServiceUnavailable:
		if c.retryConfig.RetryUnavailable {
			return c.retryRequest(
				req,
				dest,
				c.retryConfig.RetryUnavailableTimes,
				0,
				c.retryConfig.RetryUnavailableInterval,
			)
		}
	}

	return metadata, err
}

func (c Client) tryRefreshUserAccessToken(
	ctx context.Context,
	userID string,
) (oauth.UserAccessToken, error) {
	refresher, ok := c.authProvider.(authprovider.RefreshProvider)
	if !ok {
		return oauth.UserAccessToken{}, authprovider.ErrNotRefresher
	}

	freshToken, err := refresher.RefreshUserAccessToken(ctx, userID)
	if err != nil {
		return oauth.UserAccessToken{}, fmt.Errorf("refresh user access token: %w", err)
	}

	return freshToken, nil
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

// retryRequest stupidly retry request with given options. No exp. backoff or
// some other smart algorithm, just one by one requests.
func (c Client) retryRequest(
	req *http.Request,
	dest any,
	times int32,
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
		retried  int32
		err      error
		metadata api.ResponseMetadata
	)

	for {
		metadata, err = httpcore.DoAPIRequest(req, dest, c.httpClient)
		if err == nil {
			return metadata, nil
		}

		retried += 1
		if retried == times && times != RetryUntilOK {
			break
		}

		if waitInterval != nil {
			<-waitInterval.C
		}
	}

	return metadata, err
}

package twitchkit

import (
	"context"
	"net/http"
	"net/url"

	"github.com/kvizyx/twitchkit/api"
	httpcore "github.com/kvizyx/twitchkit/http-core"
)

type AppTokenResponse struct {
	AppToken
	ResponseMetadata api.ResponseMetadata
}

type AccessTokenResponse struct {
	AccessToken
	ResponseMetadata api.ResponseMetadata
}

type ValidateTokenResponse struct {
	ClientID  string   `json:"client_id"`
	Login     string   `json:"login"`
	Scopes    []string `json:"scopes"`
	UserID    string   `json:"user_id"`
	ExpiresIn int64    `json:"expires_in"`

	ResponseMetadata api.ResponseMetadata
}

type ExchangeCodeParams struct {
	ClientCredentials
	Code        string
	RedirectURL string
}

func ExchangeCode(
	ctx context.Context,
	params ExchangeCodeParams,
	httpClient ...httpcore.HTTPClient,
) (AccessTokenResponse, error) {
	const resource = "token"

	values := url.Values{}
	values.Set("client_id", params.ClientID)
	values.Set("client_secret", params.ClientSecret)
	values.Set("code", params.Code)
	values.Set("grant_type", "authorization_code")
	values.Set("redirect_url", params.RedirectURL)

	req, err := httpcore.OAuthRequestWithURLValues(ctx, resource, http.MethodPost, values, true)
	if err != nil {
		return AccessTokenResponse{}, err
	}

	var accessToken AccessTokenResponse

	metadata, err := httpcore.DoAPIRequest(req, &accessToken, httpClient...)
	accessToken.ResponseMetadata = metadata

	if err != nil {
		return accessToken, err
	}

	accessToken.ObtainedAt().SetNow(true)

	return accessToken, nil
}

func FetchAppAccessToken(
	ctx context.Context,
	credentials ClientCredentials,
	httpClient ...httpcore.HTTPClient,
) (AppTokenResponse, error) {
	const resource = "token"

	values := url.Values{}
	values.Set("client_id", credentials.ClientID)
	values.Set("client_secret", credentials.ClientSecret)
	values.Set("grant_type", "client_credentials")

	req, err := httpcore.OAuthRequestWithURLValues(ctx, resource, http.MethodPost, values, true)
	if err != nil {
		return AppTokenResponse{}, err
	}

	var appToken AppTokenResponse

	metadata, err := httpcore.DoAPIRequest(req, &appToken, httpClient...)
	appToken.ResponseMetadata = metadata

	if err != nil {
		return appToken, err
	}

	appToken.ObtainedAt().SetNow(true)

	return appToken, nil
}

func RefreshToken(
	ctx context.Context,
	credentials ClientCredentials,
	refreshToken string,
	httpClient ...httpcore.HTTPClient,
) (AccessTokenResponse, error) {
	const resource = "token"

	values := url.Values{}
	values.Set("client_id", credentials.ClientID)
	values.Set("client_secret", credentials.ClientSecret)
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", refreshToken)

	req, err := httpcore.OAuthRequestWithURLValues(ctx, resource, http.MethodPost, values, true)
	if err != nil {
		return AccessTokenResponse{}, err
	}

	var accessToken AccessTokenResponse

	metadata, err := httpcore.DoAPIRequest(req, &accessToken, httpClient...)
	accessToken.ResponseMetadata = metadata

	if err != nil {
		return accessToken, err
	}

	accessToken.ObtainedAt().SetNow(true)

	return accessToken, nil
}

func RevokeToken(
	ctx context.Context,
	clientID, accessToken string,
	httpClient ...httpcore.HTTPClient,
) (api.ResponseMetadata, error) {
	const resource = "revoke"

	values := url.Values{}
	values.Set("client_id", clientID)
	values.Set("token", accessToken)

	req, err := httpcore.OAuthRequestWithURLValues(ctx, resource, http.MethodPost, values, true)
	if err != nil {
		return api.ResponseMetadata{}, err
	}

	metadata, err := httpcore.DoAPIRequest(req, nil, httpClient...)

	return metadata, err
}

func ValidateToken(
	ctx context.Context,
	accessToken string,
	httpClient ...httpcore.HTTPClient,
) (ValidateTokenResponse, error) {
	const resource = "validate"

	req, err := httpcore.OAuthRequestEmpty(ctx, resource, http.MethodGet)
	if err != nil {
		return ValidateTokenResponse{}, err
	}

	api.SetAuthType(req, api.AuthTypeOAuth, accessToken)

	var validateToken ValidateTokenResponse

	metadata, err := httpcore.DoAPIRequest(req, &validateToken, httpClient...)
	validateToken.ResponseMetadata = metadata

	return validateToken, err
}

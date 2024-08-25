package twitchkit

import (
	"context"
	"net/http"
	"net/url"

	"github.com/kvizyx/twitchkit/api"
	httpcore "github.com/kvizyx/twitchkit/http-core"
)

type AppAccessTokenResponse struct {
	AppAccessToken
	ResponseMetadata api.ResponseMetadata
}

type UserAccessTokenResponse struct {
	UserAccessToken
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
) (UserAccessTokenResponse, error) {
	const resource = "token"

	values := url.Values{}
	values.Set("client_id", params.ClientID)
	values.Set("client_secret", params.ClientSecret)
	values.Set("code", params.Code)
	values.Set("grant_type", "authorization_code")
	values.Set("redirect_url", params.RedirectURL)

	req, err := httpcore.NewAPIRequest(ctx, httpcore.RequestOptions{
		APIScope: api.ScopeOAuth,
		Resource: resource,
		Method:   http.MethodPost,
		Body:     values,
	}, false)
	if err != nil {
		return UserAccessTokenResponse{}, err
	}

	var accessToken UserAccessTokenResponse

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
) (AppAccessTokenResponse, error) {
	const resource = "token"

	values := url.Values{}
	values.Set("client_id", credentials.ClientID)
	values.Set("client_secret", credentials.ClientSecret)
	values.Set("grant_type", "client_credentials")

	req, err := httpcore.NewAPIRequest(ctx, httpcore.RequestOptions{
		APIScope: api.ScopeOAuth,
		Resource: resource,
		Method:   http.MethodPost,
		Body:     values,
	}, false)
	if err != nil {
		return AppAccessTokenResponse{}, err
	}

	var appToken AppAccessTokenResponse

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
) (UserAccessTokenResponse, error) {
	const resource = "token"

	values := url.Values{}
	values.Set("client_id", credentials.ClientID)
	values.Set("client_secret", credentials.ClientSecret)
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", refreshToken)

	req, err := httpcore.NewAPIRequest(ctx, httpcore.RequestOptions{
		APIScope: api.ScopeOAuth,
		Resource: resource,
		Method:   http.MethodPost,
		Body:     values,
	}, false)
	if err != nil {
		return UserAccessTokenResponse{}, err
	}

	var accessToken UserAccessTokenResponse

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

	req, err := httpcore.NewAPIRequest(ctx, httpcore.RequestOptions{
		APIScope: api.ScopeOAuth,
		Resource: resource,
		Method:   http.MethodPost,
		Body:     values,
	}, false)
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
) (bool, ValidateTokenResponse, error) {
	const resource = "validate"

	req, err := httpcore.NewAPIRequest(ctx, httpcore.RequestOptions{
		APIScope: api.ScopeOAuth,
		Resource: resource,
		Method:   http.MethodGet,
	}, false)
	if err != nil {
		return false, ValidateTokenResponse{}, err
	}

	api.SetAuthHeader(req, api.AuthTypeOAuth, accessToken)

	var vt ValidateTokenResponse

	metadata, err := httpcore.DoAPIRequest(req, &vt, httpClient...)
	vt.ResponseMetadata = metadata

	return vt.ResponseMetadata.StatusCode == http.StatusOK, vt, err
}

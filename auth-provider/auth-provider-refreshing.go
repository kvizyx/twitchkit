package authprovider

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/kvizyx/twitchkit/api"
	"github.com/kvizyx/twitchkit/api/oauth"
)

const authorizationType = api.AuthTypeBearer

type (
	// OnRefreshCallback triggers on successful user access token refresh event.
	OnRefreshCallback func(userID string, token oauth.UserAccessToken)

	// OnRefreshFailureCallback triggers on failed user access token refresh event
	OnRefreshFailureCallback func(userID string, err error)
)

type RefreshingProvider struct {
	clientID     string
	clientSecret string
	redirectURI  string
	scopes       []string

	appAccessToken oauth.AppAccessToken
	appTokenLocker sync.RWMutex

	users       map[string]oauth.UserAccessToken
	usersLocker sync.RWMutex

	cbOnRefresh        OnRefreshCallback
	cbOnRefreshFailure OnRefreshFailureCallback
}

var (
	_ AuthProvider    = &RefreshingProvider{}
	_ RefreshProvider = &RefreshingProvider{}
)

type RefreshingProviderParams struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

func NewRefreshingProvider(p RefreshingProviderParams) *RefreshingProvider {
	return &RefreshingProvider{
		clientID:     p.ClientID,
		clientSecret: p.ClientSecret,
		redirectURI:  p.RedirectURI,
		scopes:       p.Scopes,
		users:        make(map[string]oauth.UserAccessToken),
	}
}

func (ap *RefreshingProvider) ClientID() string {
	return ap.clientID
}

func (ap *RefreshingProvider) AuthorizationType() api.AuthorizationType {
	return authorizationType
}

// OnRefresh sets callback for successful UserAccessToken refresh event.
func (ap *RefreshingProvider) OnRefresh(cb OnRefreshCallback) {
	ap.cbOnRefresh = cb
}

// OnRefreshFailure sets callback for failed UserAccessToken refresh event.
func (ap *RefreshingProvider) OnRefreshFailure(cb OnRefreshFailureCallback) {
	ap.cbOnRefreshFailure = cb
}

// AppAccessToken ...
func (ap *RefreshingProvider) AppAccessToken(ctx context.Context, forceNew bool) (oauth.AppAccessToken, error) {
	ap.appTokenLocker.RLock()
	if !forceNew && !oauth.IsTokenExpired(&ap.appAccessToken) && len(ap.appAccessToken.AccessToken()) > 0 {
		ap.appTokenLocker.RUnlock()
		return ap.appAccessToken, nil
	}
	ap.appTokenLocker.RUnlock()

	res, err := oauth.FetchAppAccessToken(ctx, oauth.ClientCredentials{
		ClientID:     ap.clientID,
		ClientSecret: ap.clientSecret,
	})
	if err != nil {
		return oauth.AppAccessToken{}, fmt.Errorf("fetch app access token: %w", err)
	}

	ap.appTokenLocker.Lock()
	ap.appAccessToken = res.AppAccessToken
	ap.appTokenLocker.Unlock()

	return ap.appAccessToken, nil
}

// UserAccessToken ...
func (ap *RefreshingProvider) UserAccessToken(
	ctx context.Context,
	userID string,
	scopes []string,
) (oauth.UserAccessToken, error) {
	accessToken, found := ap.user(userID)
	if !found {
		return oauth.UserAccessToken{}, ErrUserNotFound
	}

	absentScope, equal := oauth.IsScopesEqual(accessToken.Scope(), scopes)
	if !equal {
		return oauth.UserAccessToken{}, oauth.MissingScopeError(absentScope)
	}

	if len(accessToken.AccessToken()) != 0 && !oauth.IsTokenExpired(&accessToken) {
		return accessToken, nil
	}

	freshToken, err := ap.RefreshUserAccessToken(ctx, userID)
	if err != nil {
		return oauth.UserAccessToken{}, fmt.Errorf("refresh user access token: %w", err)
	}

	return freshToken, nil
}

// AnyAccessToken ...
func (ap *RefreshingProvider) AnyAccessToken(ctx context.Context, userID string) (oauth.AccessToken, error) {
	if len(userID) != 0 && ap.HasUser(userID) {
		userToken, err := ap.UserAccessToken(ctx, userID, nil)
		if err != nil {
			return nil, fmt.Errorf("get user access token: %w", err)
		}

		return &userToken, nil
	}

	appToken, err := ap.AppAccessToken(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("get app access token: %w", err)
	}

	return &appToken, nil
}

// AddUserForToken ...
func (ap *RefreshingProvider) AddUserForToken(ctx context.Context, accessToken oauth.AccessToken) (string, error) {
	tokenWithInfo := oauth.AccessTokenWithInfo{
		AnyAccessToken: accessToken,
	}

	if len(accessToken.AccessToken()) != 0 && !oauth.IsTokenExpired(accessToken) {
		_, res, err := oauth.ValidateToken(ctx, accessToken.AccessToken())
		if err != nil {
			if res.ResponseMetadata.StatusCode != http.StatusUnauthorized {
				return "", fmt.Errorf("validate token: %w", err)
			}
		}

		tokenWithInfo.TokenInfo = res.TokenInfo
	}

	if len(tokenWithInfo.TokenInfo.ClientID) == 0 {
		if len(accessToken.RefreshToken()) == 0 {
			return "", ErrEmptyRefresh
		}

		freshTokenWithInfo, err := ap.RefreshUnknownUserAccessToken(ctx, accessToken.RefreshToken())
		if err != nil {
			return "", fmt.Errorf("refresh access token for unknown user: %w", err)
		}

		tokenWithInfo = freshTokenWithInfo
	}

	if len(tokenWithInfo.TokenInfo.UserID) == 0 {
		return "", oauth.ErrUnsuitableToken
	}

	userToken, ok := tokenWithInfo.AnyAccessToken.(*oauth.UserAccessToken)
	if !ok {
		return "", oauth.ErrUnsuitableToken
	}

	if len(tokenWithInfo.Scope()) == 0 {
		userToken.ScopeValue = tokenWithInfo.TokenInfo.Scopes
	}

	if err := ap.AddUser(tokenWithInfo.TokenInfo.UserID, *userToken); err != nil {
		return "", fmt.Errorf("add user: %w", err)
	}

	return tokenWithInfo.TokenInfo.UserID, nil
}

// AddUserForCode ...
func (ap *RefreshingProvider) AddUserForCode(ctx context.Context, code string) (string, error) {
	if len(ap.redirectURI) == 0 {
		return "", oauth.ErrEmptyRedirectURI
	}

	res, err := oauth.ExchangeCode(ctx, oauth.ExchangeCodeParams{
		ClientCredentials: oauth.ClientCredentials{
			ClientID:     ap.clientID,
			ClientSecret: ap.clientSecret,
		},
		Code:        code,
		RedirectURI: ap.redirectURI,
	})
	if err != nil {
		return "", fmt.Errorf("exchange code: %w", err)
	}

	userID, err := ap.AddUserForToken(ctx, &res.UserAccessToken)
	if err != nil {
		return "", fmt.Errorf("add user for token: %w", err)
	}

	return userID, nil
}

// AddUser ...
func (ap *RefreshingProvider) AddUser(userID string, token oauth.UserAccessToken) error {
	if len(token.RefreshToken()) == 0 {
		return ErrEmptyRefresh
	}

	ap.usersLocker.Lock()
	ap.users[userID] = token
	ap.usersLocker.Unlock()

	return nil
}

// RemoveUser ...
func (ap *RefreshingProvider) RemoveUser(userID string) {
	ap.usersLocker.Lock()
	delete(ap.users, userID)
	ap.usersLocker.Unlock()
}

// HasUser returns if provider contains token for given user ID.
func (ap *RefreshingProvider) HasUser(userID string) bool {
	_, found := ap.user(userID)
	return found
}

// RefreshUnknownUserAccessToken ...
func (ap *RefreshingProvider) RefreshUnknownUserAccessToken(
	ctx context.Context,
	refreshToken string,
) (oauth.AccessTokenWithInfo, error) {
	var token oauth.AccessTokenWithInfo

	err := ap.withRefreshCallbacks(func() error {
		var err error

		token, err = ap.refreshUnknownUserToken(ctx, refreshToken)
		if err != nil {
			return err
		}

		return nil
	}, token, "")
	if err != nil {
		return oauth.AccessTokenWithInfo{}, err
	}

	return token, nil
}

// RefreshUserAccessToken ...
func (ap *RefreshingProvider) RefreshUserAccessToken(
	ctx context.Context,
	userID string,
) (oauth.UserAccessToken, error) {
	var token oauth.UserAccessToken

	err := ap.withRefreshCallbacks(func() error {
		var err error

		token, err = ap.refreshUserToken(ctx, userID)
		if err != nil {
			return err
		}

		return nil
	}, &token, userID)
	if err != nil {
		return oauth.UserAccessToken{}, err
	}

	return token, nil
}

func (ap *RefreshingProvider) refreshUnknownUserToken(
	ctx context.Context,
	refreshToken string,
) (oauth.AccessTokenWithInfo, error) {
	if len(refreshToken) == 0 {
		return oauth.AccessTokenWithInfo{}, ErrEmptyRefresh
	}

	freshToken, err := oauth.RefreshToken(ctx,
		oauth.ClientCredentials{
			ClientID:     ap.clientID,
			ClientSecret: ap.clientSecret,
		},
		refreshToken,
	)
	if err != nil {
		return oauth.AccessTokenWithInfo{}, fmt.Errorf("refresh token: %w", err)
	}

	_, tokenInfo, err := oauth.ValidateToken(ctx, freshToken.AccessToken())
	if err != nil {
		return oauth.AccessTokenWithInfo{}, fmt.Errorf("validate token: %w", err)
	}

	if !ap.HasUser(tokenInfo.UserID) {
		return oauth.AccessTokenWithInfo{}, ErrUserNotFound
	}

	_ = ap.AddUser(tokenInfo.UserID, freshToken.UserAccessToken)

	return oauth.AccessTokenWithInfo{
		AnyAccessToken: &freshToken.UserAccessToken,
		TokenInfo:      tokenInfo.TokenInfo,
	}, nil
}

func (ap *RefreshingProvider) refreshUserToken(
	ctx context.Context,
	userID string,
) (oauth.UserAccessToken, error) {
	oldToken, userExist := ap.user(userID)
	if !userExist {
		return oauth.UserAccessToken{}, ErrUserNotFound
	}

	if len(oldToken.RefreshToken()) == 0 {
		return oauth.UserAccessToken{}, ErrEmptyRefresh
	}

	freshToken, err := oauth.RefreshToken(ctx,
		oauth.ClientCredentials{
			ClientID:     ap.clientID,
			ClientSecret: ap.clientSecret,
		},
		oldToken.RefreshToken(),
	)
	if err != nil {
		return oauth.UserAccessToken{}, fmt.Errorf("refresh token: %w", err)
	}

	if !ap.HasUser(userID) {
		return oauth.UserAccessToken{}, ErrUserNotFound
	}

	_ = ap.AddUser(userID, freshToken.UserAccessToken)

	return freshToken.UserAccessToken, nil
}

func (ap *RefreshingProvider) withRefreshCallbacks(
	refresh func() error,
	token oauth.AccessToken,
	userID string,
) error {
	if err := refresh(); err != nil {
		if ap.cbOnRefreshFailure == nil {
			return err
		}

		if len(userID) == 0 {
			tokenInfo, ok := token.(oauth.AccessTokenWithInfo)
			if !ok {
				return err
			}

			if len(tokenInfo.TokenInfo.UserID) != 0 {
				go ap.cbOnRefreshFailure(tokenInfo.TokenInfo.UserID, err)
			}
			return err
		}

		go ap.cbOnRefreshFailure(userID, err)
		return err
	}

	if ap.cbOnRefresh != nil {
		userToken, ok := token.(*oauth.UserAccessToken)
		if !ok {
			return oauth.ErrUnsuitableToken
		}

		go ap.cbOnRefresh(userID, *userToken)
	}

	return nil
}

func (ap *RefreshingProvider) user(userID string) (oauth.UserAccessToken, bool) {
	ap.usersLocker.RLock()
	token, found := ap.users[userID]
	ap.usersLocker.RUnlock()

	return token, found
}

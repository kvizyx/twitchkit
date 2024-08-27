package authprovider

import (
	"context"

	"github.com/kvizyx/twitchkit/api"
	"github.com/kvizyx/twitchkit/api/oauth"
)

// TODO: make AuthProvider interface more flexible

// AuthProvider ...
type AuthProvider interface {
	ClientID() string
	AuthorizationType() api.AuthorizationType
	AnyAccessToken(ctx context.Context, userID string) (oauth.AccessToken, error)
	UserAccessToken(ctx context.Context, userID string, scopes []string) (oauth.UserAccessToken, error)
	AppAccessToken(ctx context.Context, forceNew bool) (oauth.AppAccessToken, error)
}

// RefreshProvider ...
type RefreshProvider interface {
	// RefreshUserAccessToken refreshes user access token for given user ID and save it
	// internally. It would be ideal to implement callbacks on successful and failed
	// refresh events.
	RefreshUserAccessToken(ctx context.Context, userID string) (oauth.UserAccessToken, error)
}

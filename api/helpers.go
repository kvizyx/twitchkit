package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	HelixBaseURL = "https://api.twitch.tv/helix"
	OAuthBaseURL = "https://id.twitch.tv/oauth2"
)

var (
	ErrUnknownScope = errors.New("unknown API scope was provided")
)

type AuthorizationType string

const (
	AuthTypeBearer AuthorizationType = "Bearer"
	AuthTypeOAuth  AuthorizationType = "OAuth"
)

type Scope int

const (
	ScopeHelix Scope = iota
	ScopeOAuth
)

func ComposeHelixURL(resource string) string {
	return strings.Join([]string{HelixBaseURL, resource}, "/")
}

func ComposeOAuthURL(resource string) string {
	return strings.Join([]string{OAuthBaseURL, resource}, "/")
}

// SetAuthHeader sets authorization header for the provided HTTP request.
func SetAuthHeader(req *http.Request, authType AuthorizationType, accessToken string) {
	if req == nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", authType, accessToken))
}

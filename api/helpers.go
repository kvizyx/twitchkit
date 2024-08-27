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
	ErrUnknownType = errors.New("unknown Twitch API type")
)

// Type ...
type Type int

const (
	TypeHelix Type = iota
	TypeOAuth
)

// ComposeHelixURL ...
func ComposeHelixURL(resource string) string {
	return strings.Join([]string{HelixBaseURL, resource}, "/")
}

// ComposeOAuthURL ...
func ComposeOAuthURL(resource string) string {
	return strings.Join([]string{OAuthBaseURL, resource}, "/")
}

// SetClientHeader sets app client ID header for the provider HTTP request.
func SetClientHeader(req *http.Request, clientID string) {
	if req == nil {
		return
	}

	req.Header.Set("Client-ID", clientID)
}

// SetAuthHeader sets Twitch API authorization header for the provided HTTP request.
func SetAuthHeader(req *http.Request, authType AuthorizationType, accessToken string) {
	if req == nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", authType, accessToken))
}

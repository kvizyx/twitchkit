package api

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	HelixBaseURL = "https://api.twitch.tv/helix"
	OAuthBaseURL = "https://id.twitch.tv/oauth2"
)

type AuthorizationType string

const (
	AuthTypeBearer AuthorizationType = "Bearer"
	AuthTypeOAuth  AuthorizationType = "OAuth"
)

// ResponseMetadata is metadata from Twitch API HTTP response.
type ResponseMetadata struct {
	StatusCode    int
	Header        http.Header
	TwitchError   string `json:"error"`
	TwitchStatus  int    `json:"status"`
	TwitchMessage string `json:"message"`
}

func ComposeHelixURL(resource string) string {
	return strings.Join([]string{HelixBaseURL, resource}, "/")
}

func ComposeOAuthURL(resource string) string {
	return strings.Join([]string{OAuthBaseURL, resource}, "/")
}

func SetAuthType(req *http.Request, authType AuthorizationType, accessToken string) {
	if req == nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", authType, accessToken))
}

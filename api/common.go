package api

import (
	"net/http"
	"strconv"
)

// AuthorizationType ...
type AuthorizationType string

const (
	AuthTypeBearer    AuthorizationType = "Bearer"
	AuthTypeOAuth     AuthorizationType = "OAuth"
	AuthTypeExtension AuthorizationType = "Extension"
)

// ResponseMetadata is metadata from Twitch API HTTP response.
type ResponseMetadata struct {
	StatusCode    int
	Header        http.Header
	TwitchError   string `json:"error"`
	TwitchStatus  int    `json:"status"`
	TwitchMessage string `json:"message"`
}

func (rm ResponseMetadata) RateLimit() int {
	value, _ := strconv.Atoi(rm.Header.Get("RateLimit-Limit"))
	return value
}

func (rm ResponseMetadata) RateLimitRemaining() int {
	value, _ := strconv.Atoi(rm.Header.Get("RateLimit-Remaining"))
	return value
}

func (rm ResponseMetadata) RateLimitReset() int64 {
	value, _ := strconv.ParseInt(rm.Header.Get("RateLimit-Remaining"), 10, 64)
	return value
}

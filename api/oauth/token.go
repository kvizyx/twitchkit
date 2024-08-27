package oauth

import (
	"time"
)

// ObtainTime is time in Unix timestamp when token was obtained.
type ObtainTime int64

// ExpirationToken describes token that have a lifetime.
type ExpirationToken interface {
	ExpiresIn() int64
	ObtainedAt() *ObtainTime
}

type TokenLifetime struct {
	ExpiresInValue  int64 `json:"expires_in"`
	ObtainedAtValue *ObtainTime
}

func (tl *TokenLifetime) ExpiresIn() int64 {
	return tl.ExpiresInValue
}

func (tl *TokenLifetime) ObtainedAt() *ObtainTime {
	if tl.ObtainedAtValue == nil {
		tl.ObtainedAtValue = new(ObtainTime)
	}

	return tl.ObtainedAtValue
}

// IsTokenExpired returns whether token is expired or not with 30 seconds
// spare for safer refresh handling.
func IsTokenExpired(token ExpirationToken) bool {
	return time.Now().Unix() >= (token.ObtainedAt().Int64() + token.ExpiresIn() + 30)
}

// Int64 returns ObtainTime as int64 type.
func (ot *ObtainTime) Int64() int64 {
	if ot == nil {
		return 0
	}

	return int64(*ot)
}

// SetNow sets ObtainTime to the current time in Unix timestamp
// if it doesn't already have a value or force set to true.
func (ot *ObtainTime) SetNow(force bool) {
	if ot == nil {
		return
	}

	if *ot > 0 && !force {
		return
	}

	*ot = ObtainTime(time.Now().Unix())
}

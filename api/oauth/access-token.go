package oauth

// AccessToken is a Twitch API authorization token.
type AccessToken interface {
	AccessToken() string
	RefreshToken() string
	Scope() []string
	ExpirationToken
}

// AccessTokenWithInfo is combination of AccessToken with TokenInfo.
// It may be useful for contextualized operations with access token.
type AccessTokenWithInfo struct {
	AnyAccessToken AccessToken
	TokenInfo      TokenInfo
}

func (at AccessTokenWithInfo) AccessToken() string {
	return at.AnyAccessToken.AccessToken()
}

func (at AccessTokenWithInfo) RefreshToken() string {
	return at.AnyAccessToken.RefreshToken()
}

func (at AccessTokenWithInfo) Scope() []string {
	return at.AnyAccessToken.Scope()
}

func (at AccessTokenWithInfo) ExpiresIn() int64 {
	return at.AnyAccessToken.ExpiresIn()
}

func (at AccessTokenWithInfo) ObtainedAt() *ObtainTime {
	return at.AnyAccessToken.ObtainedAt()
}

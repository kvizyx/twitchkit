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
	return at.RefreshToken()
}

func (at AccessTokenWithInfo) Scope() []string {
	return at.Scope()
}

func (at AccessTokenWithInfo) ExpiresIn() int64 {
	return at.ExpiresIn()
}

func (at AccessTokenWithInfo) ObtainedAt() *ObtainTime {
	return at.ObtainedAt()
}

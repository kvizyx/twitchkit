package oauth

// UserAccessToken ...
type UserAccessToken struct {
	AccessTokenValue  string   `json:"access_token"`
	RefreshTokenValue string   `json:"refresh_token"`
	ScopeValue        []string `json:"scope"`
	TokenLifetime
}

var _ AccessToken = &UserAccessToken{}

func (ut UserAccessToken) AccessToken() string {
	return ut.AccessTokenValue
}

func (ut UserAccessToken) RefreshToken() string {
	return ut.RefreshTokenValue
}

func (ut UserAccessToken) Scope() []string {
	return ut.ScopeValue
}

package oauth

// AppAccessToken ...
type AppAccessToken struct {
	AccessTokenValue string `json:"access_token"`
	TokenLifetime
}

var _ AccessToken = &AppAccessToken{}

func (at AppAccessToken) AccessToken() string {
	return at.AccessTokenValue
}

func (at AppAccessToken) RefreshToken() string {
	return ""
}

func (at AppAccessToken) Scope() []string {
	return nil
}

package twitchkit

var (
	_ ExpirationToken = &AppToken{}
	_ ExpirationToken = &AccessToken{}
)

type AccessToken struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenLifetime
}

type AppToken struct {
	AccessToken string `json:"access_token"`
	TokenLifetime
}

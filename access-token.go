package twitchkit

var (
	_ ExpirationToken = &AppAccessToken{}
	_ ExpirationToken = &UserAccessToken{}
)

type UserAccessToken struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenLifetime
}

type AppAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenLifetime
}

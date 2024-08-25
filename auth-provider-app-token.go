package twitchkit

// AppTokenAuthProvider ...
type AppTokenAuthProvider struct {
	clientID     string
	clientSecret string
}

var _ AuthProvider = &AppTokenAuthProvider{}

func NewAppTokenAuthProvider(clientID, clientSecret string) *AppTokenAuthProvider {
	return &AppTokenAuthProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (ap *AppTokenAuthProvider) ClientID() string {
	return ap.clientID
}

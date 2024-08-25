package twitchkit

// StaticAuthProvider ...
type StaticAuthProvider struct {
	clientID    string
	accessToken string
	scopes      []string
}

var _ AuthProvider = &StaticAuthProvider{}

func NewStaticAuthProvider(clientID, accessToken string) *StaticAuthProvider {
	return &StaticAuthProvider{
		clientID:    clientID,
		accessToken: accessToken,
	}
}

func (ap StaticAuthProvider) ClientID() string {
	return ap.clientID
}

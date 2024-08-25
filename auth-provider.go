package twitchkit

// AuthProvider ...
type AuthProvider interface {
	// ClientID returns app client ID.
	ClientID() string

	// ...
}

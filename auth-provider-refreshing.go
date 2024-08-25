package twitchkit

type (
	// OnRefreshCallback ...
	OnRefreshCallback func(userID string, token UserAccessToken)

	// OnRefreshFailureCallback ...
	OnRefreshFailureCallback func(userID string, err error)
)

// RefreshingAuthProvider ...
type RefreshingAuthProvider struct {
	clientID     string
	clientSecret string
	redirectURL  string
	scopes       []string

	cbOnRefresh        OnRefreshCallback
	cbOnRefreshFailure OnRefreshFailureCallback
}

var _ AuthProvider = &RefreshingAuthProvider{}

type RefreshingAuthProviderParams struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

func NewRefreshingAuthProvider(p RefreshingAuthProviderParams) *RefreshingAuthProvider {
	return &RefreshingAuthProvider{
		clientID:     p.ClientID,
		clientSecret: p.ClientSecret,
		redirectURL:  p.RedirectURL,
		scopes:       p.Scopes,
	}
}

func (ap *RefreshingAuthProvider) ClientID() string {
	return ap.clientID
}

// OnRefresh sets callback for successful tokens refreshing event.
func (ap *RefreshingAuthProvider) OnRefresh(cb OnRefreshCallback) {
	ap.cbOnRefresh = cb
}

// OnRefreshFailure sets callback for failed tokens refreshing event.
func (ap *RefreshingAuthProvider) OnRefreshFailure(cb OnRefreshFailureCallback) {
	ap.cbOnRefreshFailure = cb
}

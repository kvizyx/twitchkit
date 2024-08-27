package oauth

// equalScopes is an equal authorization scopes.
var equalScopes = map[string][]string{
	"channel_commercial":    {"channel:edit:commercial"},
	"channel_editor":        {"channel:manage:broadcast"},
	"channel_read":          {"channel:read:stream_key"},
	"channel_subscriptions": {"channel:read:subscriptions"},
	"user_blocks_read":      {"user:read:blocked_users"},
	"user_blocks_edit":      {"user:manage:blocked_users"},
	"user_follows_edit":     {"user:edit:follows"},
	"user_read":             {"user:read:email"},
	"user_subscriptions":    {"user:read:subscriptions"},
	"user:edit:broadcast":   {"channel:manage:broadcast", "channel:manage:extensions"},
}

// IsScopesEqual compares given scopes and returns whether scopes are equal.
// If some of requested scopes are missing in given scopes, then the first
// absent scope will be also returned.
func IsScopesEqual(scopes, requestedScopes []string) (string, bool) {
	if len(requestedScopes) <= 0 {
		return "", true
	}

	compareSet := make(map[string]struct{})

	for _, scope := range scopes {
		compareSet[scope] = struct{}{}

		if equal, ok := equalScopes[scope]; ok {
			for _, equalScope := range equal {
				compareSet[equalScope] = struct{}{}
			}
		}
	}

	for _, requestedScope := range requestedScopes {
		if _, exists := compareSet[requestedScope]; !exists {
			// absent scope
			return requestedScope, false
		}
	}

	return "", true
}

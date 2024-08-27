package oauth

import (
	"errors"
	"fmt"
)

var (
	ErrUnsuitableToken  = errors.New("access token is not suitable for this context")
	ErrEmptyRedirectURI = errors.New("redirect URI is empty")
	ErrMissingScope     = errors.New("missing scope but it's required")
)

// MissingScopeError ...
func MissingScopeError(absentScope string) error {
	return fmt.Errorf("%w (%s)", ErrMissingScope, absentScope)
}

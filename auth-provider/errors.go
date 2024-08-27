package authprovider

import (
	"errors"
)

var (
	ErrUserNotFound = errors.New("user with given ID is not in provider")
	ErrEmptyRefresh = errors.New("refresh token is empty")
	ErrNotRefresher = errors.New("cannot refresh user access token as provider is not implement it")
)

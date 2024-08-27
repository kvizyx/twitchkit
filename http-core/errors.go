package httpcore

import (
	"errors"
	"fmt"
)

var (
	ErrUnsuccessfulRequest = errors.New("unsuccessful request")
)

func UnsuccessfulRequestError(status string) error {
	return fmt.Errorf("%w: %s", ErrUnsuccessfulRequest, status)
}

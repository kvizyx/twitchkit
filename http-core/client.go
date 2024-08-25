package httpcore

import (
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func DefaultHTTPClient() HTTPClient {
	return http.DefaultClient
}

func GetOrDefaultHTTPClient(clients ...HTTPClient) HTTPClient {
	if clients == nil || len(clients) == 0 {
		return DefaultHTTPClient()
	}

	return clients[0]
}

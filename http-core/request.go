package httpcore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kvizyx/twitchkit/api"
)

const (
	lastSuccessfulStatus = 299
	lastParseableStatus  = 499
)

var (
	ErrUnknownBody          = errors.New("unknown body type")
	ErrEmptyRequest         = errors.New("request data is empty")
	ErrNoContentDestination = errors.New("no content was returned by server but destination is not empty")
)

type RequestOptions struct {
	APIScope  api.Scope
	Resource  string
	Method    string
	URLValues url.Values
	Body      any
}

func NewAPIRequest(ctx context.Context, opts RequestOptions, jsonBody bool) (*http.Request, error) {
	var endpointURL string

	switch opts.APIScope {
	case api.ScopeHelix:
		endpointURL = api.ComposeHelixURL(opts.Resource)
	case api.ScopeOAuth:
		endpointURL = api.ComposeOAuthURL(opts.Resource)
	default:
		return nil, api.ErrUnknownScope
	}

	if opts.URLValues != nil {
		endpointURL = fmt.Sprintf("%s?%s", endpointURL, opts.URLValues.Encode())
	}

	var body io.Reader

	req, err := http.NewRequestWithContext(ctx, opts.Method, endpointURL, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if opts.Body == nil {
		return req, nil
	}

	if jsonBody {
		jsonBytes, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}

		body = bytes.NewBuffer(jsonBytes)
		req.Header.Set("Content-Type", "application/json")
	} else {
		urlValues, ok := opts.Body.(url.Values)
		if !ok {
			return nil, ErrUnknownBody
		}

		body = bytes.NewBuffer([]byte(urlValues.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	req.Body = io.NopCloser(body)

	return req, nil
}

func DoAPIRequest(req *http.Request, dest any, httpClient ...HTTPClient) (api.ResponseMetadata, error) {
	client := GetOrDefaultHTTPClient(httpClient...)

	res, err := client.Do(req)
	if err != nil {
		return api.ResponseMetadata{}, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	metadata := api.ResponseMetadata{
		StatusCode: res.StatusCode,
		Header:     res.Header,
	}

	if res.StatusCode == http.StatusNoContent {
		if dest != nil {
			return metadata, ErrNoContentDestination
		}

		return metadata, nil
	}

	if res.StatusCode > lastParseableStatus {
		return metadata, UnsuccessfulRequest(res.Status)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return metadata, fmt.Errorf("read response body: %w", err)
	}

	if res.StatusCode > lastSuccessfulStatus {
		if err = json.Unmarshal(bodyBytes, &metadata); err != nil {
			return metadata, fmt.Errorf("unmarshal response body: %w", err)
		}

		return metadata, UnsuccessfulRequest(res.Status)
	}

	if dest != nil {
		if err = json.Unmarshal(bodyBytes, dest); err != nil {
			return metadata, fmt.Errorf("unmarshal response body: %w", err)
		}
	}

	return metadata, nil
}

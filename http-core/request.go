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
	ErrEmptyRequest         = errors.New("request data is empty")
	ErrNoContentDestination = errors.New("no content was returned by server but destination is not empty")
)

// HelixRequestWithURLValues ...
func HelixRequestWithURLValues(
	ctx context.Context,
	resource, method string,
	values url.Values,
	withBody bool,
) (*http.Request, error) {
	return requestWithURLValues(ctx, api.ComposeHelixURL(resource), method, values, withBody)
}

// HelixRequestWithJSON ...
func HelixRequestWithJSON(ctx context.Context, resource, method string, data any) (*http.Request, error) {
	return requestWithJSON(ctx, api.ComposeHelixURL(resource), method, data)
}

// OAuthRequestWithURLValues ...
func OAuthRequestWithURLValues(
	ctx context.Context,
	resource, method string,
	values url.Values,
	withBody bool,
) (*http.Request, error) {
	return requestWithURLValues(ctx, api.ComposeOAuthURL(resource), method, values, withBody)
}

// OAuthRequestEmpty ...
func OAuthRequestEmpty(ctx context.Context, resource, method string) (*http.Request, error) {
	return requestEmpty(ctx, api.ComposeOAuthURL(resource), method)
}

// DoAPIRequest ...
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

func requestEmpty(ctx context.Context, endpoint, method string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request with context: %w", err)
	}

	return req, nil
}

func requestWithURLValues(
	ctx context.Context,
	endpoint, method string,
	values url.Values,
	withBody bool,
) (*http.Request, error) {
	if !withBody {
		endpoint = fmt.Sprintf("%s?%s", endpoint, values.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request with context: %w", err)
	}

	if withBody {
		req.Body = io.NopCloser(bytes.NewReader([]byte(values.Encode())))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

func requestWithJSON(ctx context.Context, endpoint, method string, data any) (*http.Request, error) {
	if data == nil {
		return nil, ErrEmptyRequest
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create request with context: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if len(body) > 0 {
		req.Body = io.NopCloser(bytes.NewReader(body))
	}

	return req, nil
}

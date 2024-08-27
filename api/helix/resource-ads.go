package helix

import (
	"context"
	"net/http"

	"github.com/kvizyx/twitchkit/api"
	"github.com/kvizyx/twitchkit/http-core"
)

type AdsResource struct {
	client Client
}

func (c Client) Ads() AdsResource {
	return AdsResource{client: c}
}

type (
	StartCommercialWrapper struct {
		Data []StartCommercialOutput `json:"data"`
	}

	StartCommercialInput struct {
		BroadcasterID string `json:"broadcaster_id"`
		Length        int    `json:"length"`
	}

	StartCommercialOutput struct {
		Length           int    `json:"length"`
		Message          string `json:"message"`
		RetryAfter       int    `json:"retry_after"`
		ResponseMetadata api.ResponseMetadata
	}
)

// StartCommercial starts a commercial on the specified channel.
//
// Reference: https://dev.twitch.tv/docs/api/reference/#start-commercial
//
// Requires a user access token that includes the channel:edit:commercial scope.
func (r AdsResource) StartCommercial(ctx context.Context, input StartCommercialInput) (StartCommercialOutput, error) {
	const resource = "channels/commercial"

	req, err := httpcore.NewAPIRequest(ctx, httpcore.RequestOptions{
		APIType:  api.TypeHelix,
		Resource: resource,
		Method:   http.MethodPost,
		Body:     input,
	}, true)
	if err != nil {
		return StartCommercialOutput{}, err
	}

	var (
		wrapper StartCommercialWrapper
		output  StartCommercialOutput
	)

	metadata, err := r.client.doRequest(req, &wrapper, RequestAuthParams{
		UserID: input.BroadcasterID,
		Scopes: []string{"channel:edit:commercial"},
	})
	output.ResponseMetadata = metadata

	if err != nil {
		return output, err
	}

	if len(wrapper.Data) != 0 {
		output = wrapper.Data[0]
	}

	return output, nil
}

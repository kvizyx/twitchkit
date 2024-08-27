package helix

import (
	"context"
	"net/http"

	"github.com/kvizyx/twitchkit/api"
	"github.com/kvizyx/twitchkit/http-core"
)

type ChatResource struct {
	client Client
}

func (c Client) Chat() ChatResource {
	return ChatResource{client: c}
}

type (
	GetGlobalChatBadgesOutput struct {
		ChatBadges       []ChatBadge `json:"data"`
		ResponseMetadata api.ResponseMetadata
	}

	ChatBadge struct {
		SetID    string         `json:"set_id"`
		Versions []BadgeVersion `json:"versions"`
	}

	BadgeVersion struct {
		ID          string `json:"id"`
		ImageURL1x  string `json:"image_url_1x"`
		ImageURL2x  string `json:"image_url_2x"`
		ImageURL4x  string `json:"image_url_4x"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ClickAction string `json:"click_action"`
		ClickURL    string `json:"click_url"`
	}
)

// GetGlobalBadges gets Twitch’s list of chat badges, which users may use in any channel’s chat room.
//
// Reference: https://dev.twitch.tv/docs/api/reference/#get-global-chat-badges
//
// Requires an app access token or user access token.
func (r ChatResource) GetGlobalBadges(ctx context.Context) (GetGlobalChatBadgesOutput, error) {
	const resource = "chat/badges/global"

	req, err := httpcore.NewAPIRequest(ctx, httpcore.RequestOptions{
		APIType:  api.TypeHelix,
		Resource: resource,
		Method:   http.MethodGet,
	}, false)
	if err != nil {
		return GetGlobalChatBadgesOutput{}, err
	}

	var output GetGlobalChatBadgesOutput

	metadata, err := r.client.doRequest(req, &output, RequestAuthParams{})
	output.ResponseMetadata = metadata

	if err != nil {
		return output, err
	}

	return output, nil
}

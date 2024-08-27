package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kvizyx/twitchkit/api/helix"
	"github.com/kvizyx/twitchkit/api/oauth"
	"github.com/kvizyx/twitchkit/auth-provider"
)

func main() {
	authProvider := authprovider.NewRefreshingProvider(
		authprovider.RefreshingProviderParams{
			ClientID:     "",
			ClientSecret: "",
			RedirectURL:  "",
			Scopes:       []string{},
		},
	)

	authProvider.OnRefresh(func(userID string, token oauth.UserAccessToken) {
		fmt.Printf("refresh success for %s: %v\n", userID, token)
	})

	authProvider.OnRefreshFailure(func(userID string, err error) {
		fmt.Printf("refresh failure for %s: %s\n", userID, err)
	})

	client, err := helix.NewClient(helix.ClientConfig{
		AuthProvider: authProvider,
		RetryConfig:  helix.DefaultRetryConfig,
	})
	if err != nil {
		log.Fatalf("failed to create helix client: %s", err)
	}

	client.AsUser("", func(client helix.Client) {
		// TODO: do something
	})

	output, err := client.Ads().StartCommercial(
		context.Background(),
		helix.StartCommercialInput{
			BroadcasterID: "",
			Length:        0,
		},
	)
	if err != nil {
		log.Printf("failed to start commercial: %s\n", err)
	}

	fmt.Printf("Output: %v\n", output)
}

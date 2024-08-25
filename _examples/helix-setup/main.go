package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kvizyx/twitchkit"
	"github.com/kvizyx/twitchkit/api/helix"
)

func main() {
	authProvider := twitchkit.NewRefreshingAuthProvider(
		twitchkit.RefreshingAuthProviderParams{
			ClientID:     "<ClientID>",
			ClientSecret: "<ClientSecret>",
			RedirectURL:  "<RedirectURL>",
			Scopes:       []string{},
		},
	)

	authProvider.OnRefresh(func(userID string, token twitchkit.UserAccessToken) {
		fmt.Printf("refresh success for %s: %v\n", userID, token)
	})

	authProvider.OnRefreshFailure(func(userID string, err error) {
		fmt.Printf("refresh failure for %s: %s\n", userID, err)
	})

	client, err := helix.NewClient(helix.ClientConfig{
		AuthProvider: authProvider,
		RetryConfig: helix.RetryConfig{
			RetryAll:                   false,
			RetryOnUnavailable:         true,
			RetryOnUnavailableTimes:    3,
			RetryOnUnavailableInterval: 1 * time.Second,
		},
	})
	if err != nil {
		log.Fatalf("failed to create helix client: %s", err)
	}

	output, err := client.WithRetry().Ads().StartCommercial(
		context.TODO(),
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

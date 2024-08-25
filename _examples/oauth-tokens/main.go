package main

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchkit"
)

func main() {
	res, err := twitchkit.FetchAppAccessToken(
		context.TODO(),
		twitchkit.ClientCredentials{
			ClientID:     "<ClientID>",
			ClientSecret: "<ClientSecret>",
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("App token: %v\n", res.AppAccessToken)

	isValid, validation, err := twitchkit.ValidateToken(context.TODO(), res.AccessToken)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Validate result (%t): %v\n", isValid, validation)
}

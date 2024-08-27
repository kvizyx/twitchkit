package main

import (
	"context"
	"fmt"

	"github.com/kvizyx/twitchkit/api/oauth"
)

func main() {
	res, err := oauth.FetchAppAccessToken(
		context.TODO(),
		oauth.ClientCredentials{
			ClientID:     "<ClientID>",
			ClientSecret: "<ClientSecret>",
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("App token: %v\n", res.AppAccessToken)

	isValid, tokenInfo, err := oauth.ValidateToken(context.TODO(), res.AccessTokenValue)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Validate result (%t): %v\n", isValid, tokenInfo)
}

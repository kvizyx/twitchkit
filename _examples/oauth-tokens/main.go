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
			ClientID:     "lrbwj516b8l8nqm2fr7ft2xeihik8v",
			ClientSecret: "q1zn43p9triqfrhq98u5v2u3zy11aq",
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("App token: %v\n", res.AppToken)

	validateRes, err := twitchkit.ValidateToken(context.TODO(), res.AccessToken)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Validate result: %v\n", validateRes)
}

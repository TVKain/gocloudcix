package main

import (
	"context"
	"fmt"

	gocloudcix "github.com/TVKain/go-cloudcix"
	"github.com/TVKain/go-cloudcix/cloudcix"
)

func main() {
	opts := gocloudcix.AuthOptions{
		Username:           "",
		Password:           "",
		APIKey:             "",
		MembershipEndpoint: "",
	}

	fmt.Print(opts)

	cloudcixClient, _ := cloudcix.AuthenticatedClient(context.TODO(), opts)

	token := cloudcixClient.GetToken()

	fmt.Print(token)
}

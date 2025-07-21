package tokens

import (
	"context"
	"fmt"

	gocloudcix "github.com/TVKain/go-cloudcix"
)

type GetTokenRequest struct {
	// The API Key
	APIKey string
	// The username
	Username string
	// The password
	Password string
}

func Get(ctx context.Context, client *gocloudcix.ApplicationClient, request GetTokenRequest) (*Token, error) {
	opts := &gocloudcix.RequestOpts{
		JSONBody:     request,
		JSONResponse: &Token{},
	}

	response, err := client.Post(ctx, client.Endpoint+"token/", request, nil, opts)
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return nil, fmt.Errorf("%v", response.Error())
	}

	tokenResponse := opts.JSONResponse.(*Token)
	return tokenResponse, nil
}

package tokens

import (
	"context"
	"fmt"

	gocloudcix "github.com/TVKain/go-cloudcix"
)

// GetTokenRequest represents the request structure for getting a token.
// It includes the API key, username, and password.
// This structure is used to authenticate and retrieve a token from the membership service.
// This needs to be converted to a json object when making the request.
type GetTokenRequest struct {
	// The API Key
	APIKey string `json:"api_key"`
	// The username
	Username string `json:"email"`
	// The password
	Password string `json:"password"`
}

const authLoginPath = "/auth/login/"

func Get(ctx context.Context, client *gocloudcix.ApplicationClient, request GetTokenRequest) (*Token, error) {
	opts := &gocloudcix.RequestOpts{
		JSONBody:     request,
		JSONResponse: &Token{},
	}

	// TODO: Use other library for constructing the URL
	// This is a simple string concatenation, but it should be replaced with a proper URL
	url := client.Endpoint + authLoginPath

	response, err := client.Post(ctx, url, request, nil, opts)
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return nil, fmt.Errorf("%v", response.Error())
	}

	tokenResponse := opts.JSONResponse.(*Token)
	return tokenResponse, nil
}

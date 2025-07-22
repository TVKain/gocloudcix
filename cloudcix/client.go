package cloudcix

import (
	"context"
	"fmt"
	"net/url"

	gocloudcix "github.com/TVKain/go-cloudcix"
	tokens "github.com/TVKain/go-cloudcix/cloudcix/membership/tokens"
)

type Application struct {
	name string
}

// List of supported applications.
// To add a new application, just add it here.
var (
	Membership = Application{name: "membership"}
	Compute    = Application{name: "compute"}
	SCM        = Application{name: "scm"}
)

// TODO: Need a better way to handle this.
// Utility function to construct application endpoints
// Example: "https://dev.cloudcix.net/"
// Need a function to construct the full endpoint for an application
// Example: "https://membership.dev.cloudcix.net/"
// Example: "https://compute.dev.cloudcix.net/"
func constructApplicationEndpoint(baseEndpoint string, application Application) string {
	// Parse the base URL
	parsed, err := url.Parse(baseEndpoint)
	if err != nil {
		// If parsing fails, just return baseEndpoint as fallback
		return baseEndpoint
	}

	// Insert the app name before the existing host
	// e.g. dev.cloudcix.net => membership.dev.cloudcix.net
	newHost := application.name + "." + parsed.Host

	// Update the host in the parsed URL
	parsed.Host = newHost

	// Ensure no trailing path remains
	parsed.Path = ""

	return parsed.String()
}

// AuthenticatedClient logs in to an CloudCIX cloud
// specified by the options, acquires a token, and returns a Provider Client
// instance that's ready to operate.
// Example:
//
//	ao, err := cloudcix.AuthOptionsFromEnv()
//	provider, err := cloudcix.AuthenticatedClient(ctx, ao)
//	client, err := cloudcix.NewProject(ctx, provider)
func AuthenticatedClient(ctx context.Context, options gocloudcix.AuthOptions) (*gocloudcix.ProviderClient, error) {
	client := gocloudcix.NewProviderClient(options.BaseEndpoint)

	// Set up reauth function - this is safe because ProviderClient handles concurrent reauth attempts

	err := Authenticate(ctx, client, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Authenticate authenticates or re-authenticates against the most
// recent identity service supported at the provided endpoint.
func Authenticate(ctx context.Context, client *gocloudcix.ProviderClient, options gocloudcix.AuthOptions) error {
	// Create membership client with proper locking

	membershipClient, err := NewMembership(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to create membership client: %w", err)
	}

	client.ReauthFunc = func(ctx context.Context) error {
		newClient, err := AuthenticatedClient(ctx, options)
		if err != nil {
			return fmt.Errorf("failed to re-authenticate: %w", err)
		}
		client.SetToken(newClient.GetToken())
		return nil
	}

	// Get new token
	token, err := tokens.Get(ctx, membershipClient, tokens.GetTokenRequest{
		APIKey:   options.APIKey,
		Username: options.Username,
		Password: options.Password,
	})
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	// Validate token
	if token == nil || token.Token == "" {
		return fmt.Errorf("received empty token from authentication request")
	}

	// Set token - this is thread-safe because ProviderClient.SetToken is protected by a mutex
	client.SetToken(token.Token)
	return nil
}

func NewMembership(ctx context.Context, client *gocloudcix.ProviderClient) (*gocloudcix.ApplicationClient, error) {
	// TODO: Handle protocol setting based on client configuration
	endpoint := constructApplicationEndpoint(client.BaseEndpoint, Membership)

	return &gocloudcix.ApplicationClient{
		ProviderClient: client,
		Endpoint:       endpoint,
	}, nil
}

func NewCompute(ctx context.Context, client *gocloudcix.ProviderClient) (*gocloudcix.ApplicationClient, error) {
	// TODO: Handle protocol setting based on client configuration
	endpoint := constructApplicationEndpoint(client.BaseEndpoint, Compute)

	return &gocloudcix.ApplicationClient{
		ProviderClient: client,
		Endpoint:       endpoint,
	}, nil
}

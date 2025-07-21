package cloudcix

import (
	"context"
	"fmt"

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

// getApplicationEndpoint constructs the full endpoint for an application.
// It supports baseEndpoint as a full FQDN (e.g., "membership.dev.cloudcix.net"),
// a partial domain (e.g., "dev.cloudcix.net"), or even just a hostname/IP.
// If baseEndpoint already includes the application as a subdomain, it is returned as-is.
func getApplicationEndpoint(baseEndpoint string, application Application) (string, error) {
	if baseEndpoint == "" {
		return "", fmt.Errorf("baseEndpoint cannot be empty")
	}

	// Check if the baseEndpoint already contains the application as a subdomain
	if len(baseEndpoint) > 0 && baseEndpoint[:len(application.name)+1] == application.name+"." {
		return baseEndpoint, nil
	}

	// Construct the full endpoint
	return fmt.Sprintf("%s.%s", application.name, baseEndpoint), nil
}

func setApplicationEndpointProtocol(endpoint string, protocol string) string {
	return fmt.Sprintf("%s://%s", protocol, endpoint)
}

// NewClient prepares an unauthenticated ProviderClient instance.
// Most users will probably prefer using the AuthenticatedClient function
// instead.
//
// This is useful if you wish to explicitly control the version of the identity
// service that's used for authentication explicitly, for example.
//
// A basic example of using this would be:
//
//	ao, err := cloudcix.AuthOptionsFromEnv()
//	provider, err := cloudcix.NewClient(ao.MembershipEndpoint)
//	client, err := cloudcix.NewMembership(ctx, provider)
func NewClient(baseEndpoint string) (*gocloudcix.ProviderClient, error) {

	// Use the constructor function instead of manually creating
	client := gocloudcix.NewProviderClient(baseEndpoint)

	// Additional setup if needed
	client.BaseEndpoint = baseEndpoint

	return client, nil
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
	client, err := NewClient(options.MembershipEndpoint)
	if err != nil {
		return nil, err
	}

	err = Authenticate(ctx, client, options)
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

	// Set up reauth function - this is safe because ProviderClient handles concurrent reauth attempts
	client.ReauthFunc = func(ctx context.Context) error {
		return Authenticate(ctx, client, options)
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
	endpoint := setApplicationEndpointProtocol(client.BaseEndpoint, "https")

	return &gocloudcix.ApplicationClient{
		ProviderClient: client,
		Endpoint:       endpoint,
	}, nil
}

package gocloudcix

/*
AuthOptions stores information needed to authenticate to a CloudCIX Cloud.
You can populate one manually, or use a provider's AuthOptionsFromEnv() function
to read relevant information from the standard environment variables. Pass one
to a provider's AuthenticatedClient function to authenticate and obtain a
ProviderClient representing an active session on that provider.

Its fields are the union of those recognized by each identity implementation and
provider.

An example of manually providing authentication information:

	opts := gocloudcix.AuthOptions{
	  BaseEndpoint: "https://api.cloudcix.com/",
	  Username: "{username}",
	  Password: "{password}",
	  ApiKey: "{apikey}",
	}

	provider, err := cloudcix.AuthenticatedClient(context.TODO(), opts)

An example of using AuthOptionsFromEnv(), where the environment variables can
be read from a file, such as a standard openrc file:

	opts, err := cloudcix.AuthOptionsFromEnv()
	provider, err := cloudcix.AuthenticatedClient(context.TODO(), opts)
*/
type AuthOptions struct {
	// CloudCIX API endpoint specifies the HTTP endpoint that is required to work with
	// the Membership API.
	BaseEndpoint string `json:"-"`

	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`

	// ApiKey allows users to authenticate (possibly as another user) with an
	// authentication token ID.
	APIKey string `json:"-"`
}

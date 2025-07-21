package cloudcix

import (
	"fmt"
	"os"

	gocloudcix "github.com/TVKain/go-cloudcix"
)

// nilOptions for returning an empty struct when something is missing
var nilOptions = gocloudcix.AuthOptions{}

// AuthOptionsFromEnv reads CLOUDCIX_* environment variables and returns AuthOptions.
//
// Required variables:
//   - CLOUDCIX_USERNAME
//   - CLOUDCIX_PASSWORD
//   - CLOUDCIX_API_KEY
//
// Example usage:
//
//	opts, err := cloudcix.AuthOptionsFromEnv()
//	if err != nil { log.Fatal(err) }
//	client := cloudcix.NewClient(opts)
func AuthOptionsFromEnv() (gocloudcix.AuthOptions, error) {
	username := os.Getenv("CLOUDCIX_USERNAME")
	password := os.Getenv("CLOUDCIX_PASSWORD")
	apiKey := os.Getenv("CLOUDCIX_API_KEY")

	// Validate required values
	if username == "" {
		return nilOptions, fmt.Errorf("missing required environment variable: CLOUDCIX_USERNAME")
	}
	if password == "" {
		return nilOptions, fmt.Errorf("missing required environment variable: CLOUDCIX_PASSWORD")
	}
	if apiKey == "" {
		return nilOptions, fmt.Errorf("missing required environment variable: CLOUDCIX_API_KEY")
	}

	opts := gocloudcix.AuthOptions{
		Username: username,
		Password: password,
		APIKey:   apiKey,
	}

	return opts, nil
}

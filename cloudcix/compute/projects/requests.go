package projects

import (
	"context"
	"fmt"

	gocloudcix "github.com/TVKain/go-cloudcix"
)

// Requests for interacting with projects in the CloudCIX Compute API
func List(ctx context.Context, client *gocloudcix.ApplicationClient) (*ProjectsListResponse, error) {
	opts := &gocloudcix.RequestOpts{
		JSONResponse: &ProjectsListResponse{},
	}

	fmt.Println("Listing projects at:", client.Endpoint)

	response, err := client.Get(ctx, client.Endpoint+"/project/", nil, opts)
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return nil, fmt.Errorf("%v", response.Error())
	}

	return opts.JSONResponse.(*ProjectsListResponse), nil
}

func Get(ctx context.Context, client *gocloudcix.ApplicationClient, projectID int) (*ProjectGetResponse, error) {
	opts := &gocloudcix.RequestOpts{
		JSONResponse: &ProjectGetResponse{},
	}

	response, err := client.Get(ctx, fmt.Sprintf("%s/project/%d/", client.Endpoint, projectID), nil, opts)
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return nil, fmt.Errorf("%v", response.Error())
	}

	return opts.JSONResponse.(*ProjectGetResponse), nil
}

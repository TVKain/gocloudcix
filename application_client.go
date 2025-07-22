package gocloudcix

import (
	"context"
	"io"
	"maps"

	"github.com/go-resty/resty/v2"
)

// ApplicationClient stores details required to interact with a specific Application API implemented by a provider.
type ApplicationClient struct {
	// ProviderClient is a reference to the provider that implements this Application.
	*ProviderClient

	// Endpoint is the base URL of the Application's API
	Endpoint string

	// MoreHeaders allows users to set Application-wide headers on requests.
	MoreHeaders map[string]string
}

// RequestOpts contains options for making HTTP requests
type RequestOpts struct {
	JSONBody     any
	JSONResponse any
	RawBody      io.Reader
	MoreHeaders  map[string]string
	OKCodes      []int
}

func (client *ApplicationClient) initReqOpts(JSONBody any, JSONResponse any, opts *RequestOpts) {
	if v, ok := (JSONBody).(io.Reader); ok {
		opts.RawBody = v
	} else if JSONBody != nil {
		opts.JSONBody = JSONBody
	}

	if JSONResponse != nil {
		opts.JSONResponse = JSONResponse
	}
}

// Get calls `Request` with the "GET" HTTP verb.
func (client *ApplicationClient) Get(ctx context.Context, url string, JSONResponse any, opts *RequestOpts) (*resty.Response, error) {
	if opts == nil {
		opts = new(RequestOpts)
	}
	client.initReqOpts(nil, JSONResponse, opts)
	return client.Request(ctx, "GET", url, opts)
}

// Post calls `Request` with the "POST" HTTP verb.
func (client *ApplicationClient) Post(ctx context.Context, url string, JSONBody any, JSONResponse any, opts *RequestOpts) (*resty.Response, error) {
	if opts == nil {
		opts = new(RequestOpts)
	}
	client.initReqOpts(JSONBody, JSONResponse, opts)
	return client.Request(ctx, "POST", url, opts)
}

// Put calls `Request` with the "PUT" HTTP verb.
func (client *ApplicationClient) Put(ctx context.Context, url string, JSONBody any, JSONResponse any, opts *RequestOpts) (*resty.Response, error) {
	if opts == nil {
		opts = new(RequestOpts)
	}
	client.initReqOpts(JSONBody, JSONResponse, opts)
	return client.Request(ctx, "PUT", url, opts)
}

// Patch calls `Request` with the "PATCH" HTTP verb.
func (client *ApplicationClient) Patch(ctx context.Context, url string, JSONBody any, JSONResponse any, opts *RequestOpts) (*resty.Response, error) {
	if opts == nil {
		opts = new(RequestOpts)
	}
	client.initReqOpts(JSONBody, JSONResponse, opts)
	return client.Request(ctx, "PATCH", url, opts)
}

// Delete calls `Request` with the "DELETE" HTTP verb.
func (client *ApplicationClient) Delete(ctx context.Context, url string, opts *RequestOpts) (*resty.Response, error) {
	if opts == nil {
		opts = new(RequestOpts)
	}
	client.initReqOpts(nil, nil, opts)
	return client.Request(ctx, "DELETE", url, opts)
}

// Head calls `Request` with the "HEAD" HTTP verb.
func (client *ApplicationClient) Head(ctx context.Context, url string, opts *RequestOpts) (*resty.Response, error) {
	if opts == nil {
		opts = new(RequestOpts)
	}
	client.initReqOpts(nil, nil, opts)
	return client.Request(ctx, "HEAD", url, opts)
}

// Request carries out the HTTP operation for the Application client
func (client *ApplicationClient) Request(ctx context.Context, method, url string, opts *RequestOpts) (*resty.Response, error) {
	if opts == nil {
		opts = new(RequestOpts)
	}
	if opts.MoreHeaders == nil {
		opts.MoreHeaders = make(map[string]string)
	}

	// Merge client-level headers with request headers
	if len(client.MoreHeaders) > 0 {
		maps.Copy(opts.MoreHeaders, client.MoreHeaders)
	}

	var requestBody any
	if opts.RawBody != nil {
		requestBody = opts.RawBody
	} else {
		requestBody = opts.JSONBody
	}

	return client.ProviderClient.Request(
		ctx,
		method,
		url,
		requestBody,
		opts.JSONResponse,
		opts.MoreHeaders,
		opts.OKCodes,
	)
}

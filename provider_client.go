package gocloudcix

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type ProviderClient struct {
	BaseEndpoint string

	// Thread-safe token
	mut     sync.RWMutex
	TokenID string

	client     *resty.Client
	ReauthFunc func(ctx context.Context) error

	reauthmut  sync.Mutex
	MaxRetries int
}

// ErrUnexpectedResponseCode represents an unexpected HTTP response code
type ErrUnexpectedResponseCode struct {
	URL            string
	Method         string
	Expected       []int
	Actual         int
	Body           []byte
	ResponseHeader http.Header
}

func (e *ErrUnexpectedResponseCode) Error() string {
	return fmt.Sprintf("unexpected response code %d for %s %s (expected one of %v)",
		e.Actual, e.Method, e.URL, e.Expected)
}

func NewProviderClient(baseEndpoint string) *ProviderClient {
	pc := &ProviderClient{
		BaseEndpoint: baseEndpoint,
		MaxRetries:   5,
		client:       resty.New(),
	}

	// Configure resty
	pc.client.
		SetHeader("User-Agent", "gocloudcix/v0.1").
		SetRetryCount(pc.MaxRetries).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			if err != nil {
				return true
			}
			status := r.StatusCode()
			return status == 429 || status == 500 || status == 502 || status == 503 || status == 504
		})

	// Inject token before each request
	pc.client.OnBeforeRequest(func(_ *resty.Client, r *resty.Request) error {
		if token := pc.getToken(); token != "" {
			r.SetHeader("X-Auth-Token", token)
		}
		return nil
	})

	// Handle 401 by reauth + retry once
	pc.client.OnAfterResponse(func(_ *resty.Client, r *resty.Response) error {
		if r.StatusCode() != http.StatusUnauthorized {
			return nil
		}

		// Lock to prevent multiple concurrent reauths
		pc.reauthmut.Lock()
		defer pc.reauthmut.Unlock()

		oldToken := pc.getToken()

		if pc.ReauthFunc == nil {
			return errors.New("received 401 but no ReauthFunc is defined")
		}

		ctx := r.Request.Context()
		if err := pc.ReauthFunc(ctx); err != nil {
			return err
		}

		newToken := pc.getToken()
		if newToken == "" || newToken == oldToken {
			return errors.New("reauthentication failed to update token")
		}

		// Clone and retry the request
		return retryWithNewToken(pc, r)
	})

	return pc
}

func retryWithNewToken(pc *ProviderClient, origResp *resty.Response) error {
	origReq := origResp.Request

	// Copy headers (convert http.Header -> map[string]string)
	headers := map[string]string{}
	for k, vals := range origReq.Header {
		if len(vals) > 0 {
			headers[k] = vals[0]
		}
	}

	// Copy query params
	qp := map[string]string{}
	for k, vals := range origReq.QueryParam {
		if len(vals) > 0 {
			qp[k] = vals[0]
		}
	}

	// Copy body (if it's reusable)
	var bodyCopy []byte
	if origReq.Body != nil {
		switch b := origReq.Body.(type) {
		case []byte:
			bodyCopy = make([]byte, len(b))
			copy(bodyCopy, b)
		case string:
			bodyCopy = []byte(b)
		case *bytes.Buffer:
			bodyCopy = make([]byte, b.Len())
			copy(bodyCopy, b.Bytes())
		case io.Reader:
			// Read and buffer the body so it can be reused
			buf, err := io.ReadAll(b)
			if err != nil {
				return fmt.Errorf("failed to read request body: %w", err)
			}
			bodyCopy = buf
		}
	}

	// Build new request with updated token
	newReq := pc.client.R().
		SetContext(origReq.Context()).
		SetHeaders(headers).
		SetQueryParams(qp)

	// Inject updated token
	if token := pc.getToken(); token != "" {
		newReq.SetHeader("X-Auth-Token", token)
	}

	// Restore body if any
	if len(bodyCopy) > 0 {
		newReq.SetBody(bodyCopy)
	}

	// Retry request
	newResp, err := newReq.Execute(origReq.Method, origReq.URL)
	if err != nil {
		return fmt.Errorf("failed to retry request: %w", err)
	}

	// Replace original response content
	*origResp = *newResp
	return nil
}

// Thread-safe token access
func (pc *ProviderClient) getToken() string {
	pc.mut.RLock()
	defer pc.mut.RUnlock()
	return pc.TokenID
}

func (pc *ProviderClient) GetToken() string {
	pc.mut.RLock()
	defer pc.mut.RUnlock()
	return pc.TokenID
}

func (pc *ProviderClient) SetToken(token string) {
	pc.mut.Lock()
	defer pc.mut.Unlock()
	pc.TokenID = token
}

// Request performs a generic JSON API request
func (pc *ProviderClient) Request(
	ctx context.Context,
	method, url string,
	body any,
	response any,
	headers map[string]string,
	okCodes []int,
) (*resty.Response, error) {

	req := pc.client.R().SetContext(ctx)

	if body != nil {
		req.SetBody(body)
	}
	if response != nil {
		req.SetResult(response)
	}
	for k, v := range headers {
		req.SetHeader(k, v)
	}

	resp, err := req.Execute(method, url)
	if err != nil {
		return resp, fmt.Errorf("request failed: %w", err)
	}

	if len(okCodes) == 0 {
		okCodes = defaultOkCodes(method)
	}

	// Validate status code
	code := resp.StatusCode()
	for _, c := range okCodes {
		if c == code {
			return resp, nil
		}
	}

	return resp, &ErrUnexpectedResponseCode{
		URL:            url,
		Method:         method,
		Expected:       okCodes,
		Actual:         code,
		Body:           resp.Body(),
		ResponseHeader: resp.Header(),
	}
}

// Default OK codes
func defaultOkCodes(method string) []int {
	switch method {
	case http.MethodGet, http.MethodHead:
		return []int{200}
	case http.MethodPost:
		return []int{200, 201, 202}
	case http.MethodPut:
		return []int{200, 201, 202}
	case http.MethodPatch:
		return []int{200, 202, 204}
	case http.MethodDelete:
		return []int{200, 202, 204}
	default:
		return []int{}
	}
}

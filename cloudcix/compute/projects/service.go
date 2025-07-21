package project

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TVKain/go-cloudcix/cloudcix"
)

type Client struct {
	sc *cloudcix.ServiceClient
}

// New creates a new project service client
func New(sc *cloudcix.ServiceClient) *Client {
	return &Client{sc: sc}
}

// List all projects
func (c *Client) List() ([]Project, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/projects", c.sc.BaseURL), nil)
	req.Header.Set("Authorization", "Bearer "+c.sc.Token)

	resp, err := c.sc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}
	return projects, nil
}

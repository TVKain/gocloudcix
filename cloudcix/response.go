package cloudcix

// This file is part of the go-cloudcix library, which provides a client for interacting with CloudCIX APIs.
type Metadata struct {
	Limit        int      `json:"limit"`
	Page         int      `json:"page"`
	Order        string   `json:"order"`
	TotalRecords int      `json:"total_records"`
	Warnings     []string `json:"warnings"`
}

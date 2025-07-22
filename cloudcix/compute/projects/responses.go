package projects

import "github.com/TVKain/go-cloudcix/cloudcix"

type Project struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Note       string `json:"note"`
	RegionId   int    `json:"region_id"`
	AddressId  int    `json:"address_id"`
	ManagerId  int    `json:"manager_id"`
	ResellerId int    `json:"reseller_id"`
	CreatedAt  string `json:"created_at"`
	Updated    string `json:"updated_at"`
}

type ProjectsListResponse struct {
	// Projects is a list of projects returned by the API.
	Content  []Project         `json:"content"`
	Metadata cloudcix.Metadata `json:"_metadata"`
}

type ProjectGetResponse struct {
	// Projects is a list of projects returned by the API.
	Content Project `json:"content"`
}

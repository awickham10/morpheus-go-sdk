package morpheus

import (
	"fmt"
)

var (
	// BlueprintsPath is the API endpoint for blueprints
	BlueprintsPath = "/api/blueprints"
)

// Blueprint structures for use in request and response payloads
type Blueprint struct {
	ID                 int64       `json:"id"`
	Name               string      `json:"name"`
	Type               string      `json:"type"`
	Description        string      `json:"description"`
	Labels             []string    `json:"labels"`
	Category           string      `json:"category"`
	Visibility         string      `json:"visibility"`
	Config             interface{} `json:"config"`
	ResourcePermission struct {
		All   bool          `json:"all"`
		Sites []interface{} `json:"sites"`
	} `json:"resourcePermission"`
	Owner struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"owner"`
	Tenant struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"tenant"`
}

// ListBlueprintsResult structure parses the list blueprints response payload
type ListBlueprintsResult struct {
	Blueprints *[]Blueprint `json:"blueprints"`
	Meta       *MetaResult  `json:"meta"`
}

type GetBlueprintResult struct {
	Blueprint *Blueprint `json:"blueprint"`
}

type CreateBlueprintResult struct {
	Success   bool              `json:"success"`
	Message   string            `json:"msg"`
	Errors    map[string]string `json:"errors"`
	Blueprint *Blueprint        `json:"blueprint"`
}

type UpdateBlueprintResult struct {
	CreateBlueprintResult
}

type DeleteBlueprintResult struct {
	DeleteResult
}

// Client request methods

func (client *Client) ListBlueprints(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        BlueprintsPath,
		QueryParams: req.QueryParams,
		Result:      &ListBlueprintsResult{},
	})
}

func (client *Client) GetBlueprint(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", BlueprintsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetBlueprintResult{},
	})
}

func (client *Client) CreateBlueprint(req *Request) (*Response, error) {
	fmt.Println(req.Body)
	return client.Execute(&Request{
		Method:      "POST",
		Path:        BlueprintsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateBlueprintResult{},
	})
}

// UpdateBlueprint updates an existing blueprint
func (client *Client) UpdateBlueprint(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", BlueprintsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateBlueprintResult{},
	})
}

// UpdateBlueprintLogo updates an existing blueprint logo
func (client *Client) UpdateBlueprintLogo(id int64, filePayload []*FilePayload, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:         "POST",
		Path:           fmt.Sprintf("/api/blueprints/%d/image", id),
		IsMultiPart:    true,
		MultiPartFiles: filePayload,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Result: &UpdateBlueprintResult{},
	})
}

// DeleteBlueprint deletes an existing blueprint
func (client *Client) DeleteBlueprint(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", BlueprintsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteBlueprintResult{},
	})
}

func (client *Client) FindBlueprintByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListBlueprints(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListBlueprintsResult)
	blueprintsCount := len(*listResult.Blueprints)
	if blueprintsCount != 1 {
		return resp, fmt.Errorf("found %d Blueprints for %v", blueprintsCount, name)
	}
	firstRecord := (*listResult.Blueprints)[0]
	blueprintID := firstRecord.ID
	return client.GetBlueprint(blueprintID, &Request{})
}

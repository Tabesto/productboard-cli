package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tabesto/productboard-cli/internal/config"
)

// Client is the ProductBoard API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
	apiVersion string
}

// New creates a new ProductBoard API client from config.
func New(cfg *config.Config) (*Client, error) {
	if cfg.APIToken == "" {
		return nil, &APIError{
			StatusCode: 0,
			Message:    "No API token configured. Run `pboard configure` or set PRODUCTBOARD_API_TOKEN environment variable.",
			ExitCode:   ExitAuthError,
		}
	}

	return &Client{
		httpClient: &http.Client{},
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		token:      cfg.APIToken,
		apiVersion: cfg.APIVersion,
	}, nil
}

// Get performs a GET request and returns the raw response body.
func (c *Client) Get(path string, params map[string]string) ([]byte, error) {
	actualPath := path
	actualParams := params

	if c.apiVersion == "2" {
		var extraParams map[string]string
		actualPath, extraParams = translateV2Path(path)
		actualParams = translateV2Params(params)
		// Merge extra params from path translation
		if actualParams == nil {
			actualParams = make(map[string]string)
		}
		for k, v := range extraParams {
			actualParams[k] = v
		}
	}

	reqURL := c.baseURL + actualPath

	if len(actualParams) > 0 {
		q := url.Values{}
		for k, v := range actualParams {
			if v == "" {
				continue
			}
			// Support comma-separated values for array params (e.g., fields[])
			if strings.HasSuffix(k, "[]") && strings.Contains(v, ",") {
				for _, item := range strings.Split(v, ",") {
					q.Add(k, strings.TrimSpace(item))
				}
			} else {
				q.Add(k, v)
			}
		}
		if encoded := q.Encode(); encoded != "" {
			reqURL += "?" + encoded
		}
	}

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if c.apiVersion != "2" {
		req.Header.Set("X-Version", "1")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, NewNetworkError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, NewAPIError(resp.StatusCode, string(body))
	}

	return body, nil
}

// IsV2 returns true if the client is configured for API V2.
func (c *Client) IsV2() bool {
	return c.apiVersion == "2"
}

// GetSingle fetches a single resource and returns the data object.
func (c *Client) GetSingle(path string) (map[string]interface{}, error) {
	body, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if c.apiVersion == "2" {
		flattenV2Entity(resp.Data)
	}

	return resp.Data, nil
}

// GetList fetches a paginated list of resources with auto-pagination.
// It follows pageCursor until exhausted or limit is reached.
func (c *Client) GetList(path string, params map[string]string, limit int) ([]map[string]interface{}, error) {
	var allResults []map[string]interface{}

	if params == nil {
		params = make(map[string]string)
	}

	for {
		body, err := c.Get(path, params)
		if err != nil {
			return nil, err
		}

		var nextCursor string

		if c.apiVersion == "2" {
			var resp struct {
				Data  []map[string]interface{} `json:"data"`
				Links *struct {
					Next *string `json:"next"`
				} `json:"links"`
			}
			if err := json.Unmarshal(body, &resp); err != nil {
				return nil, fmt.Errorf("failed to parse response: %w", err)
			}
			for _, item := range resp.Data {
				flattenV2Entity(item)
			}
			allResults = append(allResults, resp.Data...)
			if resp.Links != nil && resp.Links.Next != nil && *resp.Links.Next != "" {
				nextCursor = extractPageCursorFromURL(*resp.Links.Next)
			}
		} else {
			var resp struct {
				Data       []map[string]interface{} `json:"data"`
				PageCursor *string                  `json:"pageCursor"`
			}
			if err := json.Unmarshal(body, &resp); err != nil {
				return nil, fmt.Errorf("failed to parse response: %w", err)
			}
			allResults = append(allResults, resp.Data...)
			if resp.PageCursor != nil && *resp.PageCursor != "" {
				nextCursor = *resp.PageCursor
			}
		}

		// Check limit
		if limit > 0 && len(allResults) >= limit {
			allResults = allResults[:limit]
			break
		}

		// Check for next page
		if nextCursor == "" {
			break
		}

		params["pageCursor"] = nextCursor
	}

	return allResults, nil
}

// GetLinkedResources fetches linked resources with auto-pagination.
// For V2, relationship responses contain source/target objects which are
// flattened to return just the target entity data.
func (c *Client) GetLinkedResources(path string, limit int) ([]map[string]interface{}, error) {
	if c.apiVersion != "2" {
		return c.GetList(path, nil, limit)
	}

	// For V2, use raw list fetch then flatten relationships
	var allResults []map[string]interface{}
	params := make(map[string]string)

	for {
		body, err := c.Get(path, params)
		if err != nil {
			return nil, err
		}

		var resp struct {
			Data  []map[string]interface{} `json:"data"`
			Links *struct {
				Next *string `json:"next"`
			} `json:"links"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		// Check if this looks like a relationship response
		v2Path, _ := translateV2Path(path)
		if isRelationshipPath(v2Path) {
			for _, item := range resp.Data {
				allResults = append(allResults, flattenV2Relationship(item))
			}
		} else {
			for _, item := range resp.Data {
				flattenV2Entity(item)
			}
			allResults = append(allResults, resp.Data...)
		}

		if limit > 0 && len(allResults) >= limit {
			allResults = allResults[:limit]
			break
		}

		var nextCursor string
		if resp.Links != nil && resp.Links.Next != nil && *resp.Links.Next != "" {
			nextCursor = extractPageCursorFromURL(*resp.Links.Next)
		}
		if nextCursor == "" {
			break
		}
		params["pageCursor"] = nextCursor
	}

	return allResults, nil
}

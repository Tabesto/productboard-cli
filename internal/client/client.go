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
	}, nil
}

// Get performs a GET request and returns the raw response body.
func (c *Client) Get(path string, params map[string]string) ([]byte, error) {
	reqURL := c.baseURL + path

	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			if v != "" {
				q.Set(k, v)
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
	req.Header.Set("X-Version", "1")

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

		var resp struct {
			Data       []map[string]interface{} `json:"data"`
			PageCursor *string                  `json:"pageCursor"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allResults = append(allResults, resp.Data...)

		// Check limit
		if limit > 0 && len(allResults) >= limit {
			allResults = allResults[:limit]
			break
		}

		// Check for next page
		if resp.PageCursor == nil || *resp.PageCursor == "" {
			break
		}

		params["pageCursor"] = *resp.PageCursor
	}

	return allResults, nil
}

// GetLinkedResources fetches linked resources with auto-pagination.
func (c *Client) GetLinkedResources(path string, limit int) ([]map[string]interface{}, error) {
	return c.GetList(path, nil, limit)
}

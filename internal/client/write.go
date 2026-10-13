package client

import (
	"encoding/json"
	"fmt"
)

// WriteRef is the minimal reference V2 returns after a create or update.
type WriteRef struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Link string `json:"link"`
}

// ParseWriteRef extracts data.id, data.type and data.links.self from a V2 write response.
func ParseWriteRef(body []byte) (*WriteRef, error) {
	var resp struct {
		Data struct {
			ID    string `json:"id"`
			Type  string `json:"type"`
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &WriteRef{ID: resp.Data.ID, Type: resp.Data.Type, Link: resp.Data.Links.Self}, nil
}

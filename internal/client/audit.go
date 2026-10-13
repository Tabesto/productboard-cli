package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// auditEntry is one JSON line of the write audit log. Each write produces a
// "request" line (written BEFORE the call, with the pre-update snapshot) and,
// when the call returns, a "response" line sharing the same op_id.
type auditEntry struct {
	Time   string `json:"time"`
	OpID   string `json:"op_id"`
	Event  string `json:"event"` // "request" or "response"
	Method string `json:"method,omitempty"`
	Path   string `json:"path,omitempty"`
	Body   any    `json:"body,omitempty"`
	// Previous is the record's `data` object as read just before an update.
	Previous json.RawMessage `json:"previous,omitempty"`
	// Revert is a ready-to-send body restoring the previous values of the fields an update touched.
	Revert any       `json:"revert,omitempty"`
	Result *WriteRef `json:"result,omitempty"`
	Error  string    `json:"error,omitempty"`
}

func (c *Client) appendAudit(e auditEntry) error {
	e.Time = time.Now().UTC().Format(time.RFC3339Nano)
	line, err := json.Marshal(e)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(c.auditPath), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(c.auditPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(line, '\n'))
	return err
}

// snapshot reads the current record before an update so it can be restored.
// It returns the raw `data` object and a revert body for the fields being replaced.
func (c *Client) snapshot(path string, body any) (json.RawMessage, any, error) {
	raw, err := c.do("GET", "/v2"+path, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	var resp struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil || len(resp.Data) == 0 {
		return nil, nil, fmt.Errorf("unexpected response while reading current value")
	}
	var current struct {
		Fields map[string]json.RawMessage `json:"fields"`
	}
	_ = json.Unmarshal(resp.Data, &current)

	// Revert body only for `fields` updates; `patch` updates keep the full snapshot.
	var revert any
	if b, ok := body.(map[string]any); ok {
		if d, ok := b["data"].(map[string]any); ok {
			if touched, ok := d["fields"].(map[string]any); ok {
				prev := map[string]json.RawMessage{}
				for k := range touched {
					if v, ok := current.Fields[k]; ok {
						prev[k] = v
					}
				}
				revert = map[string]any{"data": map[string]any{"fields": prev}}
			}
		}
	}
	return resp.Data, revert, nil
}

func newOpID() string { return strconv.FormatInt(time.Now().UnixNano(), 36) }

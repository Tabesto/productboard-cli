package client

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/tabesto/productboard-cli/internal/config"
)

func newTestClient(t *testing.T, srv *httptest.Server, version string, write bool) *Client {
	t.Helper()
	c, err := New(&config.Config{APIToken: "tok", BaseURL: srv.URL, APIVersion: version, WriteEnabled: write})
	if err != nil {
		t.Fatal(err)
	}
	c.auditPath = filepath.Join(t.TempDir(), "audit.log") // never touch the real audit log
	return c
}

func TestWriteRefusedWhenGateOff(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
	}))
	defer srv.Close()
	c := newTestClient(t, srv, "2", false)

	for name, call := range map[string]func() ([]byte, error){
		"post":  func() ([]byte, error) { return c.Post("/entities", map[string]any{"data": 1}) },
		"patch": func() ([]byte, error) { return c.Patch("/entities/1", map[string]any{"data": 1}) },
	} {
		_, err := call()
		var apiErr *APIError
		if !errors.As(err, &apiErr) || apiErr.ExitCode != ExitWriteDisabled {
			t.Errorf("%s: want exit %d, got %v", name, ExitWriteDisabled, err)
		}
	}
	if hits != 0 {
		t.Errorf("gate off: server got %d requests, want 0", hits)
	}
}

func TestWriteSentWhenGateOn(t *testing.T) {
	var method, path, auth, ctype, xver, body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		method, path, auth, ctype, xver, body = r.Method, r.URL.Path, r.Header.Get("Authorization"), r.Header.Get("Content-Type"), r.Header.Get("X-Version"), string(b)
		w.WriteHeader(201)
		_, _ = w.Write([]byte(`{"data":{"id":"abc","type":"feature","links":{"self":"https://x/v2/entities/abc"}}}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv, "2", true)

	resp, err := c.Post("/entities", map[string]any{"data": map[string]any{"type": "feature"}})
	if err != nil {
		t.Fatal(err)
	}
	if method != "POST" || path != "/v2/entities" || auth != "Bearer tok" || ctype != "application/json" || xver != "" {
		t.Errorf("unexpected request: %s %s auth=%q ctype=%q x-version=%q", method, path, auth, ctype, xver)
	}
	if body != `{"data":{"type":"feature"}}` {
		t.Errorf("body = %s", body)
	}
	ref, err := ParseWriteRef(resp)
	if err != nil || ref.ID != "abc" || ref.Type != "feature" || ref.Link != "https://x/v2/entities/abc" {
		t.Errorf("ref = %+v, err = %v", ref, err)
	}

	if _, err := c.Patch("/notes/9", map[string]any{"data": 1}); err != nil {
		t.Fatal(err)
	}
	if method != "PATCH" || path != "/v2/notes/9" {
		t.Errorf("patch: %s %s", method, path)
	}
}

func TestWriteRefusedOnV1(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("unexpected request")
	}))
	defer srv.Close()
	c := newTestClient(t, srv, "1", true)
	_, err := c.Post("/entities", map[string]any{})
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.ExitCode != ExitInvalidInput {
		t.Errorf("want exit %d, got %v", ExitInvalidInput, err)
	}
}

func TestGetWorksWithoutGate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s", r.Method)
		}
		_, _ = w.Write([]byte(`{"data":{}}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv, "2", false)
	if _, err := c.Get("/notes", nil); err != nil {
		t.Fatal(err)
	}
}

func readAudit(t *testing.T, c *Client) []map[string]any {
	t.Helper()
	data, err := os.ReadFile(c.auditPath)
	if err != nil {
		t.Fatal(err)
	}
	var out []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("bad audit line %q: %v", line, err)
		}
		out = append(out, m)
	}
	return out
}

func TestUpdateIsAuditedWithRevertBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			_, _ = w.Write([]byte(`{"data":{"id":"abc","type":"feature","fields":{"name":"Old","status":"x"}}}`))
			return
		}
		_, _ = w.Write([]byte(`{"data":{"id":"abc","type":"feature","links":{"self":"u"}}}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv, "2", true)

	body := map[string]any{"data": map[string]any{"fields": map[string]any{"name": "New", "owner": "z"}}}
	if _, err := c.Patch("/entities/abc", body); err != nil {
		t.Fatal(err)
	}
	log := readAudit(t, c)
	if len(log) != 2 || log[0]["event"] != "request" || log[1]["event"] != "response" || log[0]["op_id"] != log[1]["op_id"] {
		t.Fatalf("log = %v", log)
	}
	if log[0]["method"] != "PATCH" || log[0]["path"] != "/v2/entities/abc" {
		t.Errorf("request line = %v", log[0])
	}
	// revert restores only touched fields that existed before ("owner" did not)
	rev, _ := json.Marshal(log[0]["revert"])
	if string(rev) != `{"data":{"fields":{"name":"Old"}}}` {
		t.Errorf("revert = %s", rev)
	}
	if prev, _ := json.Marshal(log[0]["previous"]); !strings.Contains(string(prev), `"status":"x"`) {
		t.Errorf("previous snapshot missing: %s", prev)
	}
	if res, _ := log[1]["result"].(map[string]any); res["id"] != "abc" {
		t.Errorf("response line = %v", log[1])
	}
	if fi, _ := os.Stat(c.auditPath); fi.Mode().Perm() != 0600 {
		t.Errorf("audit mode = %v", fi.Mode().Perm())
	}
}

func TestCreateIsAuditedWithoutSnapshot(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("unexpected %s", r.Method)
		}
		_, _ = w.Write([]byte(`{"data":{"id":"n1","type":"textNote","links":{"self":"u"}}}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv, "2", true)
	if _, err := c.Post("/notes", map[string]any{"data": map[string]any{"type": "textNote"}}); err != nil {
		t.Fatal(err)
	}
	log := readAudit(t, c)
	if len(log) != 2 || log[0]["previous"] != nil {
		t.Fatalf("log = %v", log)
	}
	if res, _ := log[1]["result"].(map[string]any); res["id"] != "n1" {
		t.Errorf("response line = %v", log[1])
	}
}

func TestUpdateAbortsWhenSnapshotFails(t *testing.T) {
	var patches int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			atomic.AddInt32(&patches, 1)
		}
		w.WriteHeader(404)
	}))
	defer srv.Close()
	c := newTestClient(t, srv, "2", true)
	if _, err := c.Patch("/entities/missing", map[string]any{"data": map[string]any{}}); err == nil {
		t.Fatal("want error")
	}
	if patches != 0 {
		t.Error("PATCH sent although the snapshot failed")
	}
}

func TestWriteAbortsWhenAuditLogUnwritable(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
	}))
	defer srv.Close()
	c := newTestClient(t, srv, "2", true)
	blocker := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(blocker, nil, 0600); err != nil {
		t.Fatal(err)
	}
	c.auditPath = filepath.Join(blocker, "audit.log") // parent is a regular file
	if _, err := c.Post("/entities", map[string]any{"data": 1}); err == nil {
		t.Fatal("want error")
	}
	if hits != 0 {
		t.Errorf("request sent although audit log was unwritable (%d hits)", hits)
	}
}

func TestFeatureStatusesV2(t *testing.T) {
	var path string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		_, _ = w.Write([]byte(`{"data":{"type":"feature","fields":{"status":{"id":"s","values":{"data":[
			{"id":"1","name":"Backlog","assignedEntityTypes":["feature"]},
			{"id":"2","name":"Beta","assignedEntityTypes":["feature"]},
			{"id":"3","name":"Done","assignedEntityTypes":["feature"]}]}}}}}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv, "2", false)

	got, err := c.GetList("/feature-statuses", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	if path != "/v2/entities/configurations/feature" {
		t.Errorf("path = %s", path)
	}
	if len(got) != 3 || got[1]["id"] != "2" || got[1]["name"] != "Beta" {
		t.Errorf("got %v", got)
	}
	if got, _ = c.GetList("/feature-statuses", nil, 2); len(got) != 2 {
		t.Errorf("limit 2: got %d", len(got))
	}
}

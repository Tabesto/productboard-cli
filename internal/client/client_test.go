package client

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
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

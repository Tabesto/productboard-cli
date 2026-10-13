package cli

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseFieldPairs(t *testing.T) {
	got, err := parseFieldPairs([]string{"name=Hello world", "n=3", "ok=true", `o={"a":1}`, "eq=a=b"})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"name": "Hello world", "n": float64(3), "ok": true, "o": map[string]any{"a": float64(1)}, "eq": "a=b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v", got)
	}
	for _, bad := range []string{"noequals", "=x"} {
		if _, err := parseFieldPairs([]string{bad}); err == nil {
			t.Errorf("%q: want error", bad)
		}
	}
}

func marshal(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestBuildUpdateBody(t *testing.T) {
	body, err := buildUpdateBody([]string{"name=x"}, "", "")
	if err != nil || marshal(t, body) != `{"data":{"fields":{"name":"x"}}}` {
		t.Errorf("fields: %v %v", body, err)
	}
	body, err = buildUpdateBody(nil, `[{"op":"clear","path":"owner"}]`, "")
	if err != nil || marshal(t, body) != `{"data":{"patch":[{"op":"clear","path":"owner"}]}}` {
		t.Errorf("patch: %v %v", body, err)
	}
	body, err = buildUpdateBody(nil, "", `{"fields":{"a":1}}`)
	if err != nil || marshal(t, body) != `{"data":{"fields":{"a":1}}}` {
		t.Errorf("json: %v %v", body, err)
	}
	for name, args := range map[string][3]string{
		"none":       {"", "", ""},
		"bad patch":  {"", `{"op":1}`, ""},
		"bad json":   {"", "", `[1]`},
		"patch+json": {"", `[]`, `{}`},
	} {
		if _, err := buildUpdateBody(nil, args[1], args[2]); err == nil {
			t.Errorf("%s: want error", name)
		}
	}
	if _, err := buildUpdateBody([]string{"a=1"}, `[{"op":"clear"}]`, ""); err == nil {
		t.Error("fields+patch: want error")
	}
}

func TestBuildEntityCreateBody(t *testing.T) {
	body, err := buildEntityCreateBody("feature", []string{"name=X"}, "")
	if err != nil || marshal(t, body) != `{"data":{"fields":{"name":"X"},"type":"feature"}}` {
		t.Errorf("%v %v", body, err)
	}
	if _, err := buildEntityCreateBody("", nil, ""); err == nil {
		t.Error("missing type: want error")
	}
	if _, err := buildEntityCreateBody("", nil, `{"type":"feature"}`); err != nil {
		t.Errorf("type in json: %v", err)
	}
	body, err = buildEntityCreateBody("", []string{"b=2"}, `{"type":"feature","fields":{"a":1}}`)
	if err != nil || marshal(t, body) != `{"data":{"fields":{"a":1,"b":2},"type":"feature"}}` {
		t.Errorf("json+field merge: %s %v", marshal(t, body), err)
	}
	if _, err := buildEntityCreateBody("", nil, `{"type":"feature","fields":[1]}`); err == nil {
		t.Error("non-object fields: want error")
	}
}

func TestBuildNoteCreateBody(t *testing.T) {
	body, err := buildNoteCreateBody("textNote", "N", "C", []string{"owner={\"email\":\"a@b.c\"}"}, "")
	want := `{"data":{"fields":{"content":"C","name":"N","owner":{"email":"a@b.c"}},"type":"textNote"}}`
	if err != nil || marshal(t, body) != want {
		t.Errorf("%s %v", marshal(t, body), err)
	}
	if _, err := buildNoteCreateBody("textNote", "", "C", nil, ""); err == nil {
		t.Error("missing name: want error")
	}
	if _, err := buildNoteCreateBody("textNote", "", "", []string{"name=viaField"}, ""); err != nil {
		t.Errorf("name via --field: %v", err)
	}

	// raw --json must enforce the name and honour --name/--content/--field
	for _, raw := range []string{`{"type":"textNote","fields":{}}`, `{"type":"textNote"}`, `{"fields":{"name":""}}`} {
		if _, err := buildNoteCreateBody("textNote", "", "", nil, raw); err == nil {
			t.Errorf("raw %s without name: want error", raw)
		}
	}
	body, err = buildNoteCreateBody("textNote", "N", "C", []string{"x=1"}, `{"fields":{"x":0,"y":2}}`)
	want = `{"data":{"fields":{"content":"C","name":"N","x":1,"y":2},"type":"textNote"}}`
	if err != nil || marshal(t, body) != want {
		t.Errorf("raw merge: %s %v", marshal(t, body), err)
	}
	if _, err := buildNoteCreateBody("textNote", "", "", nil, `{"fields":{"name":"ok"}}`); err != nil {
		t.Errorf("raw with name: %v", err)
	}
}

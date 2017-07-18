package vendorfmt

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestFmt(t *testing.T) {
	decode := func(s string) map[string]interface{} {
		var m map[string]interface{}
		if err := json.NewDecoder(strings.NewReader(s)).Decode(&m); err != nil {
			t.Fatalf("got error %s want nil")
		}
		return m
	}

	in := `{
	"comment": "",
	"ignore": "test",
	"package": [
		{
			"path": "appengine",
			"revision": ""
		},
		{
			"path": "appengine_internal",
			"revision": ""
		},
		{
			"path": "appengine_internal/base",
			"revision": ""
		}
	]
}`
	want := `{
	"comment": "",
	"ignore": "test",
	"package": [
		{"path":"appengine","revision":""},
		{"path":"appengine_internal","revision":""},
		{"path":"appengine_internal/base","revision":""}
	]
}`

	got, err := FormatString(in)
	if err != nil {
		t.Fatalf("got %v want nil", err)
	}

	// verify content is equal
	if got, want := decode(got), decode(in); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}

	// verify format is ok
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}

}

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
			t.Fatalf("got error %s want nil", err)
		}
		return m
	}

	in := `{
	"comment": "",
	"ignore": "test",
	"package": [
		{
			"checksum": "a",
			"path": "x",
			"revision": true
		},
		{
			"checksum": "b",
			"path": "y",
			"revision": false
		},
		{
			"checksum": "c",
			"path": "z/x",
			"revision": true
		}
	]
}`
	want := `{
	"comment": "",
	"ignore": "test",
	"package": [
		{"path":"x","checksum":"a","revision":true},
		{"path":"y","checksum":"b","revision":false},
		{"path":"z/x","checksum":"c","revision":true}
	]
}
`

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

func TestPkgSort(t *testing.T) {
	in := `{
	"comment": "",
	"ignore": "test",
	"package": [
		{
			"checksum": "b",
			"path": "y",
			"revision":false
		},
		{
			"checksum": "a",
			"path": "x",
			"revision":true
		},
		{
			"checksum": "c",
			"path": "z/x",
			"revision":true
		}
	]
}`
	want := `{
	"comment": "",
	"ignore": "test",
	"package": [
		{"path":"x","checksum":"a","revision":true},
		{"path":"y","checksum":"b","revision":false},
		{"path":"z/x","checksum":"c","revision":true}
	]
}
`

	got, err := FormatString(in)
	if err != nil {
		t.Fatalf("got %v want nil", err)
	}

	// verify packages are sorted
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}

	// verify json is still legal
	var v interface{}
	if err := json.Unmarshal([]byte(got), &v); err != nil {
		t.Fatalf("json.Unmarshal failed: %s", err)
	}
}

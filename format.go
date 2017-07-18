package vendorfmt

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func FormatString(s string) (string, error) {
	b, err := Format([]byte(s))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Format renders a vendor.json file so that all entries of the
// 'package' element are written as single-line entries which
// simplifies merging.
func Format(b []byte) ([]byte, error) {
	return FormatIndent(b, "", "\t")
}

func FormatIndent(b []byte, prefix, indent string) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	// fail fast if the structure doesn't
	// have a non-empty 'package' array
	if len(m) == 0 || m["package"] == nil {
		return b, nil
	}
	p, ok := m["package"].([]interface{})
	if !ok || len(p) == 0 {
		return b, nil
	}

	// render the input w/o 'package'
	m["package"] = []interface{}{}
	out, err := json.MarshalIndent(m, prefix, indent)
	if err != nil {
		return nil, fmt.Errorf("cannot render to JSON: %s", err)
	}

	// render package with one line per entry
	var pb bytes.Buffer
	pb.WriteString("\"package\": [\n")
	for i, x := range p {
		b, _ := json.Marshal(x)
		pb.WriteString(indent)
		pb.WriteString(indent)
		pb.Write(b)
		if i < len(p)-1 {
			pb.WriteString(",")
		}
		pb.WriteString("\n")
	}
	pb.WriteString(indent)
	pb.WriteString("]")

	// replace "package": [] with new content
	return bytes.Replace(out, []byte(`"package": []`), pb.Bytes(), 1), nil
}

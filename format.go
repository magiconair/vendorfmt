package vendorfmt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
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

	// render packages with one line per entry
	var pb bytes.Buffer
	var pkgs []string
	for _, x := range p {
		// sort map keys and make "path" first element
		pm := x.(map[string]interface{})
		var keys []string
		for k := range pm {
			if k == "path" {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		keys = append([]string{"path"}, keys...)

		// render package entry
		pb.Reset()
		pb.WriteString(indent)
		pb.WriteString(indent)
		pb.WriteString(`{`)
		for j, k := range keys {
			v := pm[k]
			pb.WriteString(`"`)
			pb.WriteString(k)
			pb.WriteString(`":`)
			switch vv := v.(type) {
			case string:
				pb.WriteRune('"')
				pb.WriteString(vv)
				pb.WriteRune('"')
			case bool:
				if vv {
					pb.WriteString("true")
				} else {
					pb.WriteString("false")
				}
			default:
				panic("unknown type")
			}
			if j < len(keys)-1 {
				pb.WriteString(`,`)
			}
		}
		pb.WriteString("}")
		pkgs = append(pkgs, pb.String())
	}

	// sort package entries by path
	sort.Strings(pkgs)

	// render "package" array
	pb.Reset()
	pb.WriteString("\"package\": [\n")
	pb.WriteString(strings.Join(pkgs, ",\n"))
	pb.WriteString("\n")
	pb.WriteString(indent)
	pb.WriteString("]")

	// replace "package": [] with new content
	out = bytes.Replace(out, []byte(`"package": []`), pb.Bytes(), 1)

	// ensure that file ends with new line
	if !bytes.HasSuffix(out, []byte{'\n'}) {
		out = append(out, '\n')
	}

	return out, nil
}

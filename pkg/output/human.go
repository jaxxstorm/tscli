package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

type HumanPrinter struct{}

func (HumanPrinter) Print(raw []byte) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()

	var data any
	if err := dec.Decode(&data); err != nil {
		return err
	}

	switch v := data.(type) {
	case []any:
		return printArray(v)
	case map[string]any:
		return printMap(v)
	default:
		var obj any
		_ = json.Unmarshal(raw, &obj)
		s, _ := json.Marshal(obj)
		_, _ = os.Stdout.Write(s)
		return nil
	}
}

func printArray(arr []any) error {
	for i, itm := range arr {
		if m, ok := itm.(map[string]any); ok {
			if err := printMap(m); err != nil {
				return err
			}
		} else {
			// not an object → dump raw JSON element
			line, _ := json.Marshal(itm)
			fmt.Println(trunc(string(line)))
		}

		if i != len(arr)-1 {
			fmt.Println(strings.Repeat("─", 64))
		}
	}
	return nil
}

func printMap(m map[string]any) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	var keys []string
	for k := range maps.Keys(m) {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		val := fmtVal(m[k])
		fmt.Fprintf(w, "%s:\t%s\n", bold(k), val)
	}
	return w.Flush()
}

func fmtVal(v any) string {
	switch x := v.(type) {
	case nil:
		return "null"
	case string:
		if len(x) > 80 {
			return trunc(x)
		}
		return x
	case json.Number:
		return x.String()
	case bool:
		return fmt.Sprintf("%v", x)
	case []any:
		// join scalars, otherwise indicate array length
		var parts []string
		for _, el := range x {
			if s, ok := inlineVal(el); ok {
				parts = append(parts, s)
			} else {
				return fmt.Sprintf("[%d items]", len(x))
			}
		}
		return strings.Join(parts, "  ")
	case map[string]any:
		if s, ok := inlineMap(x); ok {
			return s
		}
		b, _ := json.Marshal(x)
		return trunc(string(b))
	default:
		b, _ := json.Marshal(x)
		return trunc(string(b))
	}
}

func inlineVal(v any) (string, bool) {
	switch x := v.(type) {
	case nil:
		return "null", true
	case string:
		return trunc(x), true
	case json.Number:
		return x.String(), true
	case bool:
		return fmt.Sprintf("%v", x), true
	case []any:
		var parts []string
		for _, el := range x {
			part, ok := inlineVal(el)
			if !ok {
				return "", false
			}
			parts = append(parts, part)
		}
		return "[" + strings.Join(parts, ", ") + "]", true
	case map[string]any:
		return inlineMap(x)
	default:
		return "", false
	}
}

func inlineMap(m map[string]any) (string, bool) {
	if len(m) == 0 {
		return "{}", true
	}

	var keys []string
	for k := range maps.Keys(m) {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	totalLen := 2 // {}
	for _, k := range keys {
		part, ok := inlineVal(m[k])
		if !ok {
			return "", false
		}
		entry := k + ": " + part
		totalLen += len(entry) + 2
		parts = append(parts, entry)
	}

	return trunc("{" + strings.Join(parts, ", ") + "}"), true
}

func trunc(s string) string {
	if len(s) > 80 {
		return s[:77] + "…"
	}
	return s
}

func bold(s string) string { return "\033[1m" + s + "\033[0m" }

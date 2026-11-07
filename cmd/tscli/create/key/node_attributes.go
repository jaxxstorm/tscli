package key

import (
	"encoding/json"
	"fmt"
	"strings"
)

func parseNodeAttributes(input string) (map[string]json.RawMessage, error) {
	var attributes map[string]json.RawMessage
	if err := json.Unmarshal([]byte(input), &attributes); err != nil {
		return nil, fmt.Errorf("--node-attributes must be a JSON object: %w", err)
	}
	if attributes == nil {
		return nil, fmt.Errorf("--node-attributes must be a JSON object")
	}
	if len(attributes) > 100 {
		return nil, fmt.Errorf("--node-attributes allows at most 100 attributes")
	}
	for name, raw := range attributes {
		if !strings.HasPrefix(name, "custom:") || name == "custom:" {
			return nil, fmt.Errorf("--node-attributes keys must have a nonempty name in the custom: namespace")
		}
		// Keep the raw JSON so large integers are never rounded through float64.
		decoder := json.NewDecoder(strings.NewReader(string(raw)))
		decoder.UseNumber()
		var value any
		if err := decoder.Decode(&value); err != nil {
			return nil, fmt.Errorf("--node-attributes: %w", err)
		}
		switch v := value.(type) {
		case string, bool:
			continue
		case json.Number:
			if !strings.ContainsAny(string(v), ".eE") {
				continue
			}
		}
		return nil, fmt.Errorf("--node-attributes values must be a string, integer, or boolean (use integer literals for numbers)")
	}
	return attributes, nil
}

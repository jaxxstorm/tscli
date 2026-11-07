package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type YAMLPrinter struct{}

// yamlNumbers retains JSON numeric values and types, including large integers.
func yamlNumbers(value any) any {
	switch v := value.(type) {
	case json.Number:
		tag := "!!int"
		if strings.ContainsAny(string(v), ".eE") {
			tag = "!!float"
		}
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: tag, Value: string(v)}
	case map[string]any:
		for k, child := range v {
			v[k] = yamlNumbers(child)
		}
	case []any:
		for i, child := range v {
			v[i] = yamlNumbers(child)
		}
	}
	return value
}

func (YAMLPrinter) Print(b []byte) error {
	if !json.Valid(b) {
		return fmt.Errorf("invalid JSON for YAML output")
	}
	// Decode JSON scalars as raw values, so integers never pass through float64.
	var obj any
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	if err := dec.Decode(&obj); err != nil {
		return err
	}
	y, err := yaml.Marshal(yamlNumbers(obj))
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(y)
	return err
}

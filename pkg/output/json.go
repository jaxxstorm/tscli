// pkg/output/json.go
package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type JSONPrinter struct{}

func (JSONPrinter) Print(b []byte) error {
	var obj any
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	if err := dec.Decode(&obj); err != nil || !json.Valid(b) {
		_, err = os.Stdout.Write(b)
		return err
	}

	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Println(string(out))
	return err
}

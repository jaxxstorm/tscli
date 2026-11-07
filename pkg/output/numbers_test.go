package output

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintersPreserveNumbers(t *testing.T) {
	for _, format := range []string{"json", "yaml"} {
		t.Run(format, func(t *testing.T) {
			original := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()
			os.Stdout = w
			defer func() { os.Stdout = original }()
			err = Print(format, []byte(`{"attributes":[9007199254740993,123456789012345678901234567890,-12,1.25,1e3],"text":"9007199254740993","enabled":false}`))
			w.Close()
			os.Stdout = original
			if err != nil {
				t.Fatal(err)
			}
			out, err := io.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			for _, number := range []string{"9007199254740993", "123456789012345678901234567890", "-12", "1.25", "1e3"} {
				if !strings.Contains(string(out), number) {
					t.Errorf("missing exact number %s in %s", number, out)
				}
			}
		})
	}
}

func TestYAMLPrinterRejectsTrailingJSON(t *testing.T) {
	if err := (YAMLPrinter{}).Print([]byte(`{} {}`)); err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

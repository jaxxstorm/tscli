package output

import "slices"

// Printer prints a marshalled JSON payload to the terminal in some format.
type Printer interface {
	Print(jsonBytes []byte) error
}

// Registry of available printers.
var registry = map[string]Printer{
	"json":   JSONPrinter{},
	"yaml":   YAMLPrinter{},
	"human":  HumanPrinter{},
	"pretty": PrettyPrinter{},
}

// Get returns the requested printer or falls back to plain JSON.
func Get(format string) Printer {
	if p, ok := registry[format]; ok {
		return p
	}
	return registry["json"]
}

// Register lets you plug-in new printers from init() funcs.
func Register(name string, p Printer) {
	if name == "" {
		panic("out.Register: empty name")
	}
	if _, dup := registry[name]; dup {
		panic("out.Register: duplicate name " + name)
	}
	registry[name] = p
}

// Convenience helper used by commands.
func Print(format string, jsonBytes []byte) error {
	return Get(format).Print(jsonBytes)
}

// user-facing list, e.g. for --help
func Available() (out []string) {
	for k := range registry {
		out = append(out, k)
	}
	slices.Sort(out)
	return
}

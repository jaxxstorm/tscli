package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/jaxxstorm/tscli/internal/cli"
	"github.com/jaxxstorm/tscli/pkg/contract"
	pkgversion "github.com/jaxxstorm/tscli/pkg/version"
)

func main() {
	if err := fang.Execute(context.Background(), cli.Configure(), fang.WithVersion(pkgversion.GetVersion())); err != nil {
		contract.IgnoreIoError(fmt.Fprintf(os.Stderr, "%v\n", err))
		os.Exit(1)
	}
}

package cli_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/jaxxstorm/tscli/internal/cli"
	"github.com/spf13/viper"
)

type execResult struct {
	stdout string
	stderr string
	err    error
}

var stdioMu sync.Mutex

func executeCLI(t *testing.T, args []string, env map[string]string) execResult {
	t.Helper()
	return executeCLIWithInputAndDefaults(t, args, env, "", true)
}

func executeCLINoDefaults(t *testing.T, args []string, env map[string]string) execResult {
	t.Helper()
	return executeCLIWithInputAndDefaults(t, args, env, "", false)
}

func executeCLIWithInput(t *testing.T, args []string, env map[string]string, input string) execResult {
	t.Helper()
	return executeCLIWithInputAndDefaults(t, args, env, input, true)
}

func executeCLINoDefaultsWithInput(t *testing.T, args []string, env map[string]string, input string) execResult {
	t.Helper()
	return executeCLIWithInputAndDefaults(t, args, env, input, false)
}

func executeCLIWithInputAndDefaults(t *testing.T, args []string, env map[string]string, input string, useDefaults bool) execResult {
	t.Helper()

	viper.Reset()

	if _, ok := env["HOME"]; !ok {
		home := t.TempDir()
		t.Setenv("HOME", home)
		cfgPath := filepath.Join(home, ".tscli.yaml")
		_ = os.WriteFile(cfgPath, []byte("output: json\n"), 0o600)
	}

	if useDefaults {
		t.Setenv("TAILSCALE_API_KEY", "tskey-test")
		t.Setenv("TAILSCALE_TAILNET", "example.com")
		t.Setenv("TSCLI_OUTPUT", "json")
	}

	for k, v := range env {
		t.Setenv(k, v)
	}

	cmd := cli.Configure()
	cmd.SetArgs(args)
	cmd.SetIn(strings.NewReader(input))

	stdout, stderr, err := captureStdIO(func() error {
		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stderr)
		return cmd.ExecuteContext(context.Background())
	})

	return execResult{
		stdout: strings.TrimSpace(stdout),
		stderr: strings.TrimSpace(stderr),
		err:    err,
	}
}

func captureStdIO(run func() error) (string, string, error) {
	stdioMu.Lock()
	defer stdioMu.Unlock()

	origOut := os.Stdout
	origErr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	doneOut := make(chan struct{})
	doneErr := make(chan struct{})

	go func() {
		_, _ = io.Copy(&outBuf, rOut)
		close(doneOut)
	}()
	go func() {
		_, _ = io.Copy(&errBuf, rErr)
		close(doneErr)
	}()

	runErr := run()

	_ = wOut.Close()
	_ = wErr.Close()
	<-doneOut
	<-doneErr

	os.Stdout = origOut
	os.Stderr = origErr
	_ = rOut.Close()
	_ = rErr.Close()

	return outBuf.String(), errBuf.String(), runErr
}

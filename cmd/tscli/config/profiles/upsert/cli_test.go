package upsert

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestPromptForValueUsesInlinePrompt(t *testing.T) {
	cmd := &cobra.Command{}
	input := strings.NewReader("cid\n")
	output := &bytes.Buffer{}
	cmd.SetIn(input)
	cmd.SetOut(output)

	reader := bufio.NewReader(input)
	var value string
	if err := promptForValue(cmd, reader, "OAuth client ID: ", &value); err != nil {
		t.Fatalf("promptForValue: %v", err)
	}
	if value != "cid" {
		t.Fatalf("expected value %q, got %q", "cid", value)
	}
	if output.String() != "OAuth client ID: " {
		t.Fatalf("expected inline prompt output, got %q", output.String())
	}
}

func TestPromptForSecretValueFallsBackToReaderWhenNotTerminal(t *testing.T) {
	cmd := &cobra.Command{}
	input := strings.NewReader("secret\n")
	output := &bytes.Buffer{}
	cmd.SetIn(input)
	cmd.SetOut(output)

	reader := bufio.NewReader(input)
	var value string
	if err := promptForSecretValue(cmd, reader, "OAuth client secret: ", &value); err != nil {
		t.Fatalf("promptForSecretValue: %v", err)
	}
	if value != "secret" {
		t.Fatalf("expected value %q, got %q", "secret", value)
	}
	if output.String() != "OAuth client secret: " {
		t.Fatalf("expected inline secret prompt output, got %q", output.String())
	}
}

func TestPromptForProfileAuthPromptsInlineForOAuthFlow(t *testing.T) {
	cmd := &cobra.Command{}
	input := strings.NewReader("oauth\ncid\nsecret\n")
	output := &bytes.Buffer{}
	cmd.SetIn(input)
	cmd.SetOut(output)

	reader := bufio.NewReader(input)
	var apiKey, oauthClientID, oauthClientSecret string
	if err := promptForProfileAuth(cmd, reader, &apiKey, &oauthClientID, &oauthClientSecret); err != nil {
		t.Fatalf("promptForProfileAuth: %v", err)
	}
	if apiKey != "" {
		t.Fatalf("expected api key to remain empty, got %q", apiKey)
	}
	if oauthClientID != "cid" || oauthClientSecret != "secret" {
		t.Fatalf("expected oauth credentials to be captured, got id=%q secret=%q", oauthClientID, oauthClientSecret)
	}
	if output.String() != "Auth type [api-key|oauth]: OAuth client ID: OAuth client secret: " {
		t.Fatalf("expected inline prompt sequence, got %q", output.String())
	}
}

package oauthapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

func CreateCommand() *cobra.Command {
	var name, description string
	var redirectURIs, scopes, allowedNodeAttributes []string
	cmd := &cobra.Command{Use: "oauth-app", Short: "Create an OAuth app", PreRunE: func(_ *cobra.Command, _ []string) error {
		return validateAppInput(name, description, redirectURIs, scopes, allowedNodeAttributes)
	}, RunE: func(_ *cobra.Command, _ []string) error {
		client, err := tscli.New()
		if err != nil {
			return err
		}
		return doJSON(client, http.MethodPost, "/tailnet/{tailnet}/oauth-apps", appPayload(name, description, redirectURIs, scopes, allowedNodeAttributes), "create OAuth app")
	}}
	appFlags(cmd, &name, &description, &redirectURIs, &scopes, &allowedNodeAttributes)
	return cmd
}

func GetCommand() *cobra.Command {
	var appID string
	cmd := &cobra.Command{Use: "oauth-app", Short: "Get an OAuth app", RunE: func(_ *cobra.Command, _ []string) error {
		client, err := tscli.New()
		if err != nil {
			return err
		}
		return doJSON(client, http.MethodGet, appPath(appID), nil, "get OAuth app")
	}}
	appIDFlag(cmd, &appID)
	return cmd
}

func ListCommand() *cobra.Command {
	return &cobra.Command{Use: "oauth-apps", Short: "List OAuth apps", RunE: func(_ *cobra.Command, _ []string) error {
		client, err := tscli.New()
		if err != nil {
			return err
		}
		return doJSON(client, http.MethodGet, "/tailnet/{tailnet}/oauth-apps", nil, "list OAuth apps")
	}}
}

func SetCommand() *cobra.Command {
	var appID, name, description string
	var redirectURIs, scopes, allowedNodeAttributes []string
	cmd := &cobra.Command{Use: "oauth-app", Short: "Update an OAuth app", PreRunE: func(_ *cobra.Command, _ []string) error {
		return validateAppInput(name, description, redirectURIs, scopes, allowedNodeAttributes)
	}, RunE: func(_ *cobra.Command, _ []string) error {
		client, err := tscli.New()
		if err != nil {
			return err
		}
		return doJSON(client, http.MethodPut, appPath(appID), appPayload(name, description, redirectURIs, scopes, allowedNodeAttributes), "update OAuth app")
	}}
	appIDFlag(cmd, &appID)
	appFlags(cmd, &name, &description, &redirectURIs, &scopes, &allowedNodeAttributes)
	return cmd
}

func DeleteCommand() *cobra.Command {
	var appID string
	cmd := &cobra.Command{Use: "oauth-app", Short: "Delete an OAuth app", RunE: func(_ *cobra.Command, _ []string) error {
		client, err := tscli.New()
		if err != nil {
			return err
		}
		if _, err := tscli.Do(context.Background(), client, http.MethodDelete, appPath(appID), nil, nil); err != nil {
			return fmt.Errorf("failed to delete OAuth app: %w", err)
		}
		out, _ := json.MarshalIndent(map[string]string{"result": fmt.Sprintf("OAuth app %s deleted", appID)}, "", "  ")
		return output.Print(viper.GetString("output"), out)
	}}
	appIDFlag(cmd, &appID)
	return cmd
}

func appIDFlag(cmd *cobra.Command, appID *string) {
	cmd.Flags().StringVar(appID, "id", "", "OAuth app ID")
	_ = cmd.MarkFlagRequired("id")
}

func appFlags(cmd *cobra.Command, name, description *string, redirectURIs, scopes, allowedNodeAttributes *[]string) {
	cmd.Flags().StringVar(name, "name", "", "OAuth app name")
	cmd.Flags().StringVar(description, "description", "", "OAuth app description")
	cmd.Flags().StringSliceVar(redirectURIs, "redirect-uri", nil, "Permitted redirect URI (repeatable)")
	cmd.Flags().StringSliceVar(scopes, "scope", nil, "OAuth scope (repeatable)")
	cmd.Flags().StringSliceVar(allowedNodeAttributes, "allowed-node-attribute", nil, "Custom node attribute the app may set (repeatable)")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("redirect-uri")
	_ = cmd.MarkFlagRequired("scope")
}

func appPayload(name, description string, redirectURIs, scopes, allowedNodeAttributes []string) map[string]any {
	payload := map[string]any{"name": name, "redirectURIs": redirectURIs, "scopes": scopes}
	if description != "" {
		payload["description"] = description
	}
	if len(allowedNodeAttributes) > 0 {
		payload["allowedNodeAttributes"] = allowedNodeAttributes
	}
	return payload
}

func validateAppInput(name, description string, redirectURIs, scopes, allowedNodeAttributes []string) error {
	if len(name) < 3 || len(name) > 50 {
		return fmt.Errorf("--name must be between 3 and 50 characters")
	}
	if len(description) > 300 {
		return fmt.Errorf("--description must be at most 300 characters")
	}
	for _, redirectURI := range redirectURIs {
		u, err := url.ParseRequestURI(redirectURI)
		if err != nil || u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("invalid --redirect-uri %q", redirectURI)
		}
	}
	for _, scope := range scopes {
		if strings.TrimSpace(scope) == "" {
			return fmt.Errorf("--scope cannot be empty")
		}
	}
	for _, attribute := range allowedNodeAttributes {
		if strings.TrimSpace(attribute) == "" {
			return fmt.Errorf("--allowed-node-attribute cannot be empty")
		}
	}
	return nil
}

func appPath(appID string) string {
	return fmt.Sprintf("/tailnet/{tailnet}/oauth-apps/%s", url.PathEscape(appID))
}

func doJSON(client *tsapi.Client, method, path string, payload any, action string) error {
	var response json.RawMessage
	if _, err := tscli.Do(context.Background(), client, method, path, payload, &response); err != nil {
		return fmt.Errorf("failed to %s: %w", action, err)
	}
	out, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}
	return output.Print(viper.GetString("output"), out)
}

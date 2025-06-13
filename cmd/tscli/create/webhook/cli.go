// cmd/tscli/create/webhook/cli.go
//
// Create a webhook (generic provider or Slack, etc.).
//
// Example
//
//	tscli create webhook \
//	      --url https://hooks.example.com/tailscale \
//	      --provider generic \
//	      --subscription device_authorized --subscription device_deleted
package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/jaxxstorm/tscli/pkg/output"

	"github.com/jaxxstorm/tscli/pkg/tscli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tsapi "tailscale.com/client/tailscale/v2"
)

var providerEnum = map[string]tsapi.WebhookProviderType{
	"generic":    tsapi.WebhookProviderType("generic"),
	"slack":      tsapi.WebhookProviderType("slack"),
	"mattermost": tsapi.WebhookProviderType("mattermost"),
	"googlechat": tsapi.WebhookProviderType("googlechat"),
	"discord":    tsapi.WebhookProviderType("discord"),
}

var subscriptionEnum = map[string]tsapi.WebhookSubscriptionType{
	"nodeCreated":             tsapi.WebhookSubscriptionType("nodeCreated"),
	"nodeNeedsApproval":       tsapi.WebhookSubscriptionType("nodeNeedsApproval"),
	"nodeApproved":            tsapi.WebhookSubscriptionType("nodeApproved"),
	"nodeKeyExpiringInOneDay": tsapi.WebhookSubscriptionType("nodeKeyExpiringInOneDay"),
	"nodeKeyExpired":          tsapi.WebhookSubscriptionType("nodeKeyExpired"),
	"nodeDeleted":             tsapi.WebhookSubscriptionType("nodeDeleted"),

	"policyUpdate": tsapi.WebhookSubscriptionType("policyUpdate"),

	"userCreated":       tsapi.WebhookSubscriptionType("userCreated"),
	"userNeedsApproval": tsapi.WebhookSubscriptionType("userNeedsApproval"),
	"userSuspended":     tsapi.WebhookSubscriptionType("userSuspended"),
	"userRestored":      tsapi.WebhookSubscriptionType("userRestored"),
	"userDeleted":       tsapi.WebhookSubscriptionType("userDeleted"),
	"userApproved":      tsapi.WebhookSubscriptionType("userApproved"),
	"userRoleUpdated":   tsapi.WebhookSubscriptionType("userRoleUpdated"),

	"subnetIPForwardingNotEnabled":   tsapi.WebhookSubscriptionType("subnetIPForwardingNotEnabled"),
	"exitNodeIPForwardingNotEnabled": tsapi.WebhookSubscriptionType("exitNodeIPForwardingNotEnabled"),
}

func Command() *cobra.Command {
	var (
		rawURL      string
		providerStr string
		subStrs     []string
	)

	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Create a new tailnet webhook",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// validate URL
			if _, err := url.ParseRequestURI(rawURL); err != nil {
				return fmt.Errorf("invalid --url: %v", err)
			}

			// provider validation
			providerStr = strings.ToLower(providerStr)
			if _, ok := providerEnum[providerStr]; !ok {
				valid := keys(providerEnum)
				return fmt.Errorf("invalid --provider %q (valid: %s)", providerStr, strings.Join(valid, ", "))
			}

			// subscription validation
			if len(subStrs) == 0 {
				return fmt.Errorf("at least one --subscription is required")
			}
			for _, s := range subStrs {
				if _, ok := subscriptionEnum[s]; !ok {
					valid := keys(subscriptionEnum)
					return fmt.Errorf("invalid subscription %q (valid: %s)", s, strings.Join(valid, ", "))
				}
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, _ []string) error {
			client, err := tscli.New()
			if err != nil {
				return err
			}

			// convert slices to typed enums
			subs := make([]tsapi.WebhookSubscriptionType, 0, len(subStrs))
			for _, s := range subStrs {
				subs = append(subs, subscriptionEnum[s])
			}

			req := tsapi.CreateWebhookRequest{
				EndpointURL:   rawURL,
				ProviderType:  providerEnum[providerStr],
				Subscriptions: subs,
			}

			hook, err := client.Webhooks().Create(context.Background(), req)
			if err != nil {
				return fmt.Errorf("failed to create webhook: %w", err)
			}

			out, _ := json.MarshalIndent(hook, "", "  ")
			format := viper.GetString("format")
			output.Print(format, out)
			return nil
		},
	}

	cmd.Flags().StringVar(&rawURL, "url", "", "Destination endpoint URL (required)")
	cmd.Flags().StringVar(&providerStr, "provider", "generic", "Provider type: generic|slack|mattermost|googlechat|discord")
	cmd.Flags().StringSliceVar(&subStrs, "subscription", nil, "Event subscription (repeatable)")

	_ = cmd.MarkFlagRequired("url")
	_ = cmd.MarkFlagRequired("subscription")

	return cmd
}

func keys[T any](m map[string]T) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

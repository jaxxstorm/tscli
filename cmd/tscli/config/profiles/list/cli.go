package list

import (
	"encoding/json"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/jaxxstorm/tscli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type listProfile struct {
	Name     string `json:"name"`
	Tailnet  string `json:"tailnet"`
	AuthType string `json:"auth-type"`
	Active   bool   `json:"active"`
}

func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured tailnet profiles",
		Long:  "List all configured tailnet profiles and indicate which profile is active.",
		RunE: func(_ *cobra.Command, _ []string) error {
			state, err := config.ListTailnetProfiles()
			if err != nil {
				return err
			}

			profiles := make([]listProfile, 0, len(state.Tailnets))
			for _, profile := range state.Tailnets {
				profiles = append(profiles, listProfile{
					Name:     profile.Name,
					Tailnet:  profile.EffectiveTailnet(),
					AuthType: profile.AuthType(),
					Active:   profile.Name == state.ActiveTailnet,
				})
			}

			payload := map[string]any{
				"active-tailnet": state.ActiveTailnet,
				"tailnets":       profiles,
			}
			out, _ := json.Marshal(payload)

			return output.Print(viper.GetString("output"), out)
		},
	}
}

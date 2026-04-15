package setup

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	var (
		publicKey         string
		privateKeySource  string
		privateKeyPath    string
		privateKeyCommand string
	)

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configure AGE encryption for persisted secrets",
		Long:  "Prompt for an AGE public key and choose how the AGE private key should be provided at runtime.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			reader := bufio.NewReader(cmd.InOrStdin())

			if strings.TrimSpace(publicKey) == "" {
				fmt.Fprint(cmd.OutOrStdout(), "AGE public key: ")
				value, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				publicKey = strings.TrimSpace(value)
			}

			if strings.TrimSpace(privateKeySource) == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "Private key source [config|env|command]: ")
				value, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				privateKeySource = strings.TrimSpace(value)
			}

			cfg := config.AgeEncryptionConfig{PublicKey: strings.TrimSpace(publicKey)}
			switch strings.ToLower(strings.TrimSpace(privateKeySource)) {
			case "path":
				if strings.TrimSpace(privateKeyPath) == "" {
					fmt.Fprint(cmd.OutOrStdout(), "AGE private key path: ")
					value, err := reader.ReadString('\n')
					if err != nil {
						return err
					}
					privateKeyPath = strings.TrimSpace(value)
				}
				cfg.PrivateKeyPath = strings.TrimSpace(privateKeyPath)
			case "command":
				if strings.TrimSpace(privateKeyCommand) == "" {
					fmt.Fprint(cmd.OutOrStdout(), "Private key command: ")
					value, err := reader.ReadString('\n')
					if err != nil {
						return err
					}
					privateKeyCommand = strings.TrimSpace(value)
				}
				cfg.PrivateKeyCommand = strings.TrimSpace(privateKeyCommand)
			case "env":
				// Public key only; runtime decryption will use TSCLI_AGE_PRIVATE_KEY.
			case "":
				return fmt.Errorf("private key source is required")
			default:
				return fmt.Errorf("private key source must be one of: path, env, command")
			}

			if err := config.SetAgeEncryptionConfig(cfg); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "config encryption saved")
			return nil
		},
	}

	cmd.Flags().StringVar(&publicKey, "public-key", "", "AGE public key used to encrypt persisted secrets")
	cmd.Flags().StringVar(&privateKeySource, "private-key-source", "", "How to provide the AGE private key: path, env, or command")
	cmd.Flags().StringVar(&privateKeyPath, "private-key-path", "", "Path to an AGE private key file when --private-key-source=path")
	cmd.Flags().StringVar(&privateKeyCommand, "private-key-command", "", "Command that returns the AGE private key when --private-key-source=command")

	return cmd
}

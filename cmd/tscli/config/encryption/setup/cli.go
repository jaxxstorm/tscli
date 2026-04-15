package setup

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jaxxstorm/tscli/pkg/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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
			interactive := isTerminalReader(cmd.InOrStdin()) && isTerminalWriter(cmd.OutOrStdout())

			if strings.TrimSpace(privateKeySource) == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "Private key source [path|env|command]: ")
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
				resolvedCfg, message, err := resolvePathBasedConfig(reader, interactive, cmd.OutOrStdout(), strings.TrimSpace(publicKey), cfg.PrivateKeyPath)
				if err != nil {
					return err
				}
				cfg = resolvedCfg
				if message != "" {
					fmt.Fprintln(cmd.OutOrStdout(), message)
				}
			case "command":
				if strings.TrimSpace(publicKey) == "" {
					fmt.Fprint(cmd.OutOrStdout(), "AGE public key: ")
					value, err := reader.ReadString('\n')
					if err != nil {
						return err
					}
					publicKey = strings.TrimSpace(value)
				}
				cfg.PublicKey = strings.TrimSpace(publicKey)
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
				if strings.TrimSpace(publicKey) == "" {
					fmt.Fprint(cmd.OutOrStdout(), "AGE public key: ")
					value, err := reader.ReadString('\n')
					if err != nil {
						return err
					}
					publicKey = strings.TrimSpace(value)
				}
				cfg.PublicKey = strings.TrimSpace(publicKey)
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
	_ = cmd.RegisterFlagCompletionFunc("private-key-source", func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"path", "env", "command"}, cobra.ShellCompDirectiveNoFileComp
	})
	_ = cmd.MarkFlagFilename("private-key-path")

	return cmd
}

func resolvePathBasedConfig(reader *bufio.Reader, interactive bool, out io.Writer, publicKey string, privateKeyPath string) (config.AgeEncryptionConfig, string, error) {
	privateKeyPath = strings.TrimSpace(privateKeyPath)
	cfg := config.AgeEncryptionConfig{PrivateKeyPath: privateKeyPath}

	inspected, err := config.InspectAgeIdentityFile(privateKeyPath)
	switch {
	case err == nil:
		reuseExisting := !interactive
		if interactive {
			fmt.Fprintf(out, "Reuse existing AGE identity at %s? [yes|no]: ", inspected.Path)
			value, readErr := reader.ReadString('\n')
			if readErr != nil {
				return config.AgeEncryptionConfig{}, "", readErr
			}
			reuseExisting = strings.EqualFold(strings.TrimSpace(value), "yes") || strings.EqualFold(strings.TrimSpace(value), "y")
		}

		if reuseExisting {
			if publicKey != "" && strings.TrimSpace(publicKey) != inspected.PublicKey {
				return config.AgeEncryptionConfig{}, "", fmt.Errorf("provided AGE public key does not match existing identity at %s", inspected.Path)
			}
			cfg.PublicKey = inspected.PublicKey
			cfg.PrivateKeyPath = inspected.Path
			return cfg, fmt.Sprintf("Reusing existing AGE identity at %s", inspected.Path), nil
		}
	case os.IsNotExist(err):
		// No reusable file exists; continue with standard path-based setup.
	case err != nil:
		fmt.Fprintf(out, "Existing key file at %s could not be reused: %v\n", privateKeyPath, err)
	}

	if strings.TrimSpace(publicKey) == "" {
		fmt.Fprint(out, "AGE public key: ")
		value, err := reader.ReadString('\n')
		if err != nil {
			return config.AgeEncryptionConfig{}, "", err
		}
		publicKey = strings.TrimSpace(value)
	}
	cfg.PublicKey = strings.TrimSpace(publicKey)
	return cfg, "", nil
}

func isTerminalReader(r io.Reader) bool {
	file, ok := r.(*os.File)
	return ok && term.IsTerminal(int(file.Fd()))
}

func isTerminalWriter(w io.Writer) bool {
	file, ok := w.(*os.File)
	return ok && term.IsTerminal(int(file.Fd()))
}

package cmd

import (
	"fmt"

	"github.com/chawza/memos-cli/internal/api"
	"github.com/chawza/memos-cli/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication configuration",
	}

	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Save authentication credentials to config file",
		RunE:  runAuthSet,
	}
	setCmd.Flags().String("base-url", "", "Memos instance base URL")
	setCmd.Flags().String("token", "", "Memos access token")
	setCmd.Flags().Int("timeout", 0, "HTTP timeout in seconds (default: 30)")
	setCmd.Flags().Bool("tls-skip-verify", false, "Skip TLS certificate verification")
	_ = setCmd.MarkFlagRequired("base-url")
	_ = setCmd.MarkFlagRequired("token")

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Verify saved configuration is valid",
		RunE:  runAuthCheck,
	}
	checkCmd.Flags().String("base-url", "", "Override base URL (optional)")
	checkCmd.Flags().String("token", "", "Override token (optional)")

	authCmd.AddCommand(setCmd, checkCmd)
	rootCmd.AddCommand(authCmd)
}

func runAuthSet(cmd *cobra.Command, args []string) error {
	baseURL, _ := cmd.Flags().GetString("base-url")
	token, _ := cmd.Flags().GetString("token")
	timeout, _ := cmd.Flags().GetInt("timeout")
	tlsSkipVerify, _ := cmd.Flags().GetBool("tls-skip-verify")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if cfg == nil {
		cfg = &config.Config{}
	}

	cfg.BaseURL = baseURL
	cfg.Token = token
	cfg.Timeout = timeout
	cfg.TLSSkipVerify = tlsSkipVerify

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	opts := clientOptionsFromConfig(cfg)
	client := api.NewClient(baseURL, token, opts...)
	if err := client.Ping(); err != nil {
		cmd.Printf("Config saved to ~/.config/memos-cli/config.toml\n")
		cmd.Printf("Warning: connectivity check failed: %v\n", err)
		return nil
	}

	cmd.Printf("Config saved to ~/.config/memos-cli/config.toml\n")
	cmd.Printf("Successfully authenticated to %s\n", baseURL)
	return nil
}

func runAuthCheck(cmd *cobra.Command, args []string) error {
	baseURL, _ := cmd.Flags().GetString("base-url")
	token, _ := cmd.Flags().GetString("token")

	var opts []api.ClientOption

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if cfg != nil {
		if baseURL == "" {
			baseURL = cfg.BaseURL
		}
		if token == "" {
			token = cfg.Token
		}
		opts = clientOptionsFromConfig(cfg)
	}

	if baseURL == "" || token == "" {
		return fmt.Errorf("no configuration found. Run `memos auth set` first")
	}

	client := api.NewClient(baseURL, token, opts...)
	if err := client.Ping(); err != nil {
		cmd.Printf("Configuration is INVALID: %v\n", err)
		return nil
	}
	cmd.Printf("Configuration is valid. Connected to %s\n", baseURL)
	return nil
}

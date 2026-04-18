package cmd

import (
	"fmt"
	"os"

	"github.com/chawza/memos-cli/internal/api"
	"github.com/chawza/memos-cli/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "memos",
	Short: "CLI wrapper for the Memos API",
	Long: `memos-cli is a fast, minimal CLI for managing your self-hosted Memos notes.

Configure using environment variables, flags, or a config file:
  export MEMOS_BASE_URL="https://memos.example.com"
  export MEMOS_TOKEN="your-access-token"

Or use flags:
  memos --base-url https://memos.example.com --token xxx list`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().String("base-url", "", "Memos instance base URL (env: MEMOS_BASE_URL)")
	rootCmd.PersistentFlags().String("token", "", "Memos access token (env: MEMOS_TOKEN)")
}

func resolveClient(cmd *cobra.Command) (*api.Client, error) {
	baseURL, _ := cmd.Flags().GetString("base-url")
	token, _ := cmd.Flags().GetString("token")

	if baseURL == "" {
		baseURL = os.Getenv("MEMOS_BASE_URL")
	}
	if token == "" {
		token = os.Getenv("MEMOS_TOKEN")
	}

	if baseURL == "" || token == "" {
		cfg, err := config.Load()
		if err != nil {
			return nil, fmt.Errorf("load config: %w", err)
		}
		if cfg != nil {
			if baseURL == "" {
				baseURL = cfg.BaseURL
			}
			if token == "" {
				token = cfg.Token
			}
		}
	}

	if baseURL == "" {
		return nil, fmt.Errorf("base URL not set: use --base-url, MEMOS_BASE_URL, or config file")
	}
	if token == "" {
		return nil, fmt.Errorf("token not set: use --token, MEMOS_TOKEN, or config file")
	}

	return api.NewClient(baseURL, token), nil
}

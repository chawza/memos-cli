package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "memos",
	Short: "memos-cli — CLI wrapper for the Memos API",
	Long: `memos-cli is a fast, minimal CLI for managing your self-hosted Memos notes.

Configure once using environment variables or flags:
  export MEMOS_BASE_URL="https://memos.example.com"
  export MEMOS_TOKEN="your-access-token"

Then run any memos command.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().String("base-url", envOr("MEMOS_BASE_URL", ""), "Memos instance base URL")
	rootCmd.PersistentFlags().String("token", envOr("MEMOS_TOKEN", ""), "Memos access token")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// resolveClient returns the API client from config + flags.
func resolveClient(cmd *cobra.Command) (*client, error) {
	baseURL, _ := cmd.Flags().GetString("base-url")
	token, _ := cmd.Flags().GetString("token")

	if baseURL == "" || token == "" {
		cfg, err := loadConfig()
		if err != nil {
			return nil, err
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
		return nil, fmt.Errorf("base URL not set: use --base-url or MEMOS_BASE_URL")
	}
	if token == "" {
		return nil, fmt.Errorf("token not set: use --token or MEMOS_TOKEN")
	}

	return newClient(baseURL, token), nil
}

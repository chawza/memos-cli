package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var searchCmd = &cobra.Command{
		Use:   "search <query>",
		Short: "Search memos by keyword",
		Args:  cobra.ExactArgs(1),
		RunE:  runSearch,
	}
	searchCmd.Flags().Int("limit", 20, "Max results")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")

	// Memos doesn't have full-text search built into the API filter (AIP-160),
	// so we list with a generous limit and filter client-side.
	// The API supports: creatorId, visibility, rowStatus, pinned filters only.
	memos, _, err := c.ListMemos(limit, "", "")
	if err != nil {
		return err
	}

	query := strings.ToLower(args[0])
	for _, m := range memos {
		if strings.Contains(strings.ToLower(m.Content), query) {
			prefix := "  "
			if m.Pinned {
				prefix = "📌 "
			}
			cmd.Printf("%s[%s] %s\n", prefix, m.Visibility, m.Content)
		}
	}
	return nil
}

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List memos",
		RunE:  runList,
	}
	listCmd.Flags().Int("limit", 20, "Max memos to return")
	listCmd.Flags().String("filter", "", "CEL filter expression")
	listCmd.Flags().String("state", "", "State filter: NORMAL or ARCHIVED")
	listCmd.Flags().StringP("output", "o", "text", "Output format: text, json, table")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	limit, _ := cmd.Flags().GetInt("limit")
	filter, _ := cmd.Flags().GetString("filter")
	state, _ := cmd.Flags().GetString("state")
	output, _ := cmd.Flags().GetString("output")

	memos, _, err := c.ListMemos(limit, "", filter, state)
	if err != nil {
		return err
	}

	return printMemoList(cmd, memos, output)
}

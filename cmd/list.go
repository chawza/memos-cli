package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func init() {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List memos",
		RunE:  runList,
	}
	listCmd.Flags().Int("limit", 20, "Max memos to return")
	listCmd.Flags().String("filter", "", "AIP-160 filter expression")
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
	output, _ := cmd.Flags().GetString("output")

	memos, _, err := c.ListMemos(limit, "", filter)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(memos)
	case "table":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tCONTENT\tVISIBILITY\tPINNED")
		for _, m := range memos {
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", m.ID, m.Content, m.Visibility, m.Pinned)
		}
		return w.Flush()
	default: // text
		for _, m := range memos {
			prefix := "  "
			if m.Pinned {
				prefix = "📌 "
			}
			fmt.Printf("%s[%s] %s\n", prefix, m.Visibility, m.Content)
		}
	}
	return nil
}

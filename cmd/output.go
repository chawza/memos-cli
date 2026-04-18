package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/chawza/memos-cli/internal/api"
	"github.com/spf13/cobra"
)

func printMemo(cmd *cobra.Command, memo *api.Memo, format string) error {
	switch format {
	case "json":
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(memo)
	default:
		fmt.Fprintf(cmd.OutOrStdout(), "Name:       %s\n", memo.Name)
		fmt.Fprintf(cmd.OutOrStdout(), "State:      %s\n", memo.State)
		fmt.Fprintf(cmd.OutOrStdout(), "Visibility: %s\n", memo.Visibility)
		fmt.Fprintf(cmd.OutOrStdout(), "Pinned:     %v\n", memo.Pinned)
		fmt.Fprintf(cmd.OutOrStdout(), "Creator:    %s\n", memo.Creator)
		if memo.CreateTime != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Created:    %s\n", memo.CreateTime)
		}
		if memo.UpdateTime != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "Updated:    %s\n", memo.UpdateTime)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", memo.Content)
		return nil
	}
}

func printMemoList(cmd *cobra.Command, memos []api.Memo, format string) error {
	switch format {
	case "json":
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		return enc.Encode(memos)
	case "table":
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVISIBILITY\tPINNED\tCONTENT")
		for _, m := range memos {
			content := m.Content
			if len(content) > 60 {
				content = content[:57] + "..."
			}
			content = strings.ReplaceAll(content, "\n", " ")
			fmt.Fprintf(w, "%s\t%s\t%v\t%s\n", m.Name, m.Visibility, m.Pinned, content)
		}
		return w.Flush()
	default:
		for _, m := range memos {
			pin := "  "
			if m.Pinned {
				pin = "* "
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s[%s] %s\n", pin, m.Visibility, strings.ReplaceAll(m.Content, "\n", " "))
		}
		return nil
	}
}

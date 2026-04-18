package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	var getCmd = &cobra.Command{
		Use:   "get <id>",
		Short: "Get a memo by ID",
		Args:  cobra.ExactArgs(1),
		RunE:  runGet,
	}
	rootCmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	memo, err := c.GetMemo(args[0])
	if err != nil {
		return err
	}

	fmt.Printf("ID:         %s\n", memo.ID)
	fmt.Printf("Visibility: %s\n", memo.Visibility)
	fmt.Printf("Pinned:     %v\n", memo.Pinned)
	fmt.Printf("Created:    %s\n", memo.CreateTime)
	fmt.Printf("Updated:    %s\n", memo.UpdateTime)
	fmt.Printf("\n%s\n", memo.Content)
	return nil
}

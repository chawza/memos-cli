package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	var deleteCmd = &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a memo by ID",
		Args:  cobra.ExactArgs(1),
		RunE:  runDelete,
	}
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	if err := c.DeleteMemo(args[0]); err != nil {
		return err
	}

	cmd.Printf("Deleted memo %s\n", args[0])
	return nil
}

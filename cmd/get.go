package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	var getCmd = &cobra.Command{
		Use:   "get <id>",
		Short: "Get a memo by ID",
		Args:  cobra.ExactArgs(1),
		RunE:  runGet,
	}
	getCmd.Flags().StringP("output", "o", "text", "Output format: text, json")
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

	output, _ := cmd.Flags().GetString("output")
	return printMemo(cmd, memo, output)
}

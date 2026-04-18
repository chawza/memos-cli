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
	getCmd.Flags().Bool("include-comments", false, "Include comments")
	getCmd.Flags().Bool("include-reactions", false, "Include reactions")
	getCmd.Flags().Bool("include-attachments", false, "Include attachments")
	getCmd.Flags().Bool("a", false, "Include all (comments, reactions, attachments)")
	memoCmd.AddCommand(getCmd)
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

	includeComments, _ := cmd.Flags().GetBool("include-comments")
	includeReactions, _ := cmd.Flags().GetBool("include-reactions")
	includeAttachments, _ := cmd.Flags().GetBool("include-attachments")
	includeAll, _ := cmd.Flags().GetBool("a")

	if err := printMemo(cmd, memo, output); err != nil {
		return err
	}

	if includeAll {
		includeComments = true
		includeReactions = true
		includeAttachments = true
	}

	if includeComments {
		comments, _, err := c.ListMemoComments(args[0])
		if err != nil {
			return err
		}
		if err := printComments(cmd, comments); err != nil {
			return err
		}
	}

	if includeReactions {
		reactions, _, err := c.ListMemoReactions(args[0])
		if err != nil {
			return err
		}
		if err := printReactions(cmd, reactions); err != nil {
			return err
		}
	}

	if includeAttachments {
		attachments, _, err := c.ListMemoAttachments(args[0])
		if err != nil {
			return err
		}
		if err := printAttachments(cmd, attachments); err != nil {
			return err
		}
	}

	return nil
}
package cmd

import (
	"github.com/chawza/memos-cli/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	var reactionsCmd = &cobra.Command{
		Use:   "reactions",
		Short: "Manage memo reactions",
	}
	rootCmd.AddCommand(reactionsCmd)

	var listReactionsCmd = &cobra.Command{
		Use:   "list <memo-id>",
		Short: "List reactions for a memo",
		Args:  cobra.ExactArgs(1),
		RunE:  runListReactions,
	}
	reactionsCmd.AddCommand(listReactionsCmd)

	var createReactionCmd = &cobra.Command{
		Use:   "create <memo-id> --type <emoji>",
		Short: "Add a reaction to a memo",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreateReaction,
	}
	createReactionCmd.Flags().String("type", "", "Reaction type (emoji)")
	createReactionCmd.MarkFlagRequired("type")
	reactionsCmd.AddCommand(createReactionCmd)

	var deleteReactionCmd = &cobra.Command{
		Use:   "delete <memo-id> <reaction-id>",
		Short: "Remove a reaction from a memo",
		Args:  cobra.ExactArgs(2),
		RunE:  runDeleteReaction,
	}
	reactionsCmd.AddCommand(deleteReactionCmd)
}

func runListReactions(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	reactions, _, err := c.ListMemoReactions(args[0])
	if err != nil {
		return err
	}

	return printReactions(cmd, reactions)
}

func runCreateReaction(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	reactionType, _ := cmd.Flags().GetString("type")
	reaction, err := c.UpsertMemoReaction(args[0], &api.UpsertReaction{ReactionType: reactionType})
	if err != nil {
		return err
	}

	cmd.Println("Reaction added:", reaction.ReactionType)
	return nil
}

func runDeleteReaction(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	return c.DeleteMemoReaction(args[0], args[1])
}
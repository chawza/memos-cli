package cmd

import (
	"github.com/chawza/memos-cli/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	var commentsCmd = &cobra.Command{
		Use:   "comments",
		Short: "Manage memo comments",
	}
	rootCmd.AddCommand(commentsCmd)

	var listCommentsCmd = &cobra.Command{
		Use:   "list <memo-id>",
		Short: "List comments for a memo",
		Args:  cobra.ExactArgs(1),
		RunE:  runListComments,
	}
	commentsCmd.AddCommand(listCommentsCmd)

	var createCommentCmd = &cobra.Command{
		Use:   "create <memo-id> --content <text>",
		Short: "Create a comment on a memo",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreateComment,
	}
	createCommentCmd.Flags().String("content", "", "Comment content")
	createCommentCmd.MarkFlagRequired("content")
	commentsCmd.AddCommand(createCommentCmd)

	var deleteCommentCmd = &cobra.Command{
		Use:   "delete <comment-id>",
		Short: "Delete a comment",
		Args:  cobra.ExactArgs(1),
		RunE:  runDeleteComment,
	}
	commentsCmd.AddCommand(deleteCommentCmd)
}

func runListComments(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	comments, _, err := c.ListMemoComments(args[0])
	if err != nil {
		return err
	}

	return printComments(cmd, comments)
}

func runCreateComment(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	content, _ := cmd.Flags().GetString("content")
	comment, err := c.CreateMemoComment(args[0], &api.CreateMemo{Content: content})
	if err != nil {
		return err
	}

	return printMemo(cmd, comment, "text")
}

func runDeleteComment(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	return c.DeleteMemo(args[0])
}
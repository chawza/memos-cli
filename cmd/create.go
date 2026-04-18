package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new memo",
		RunE:  runCreate,
	}
	createCmd.Flags().StringP("content", "c", "", "Memo content (markdown)")
	createCmd.Flags().String("visibility", "PRIVATE", "Visibility: PRIVATE, PROTECTED, PUBLIC")
	createCmd.Flags().Bool("pinned", false, "Pin the memo")
	createCmd.MarkFlagRequired("content")
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	content, _ := cmd.Flags().GetString("content")
	visibility, _ := cmd.Flags().GetString("visibility")
	pinned, _ := cmd.Flags().GetBool("pinned")

	memo, err := c.CreateMemo(&createMemoRequest{
		Content:    content,
		Visibility: visibility,
		Pinned:     pinned,
	})
	if err != nil {
		return err
	}

	cmd.Printf("Created memo %s\n", memo.ID)
	return nil
}

type createMemoRequest = api.CreateMemoRequest

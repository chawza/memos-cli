package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Note: join() is defined in helpers.go

func init() {
	var updateCmd = &cobra.Command{
		Use:   "update <id>",
		Short: "Update a memo",
		Args:  cobra.ExactArgs(1),
		RunE:  runUpdate,
	}
	updateCmd.Flags().StringP("content", "c", "", "New content")
	updateCmd.Flags().String("visibility", "", "New visibility")
	updateCmd.Flags().Bool("pinned", false, "Set pinned state")
	updateCmd.Flags().Bool("unpin", false, "Unpin the memo")
	updateCmd.Flags().String("state", "", "State: NORMAL or ARCHIVED")
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	content, _ := cmd.Flags().GetString("content")
	visibility, _ := cmd.Flags().GetString("visibility")
	pinned, _ := cmd.Flags().GetBool("pinned")
	unpin, _ := cmd.Flags().GetBool("unpin")
	state, _ := cmd.Flags().GetString("state")

	req := &updateMemoRequest{}
	var masks []string

	if cmd.Flags().Changed("content") {
		req.Content = &content
		masks = append(masks, "content")
	}
	if cmd.Flags().Changed("visibility") {
		req.Visibility = &visibility
		masks = append(masks, "visibility")
	}
	if cmd.Flags().Changed("pinned") {
		req.Pinned = &pinned
		masks = append(masks, "pinned")
	}
	if cmd.Flags().Changed("unpin") {
		p := false
		req.Pinned = &p
		masks = append(masks, "pinned")
	}
	if cmd.Flags().Changed("state") {
		req.State = &state
		masks = append(masks, "state")
	}

	if len(masks) == 0 {
		return fmt.Errorf("no fields to update (use flags)")
	}

	memo, err := c.UpdateMemo(args[0], req, join(masks, ","))
	if err != nil {
		return err
	}

	fmt.Printf("Updated memo %s\n", memo.ID)
	return nil
}

type updateMemoRequest = api.UpdateMemoRequest

package cmd

import (
	"fmt"

	"github.com/chawza/memos-cli/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	var updateCmd = &cobra.Command{
		Use:   "update <id>",
		Short: "Update a memo",
		Args:  cobra.ExactArgs(1),
		RunE:  runUpdate,
	}
	updateCmd.Flags().StringP("content", "c", "", "New content")
	updateCmd.Flags().String("visibility", "", "New visibility: PRIVATE, PROTECTED, PUBLIC")
	updateCmd.Flags().Bool("pinned", false, "Pin the memo")
	updateCmd.Flags().Bool("unpin", false, "Unpin the memo")
	updateCmd.Flags().String("state", "", "State: NORMAL or ARCHIVED")
	memoCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	content, _ := cmd.Flags().GetString("content")
	visibility, _ := cmd.Flags().GetString("visibility")
	pinned, _ := cmd.Flags().GetBool("pinned")
	state, _ := cmd.Flags().GetString("state")

	update := &api.UpdateMemo{
		Name: "memos/" + args[0],
	}
	var masks []string

	if cmd.Flags().Changed("content") {
		update.Content = &content
		masks = append(masks, "content")
	}
	if cmd.Flags().Changed("visibility") {
		update.Visibility = &visibility
		masks = append(masks, "visibility")
	}
	if cmd.Flags().Changed("pinned") {
		update.Pinned = &pinned
		masks = append(masks, "pinned")
	}
	if cmd.Flags().Changed("unpin") {
		val := false
		update.Pinned = &val
		masks = append(masks, "pinned")
	}
	if cmd.Flags().Changed("state") {
		update.State = &state
		masks = append(masks, "state")
	}

	if len(masks) == 0 {
		return fmt.Errorf("no fields to update; use flags to specify what to change")
	}

	memo, err := c.UpdateMemo(args[0], update, joinMasks(masks))
	if err != nil {
		return err
	}

	cmd.Printf("Updated memo %s\n", memo.Name)
	return nil
}

func joinMasks(masks []string) string {
	result := masks[0]
	for _, m := range masks[1:] {
		result += "," + m
	}
	return result
}

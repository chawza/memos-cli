package cmd

import (
	"github.com/chawza/memos-cli/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	var attachmentsCmd = &cobra.Command{
		Use:   "attachments",
		Short: "Manage memo attachments",
	}
	rootCmd.AddCommand(attachmentsCmd)

	var listAttachmentsCmd = &cobra.Command{
		Use:   "list <memo-id>",
		Short: "List attachments for a memo",
		Args:  cobra.ExactArgs(1),
		RunE:  runListAttachments,
	}
	attachmentsCmd.AddCommand(listAttachmentsCmd)

	var setAttachmentsCmd = &cobra.Command{
		Use:   "set <memo-id> --file <path>",
		Short: "Set attachments for a memo (replaces all existing)",
		Args:  cobra.ExactArgs(1),
		RunE:  runSetAttachments,
	}
	setAttachmentsCmd.Flags().StringArray("file", []string{}, "File path(s) to attach")
	setAttachmentsCmd.MarkFlagRequired("file")
	attachmentsCmd.AddCommand(setAttachmentsCmd)
}

func runListAttachments(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	attachments, _, err := c.ListMemoAttachments(args[0])
	if err != nil {
		return err
	}

	return printAttachments(cmd, attachments)
}

func runSetAttachments(cmd *cobra.Command, args []string) error {
	c, err := resolveClient(cmd)
	if err != nil {
		return err
	}

	files, _ := cmd.Flags().GetStringArray("file")
	if len(files) == 0 {
		cmd.Println("No files to attach")
		return nil
	}

	var attachments []api.Attachment
	for _, f := range files {
		attachments = append(attachments, api.Attachment{
			Filename: f,
		})
	}

	return c.SetMemoAttachments(args[0], attachments)
}
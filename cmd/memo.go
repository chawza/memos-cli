package cmd

import "github.com/spf13/cobra"

var memoCmd = &cobra.Command{
	Use:   "memo",
	Short: "Manage memos",
}

func init() {
	rootCmd.AddCommand(memoCmd)
}
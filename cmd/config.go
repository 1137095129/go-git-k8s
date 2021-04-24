package cmd

import "github.com/spf13/cobra"

var configCmd = &cobra.Command{
	Use: "config",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
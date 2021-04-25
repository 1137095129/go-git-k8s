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

var mainCmd = &cobra.Command{
	Use: "add",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(mainCmd)
}
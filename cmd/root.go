package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use: "git-k8s",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {

	},
}

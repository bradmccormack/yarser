package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of yarser",
	Long:  `All software has versions. This is yarser's`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO Make this use semver + git commit SHA or tag.
		fmt.Println("0.0.2")
	},
}

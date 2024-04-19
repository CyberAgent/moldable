package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	pkgVersion string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of Moldable",
	Run: func(cmd *cobra.Command, args []string) {
		getVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func getVersion() {
	fmt.Println(pkgVersion)
}

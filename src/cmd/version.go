package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
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
	log.Println(os.Getenv("PKG_VERSION"))
}

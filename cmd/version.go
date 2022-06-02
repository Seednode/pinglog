package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string = "0.5.0"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  "Print the version number of pinglog",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pinglog v" + Version)
	},
}

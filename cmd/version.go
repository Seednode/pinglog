/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var Version = "0.19.1"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  "Print the version number of pinglog",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := fmt.Println("pinglog v" + Version)
		if err != nil {
			log.Fatal(err)
		}
	},
}

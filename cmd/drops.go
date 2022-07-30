/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"github.com/spf13/cobra"
)

func DroppedPackets(arguments []string) {
}

var droppedCmd = &cobra.Command{
	Use:   "dropped <file1> [file2]...",
	Short: "Parse out timestamps of drops from log file(s)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		DroppedPackets(args)
	},
}

//func init() {
//	rootCmd.AddCommand(stripCmd)
//}

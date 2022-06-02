/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/acarl005/stripansi"
	"github.com/spf13/cobra"
)

func StripFile(logFile string) {
	fmt.Println("Stripping file " + logFile)

	file, err := os.Open(logFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strippedLine := stripansi.Strip(scanner.Text())
		fmt.Println(strippedLine)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func StripColors(arguments []string) {
	for file := 0; file < len(arguments); file++ {
		StripFile(arguments[file])
	}
}

var stripCmd = &cobra.Command{
	Use:   "strip <file1> [file2]...",
	Short: "Strip ANSI color codes from log file(s)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		StripColors(args)
	},
}

func init() {
	rootCmd.AddCommand(stripCmd)
}

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

func StripColors(args []string) error {
	file, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strippedLine := stripansi.Strip(scanner.Text())
		_, err := fmt.Println(strippedLine)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

var stripCmd = &cobra.Command{
	Use:   "strip <file>",
	Short: "Strip ANSI color codes from log file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := StripColors(args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(stripCmd)
}

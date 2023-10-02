/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

const escapeSequences = "\x1b[[0-9;]*m|\x07"

func Strip(input string, regex *regexp.Regexp) string {
	return regex.ReplaceAllString(input, "")
}

func StripColors(args []string) error {
	file, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer file.Close()

	regex := regexp.MustCompile(escapeSequences)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strippedLine := Strip(scanner.Text(), regex)

		_, err := fmt.Println(strippedLine)
		if err != nil {
			return err
		}
	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	return nil
}

var stripCmd = &cobra.Command{
	Use:   "strip <file>",
	Short: "Strip ANSI color codes from log file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := StripColors(args)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stripCmd)
}

/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

const escapeSequences = "\x1b[[0-9;]*m|\x07"

var re = regexp.MustCompile(escapeSequences)

func Strip(input string) string {
	return re.ReplaceAllString(input, "")
}

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
		strippedLine := Strip(scanner.Text())
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

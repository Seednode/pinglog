/*
Copyright © 2025 Seednode <seednode@seedno.de>
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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

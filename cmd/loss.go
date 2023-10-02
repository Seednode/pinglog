/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func parseTime(line string) (string, error) {
	fields := strings.Fields(line)

	t, err := time.Parse(DATE, strings.Join(fields[:3], " "))
	if err != nil {
		return "", err
	}

	return t.Format(DATE), nil
}

func calculateLoss(logFile string) error {
	file, err := os.Open(logFile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Printf("%v:\n", logFile)
	if err != nil {
		return err
	}

	var lostLastPacket bool
	var lastTimestamp string
	var lostPacketCount int
	var startTime string
	var endTime string

	regex := regexp.MustCompile(escapeSequences)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stripped := Strip(scanner.Text(), regex)

		if !strings.Contains(stripped, "|") {
			continue
		}

		timestamp, err := parseTime(stripped)
		if err != nil {
			return err
		}

		lostThisPacket := strings.Contains(stripped, "lost or arrived out of order")

		switch {
		case lostThisPacket && !lostLastPacket:
			startTime = lastTimestamp
			endTime = startTime
			if err != nil {
				return err
			}

			lostPacketCount = 1
			lostLastPacket = true
		case lostThisPacket && lostLastPacket:
			lostPacketCount++
			lostLastPacket = true
		case !lostThisPacket && lostLastPacket:
			endTime = timestamp
			if err != nil {
				return err
			}

			fmt.Printf("%v => %v [%v packet(s) lost]\n", startTime, endTime, lostPacketCount)
			lostLastPacket = false
		}

		lastTimestamp = timestamp
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if lostPacketCount == 0 {
		fmt.Println("No dropped packets found")
	} else {
		fmt.Println()
	}

	return nil
}

var lossCmd = &cobra.Command{
	Use:   "loss <file1> [file2]...",
	Short: "Calculate periods of packet loss from log file(s)",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for file := 0; file < len(args); file++ {
			err := calculateLoss(args[file])
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lossCmd)
}

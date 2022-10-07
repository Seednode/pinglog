/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/spf13/cobra"
)

func parseTime(line string) (time.Time, error) {
	fields := strings.Fields(line)

	t, err := time.Parse(DATE, fields[0]+" "+fields[1]+" "+fields[2])
	if err != nil {
		fmt.Printf("Failed to parse time from '%v'\n", line)
		return time.Time{}, err
	}

	return t, nil
}

func CalculateLoss(logFile string) error {
	file, err := os.Open(logFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	_, err = fmt.Printf("%v:\n", logFile)
	if err != nil {
		return err
	}

	var lostLastPacket bool = false
	var lastTimestamp = time.Time{}
	var lostPacketCount int = 0
	var startTime = time.Time{}
	var endTime = time.Time{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stripped := stripansi.Strip(scanner.Text())

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
			endTime, err = parseTime(stripped)
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

	switch lostPacketCount {
	case 0:
		fmt.Println("No dropped packets found")
	default:
		fmt.Println()
	}

	return nil
}

func Loss(arguments []string) error {
	for file := 0; file < len(arguments); file++ {
		err := CalculateLoss(arguments[file])
		if err != nil {
			return err
		}
	}

	return nil
}

var lossCmd = &cobra.Command{
	Use:   "loss <file1> [file2]...",
	Short: "Calculate periods of packet loss from log file(s)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := Loss(args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lossCmd)
}

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

	var lastHadLoss bool = false
	var lastTimestamp = time.Time{}
	var lostPackets int = 0
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

		if strings.Contains(stripped, "lost or arrived out of order") {
			if !lastHadLoss {
				startTime = lastTimestamp
				endTime = startTime
				if err != nil {
					return err
				}

				lostPackets = 1
			} else {
				endTime, err = parseTime(stripped)
				if err != nil {
					return err
				}

				lostPackets++
			}

			lastHadLoss = true
		} else if lastHadLoss {
			endTime, err = parseTime(stripped)
			if err != nil {
				return err
			}

			if lostPackets == 1 {
				fmt.Printf("%v\n", endTime)
			} else {
				fmt.Printf("%v => %v [%v packets lost]\n", startTime, endTime, lostPackets)
			}

			lastHadLoss = false
		}

		lastTimestamp = timestamp
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	switch lostPackets {
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

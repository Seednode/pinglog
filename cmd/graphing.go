/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"fmt"
	"math"
	"net/netip"
	"os"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/go-chart"
)

func calculateXAxis(slice []time.Time) []time.Time {
	r := []time.Time{}

	for i := 0; i < len(slice); i++ {
		r = append(r, slice[i])
	}

	return r
}

func parseFile(logFile string) (string, netip.Addr, []time.Time, []float64, error) {
	var hostname string
	var ipaddr netip.Addr
	var timestamps []time.Time
	var rtts []float64

	file, err := os.Open(logFile)
	if err != nil {
		return "", netip.Addr{}, nil, nil, err
	}
	defer func(file *os.File) error {
		err := file.Close()
		if err != nil {
			return err
		}

		return nil
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		stripped := stripansi.Strip(scanner.Text())

		if strings.Contains(stripped, "PING") {
			fields := strings.Fields(stripped)

			hostname = fields[1]
			ipaddr, err = netip.ParseAddr(strings.Trim(fields[2], "()"))
			if err != nil {
				return "", netip.Addr{}, nil, nil, err
			}
		}

		if strings.Contains(stripped, "|") && !strings.Contains(stripped, "lost or arrived out of order") {
			fields := strings.Fields(stripped)

			t, err := time.Parse(DATE, fields[0]+" "+fields[1]+" "+fields[2])
			if err != nil {
				return "", netip.Addr{}, nil, nil, err
			}
			timestamps = append(timestamps, t)

			r, err := time.ParseDuration(strings.TrimPrefix(fields[10], "time="))
			p := r.Milliseconds()
			if err != nil {
				return "", netip.Addr{}, nil, nil, err
			}
			rtts = append(rtts, float64(p))
		}
	}

	if err := scanner.Err(); err != nil {
		return "", netip.Addr{}, nil, nil, err
	}

	return hostname, ipaddr, timestamps, rtts, nil
}

func createLineChart(args []string) {
	hostname, ipaddr, timestamps, rtts, err := parseFile(args[0])
	if err != nil {
		panic(err)
	}

	graph := chart.Chart{
		Title:      fmt.Sprintf("Pings to %v (%v):", hostname, ipaddr),
		TitleStyle: chart.Style{Show: true},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 100,
			},
		},
		XAxis: chart.XAxis{
			Name:           "Time",
			NameStyle:      chart.Style{Show: true},
			Style:          chart.Style{Show: true},
			ValueFormatter: chart.TimeMinuteValueFormatter,
		},
		YAxis: chart.YAxis{
			Name:      "Round-trip time (ms)",
			NameStyle: chart.Style{Show: true},
			Style:     chart.Style{Show: true},
			ValueFormatter: func(v interface{}) string {
				if vf, isFloat := v.(float64); isFloat {
					return fmt.Sprint(math.Round(vf))
				}
				return ""
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Style: chart.Style{
					Show: true,
				},
				XValues: timestamps,
				YValues: rtts,
			},
		},
	}

	output, err := os.Create(args[1])
	if err != nil {
		panic(err)
	}
	defer output.Close()

	err = graph.Render(chart.PNG, output)
	if err != nil {
		panic(err)
	}
}

var graphCmd = &cobra.Command{
	Use:   "graph <input file> <output file>",
	Short: "Generate graph from log file",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		createLineChart(args)
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
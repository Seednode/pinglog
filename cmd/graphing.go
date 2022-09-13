/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"fmt"
	"net/netip"
	"os"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/spf13/cobra"
)

func calculateXAxis(slice []time.Time) []time.Time {
	r := []time.Time{}

	for i := 0; i < len(slice); i++ {
		r = append(r, slice[i])
	}

	return r
}

func parseFile(logFile string) (string, netip.Addr, []time.Time, []time.Duration, error) {
	var hostname string
	var ipaddr netip.Addr
	var timestamps []time.Time
	var rtts []time.Duration

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
			if err != nil {
				return "", netip.Addr{}, nil, nil, err
			}
			rtts = append(rtts, r)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", netip.Addr{}, nil, nil, err
	}

	return hostname, ipaddr, timestamps, rtts, nil
}

// generate random data for line chart
func generateLineItems(data []time.Duration) []opts.LineData {
	items := make([]opts.LineData, 0)

	for i := 0; i < len(data); i++ {
		items = append(items, opts.LineData{Value: data[i]})
	}

	return items
}

func createLineChart(args []string) {
	hostname, ipaddr, timestamps, rtts, err := parseFile(args[0])
	if err != nil {
		panic(err)
	}

	chart := charts.NewLine()

	chart.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeInfographic,
		}),
		charts.WithTitleOpts(opts.Title{
			Title: fmt.Sprintf("Pings to %v (%v)", hostname, ipaddr.String()),
		}),
	)

	chart.SetXAxis(calculateXAxis(timestamps)).
		AddSeries("RTTs", generateLineItems(rtts)).
		SetSeriesOptions((charts.WithLineChartOpts(opts.LineChart{Smooth: true})))

	output, err := os.Create(args[1])
	if err != nil {
		panic(err)
	}

	_ = chart.Render(output)
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

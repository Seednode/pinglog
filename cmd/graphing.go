/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"fmt"
	"net/netip"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/spf13/cobra"
)

type ChartData struct {
	host      netip.Addr
	timestamp time.Time
	rtt       time.Duration
}

func calculateXAxisOld(slice []ChartData) (time.Time, time.Time) {
	s := append([]ChartData{}, slice...)

	sort.Slice(s, func(i, j int) bool {
		return s[i].timestamp.Before(s[j].timestamp)
	})

	rangeStart := s[0].timestamp
	rangeEnd := s[len(s)-1].timestamp

	return rangeStart, rangeEnd
}

func calculateXAxis(slice []ChartData) []time.Time {
	r := []time.Time{}

	for i := 0; i < len(slice); i++ {
		r = append(r, slice[i].timestamp)
	}

	return r
}

// Receives log file, parses out values, and returns slice of ChartData
func parseFile(logFile string) ([]ChartData, error) {
	results := []ChartData{}

	file, err := os.Open(logFile)
	if err != nil {
		return nil, err
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
		if strings.Contains(scanner.Text(), "|") && !strings.Contains(scanner.Text(), "lost or arrived out of order") {
			data := ChartData{}

			fields := strings.Fields(stripansi.Strip(scanner.Text()))

			data.host, err = netip.ParseAddr(strings.TrimSuffix(fields[7], ":"))
			if err != nil {
				return nil, err
			}

			t := fields[0] + " " + fields[1] + " " + fields[2]
			data.timestamp, err = time.Parse(DATE, t)
			if err != nil {
				return nil, err
			}

			data.rtt, err = time.ParseDuration(strings.TrimPrefix(fields[10], "time="))
			if err != nil {
				return nil, err
			}

			results = append(results, data)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// generate random data for line chart
func generateLineItems(data []ChartData) []opts.LineData {
	items := make([]opts.LineData, 0)

	for i := 0; i < len(data); i++ {
		items = append(items, opts.LineData{Value: data[i].rtt})
	}

	return items
}

func createLineChart(args []string) {
	parsed := []ChartData{}

	for i := 0; i < len(args); i++ {
		results, err := parseFile(args[i])
		if err != nil {
			panic(err)
			// return err
		}
		parsed = append(parsed, results...)
	}

	for i := 0; i < len(parsed); i++ {
		fmt.Println(parsed[i].rtt)
	}

	chart := charts.NewLine()

	chart.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeInfographic,
		}),
		charts.WithTitleOpts(opts.Title{
			Title:    "I'm not sure I'm doing this right",
			Subtitle: "But I'll try anyway",
		}),
	)

	chart.SetXAxis(calculateXAxis(parsed)).
		AddSeries("RTTs", generateLineItems(parsed)).
		SetSeriesOptions((charts.WithLineChartOpts(opts.LineChart{Smooth: true})))

	f, _ := os.Create("chart.html")
	_ = chart.Render(f)

	//return nil
}

var graphCmd = &cobra.Command{
	Use:   "graph <file1> [file2] .. [fileN]",
	Short: "Generate graph from log file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createLineChart(args)
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}

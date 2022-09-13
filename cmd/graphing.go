/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/netip"
	"os"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/spf13/cobra"
)

type ChartData struct {
	host      netip.Addr
	timestamp time.Time
	rtt       time.Duration
}

// generate random data for bar chart
func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < 6; i++ {
		items = append(items, opts.BarData{Value: rand.Intn(500)})
	}
	return items
}

func createBarChart() {
	// create a new bar instance
	bar := charts.NewBar()

	// Set global options
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Bar chart in Go",
		Subtitle: "This is fun to use!",
	}))

	// Put data into instance
	bar.SetXAxis([]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}).
		AddSeries("Category A", generateBarItems()).
		AddSeries("Category B", generateBarItems())
	f, _ := os.Create("bar.html")
	_ = bar.Render(f)
}

// Receives log file and pointer to slice, parses out values from log file, and returns slice of ChartData
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
		if strings.Contains(scanner.Text(), "|") {
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

func createScatterChart(args []string) {
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

	fmt.Println("Exiting")
	os.Exit(0)

	chart := charts.NewScatter()

	chart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "I'm not sure I'm doing this right",
		Subtitle: "But I'll try anyway",
	}))

	//return nil
}

var graphCmd = &cobra.Command{
	Use:   "graph <file1> [file2] .. [fileN]",
	Short: "Generate graph from log file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createScatterChart(args)
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}

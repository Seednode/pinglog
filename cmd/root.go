/*
Copyright © 2022 Seednode <seednode@seedno.de>

*/
package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

var Version string = "0.1"

var Color bool
var Count int
var Interval time.Duration
var Size int
var Timeout time.Duration
var Quiet bool
var Privileged bool
var Timestamp bool
var Dropped bool
var NoRTT bool
var Output string

var rootCmd = &cobra.Command{
	Use:   "pinglog [flags] <host>",
	Short: "A more featureful ping tool.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pingCmd(args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&Color, "color", false, "colorize output")
	rootCmd.Flags().IntVarP(&Count, "count", "c", -1, "number of packets to send")
	rootCmd.Flags().DurationVarP(&Interval, "interval", "i", time.Second, "time between packets")
	rootCmd.Flags().IntVarP(&Size, "size", "s", 56, "size of packets, in bytes")
	rootCmd.Flags().DurationVarP(&Timeout, "timeout", "t", time.Minute*15, "connection timeout")
	rootCmd.Flags().BoolVarP(&Quiet, "quiet", "q", false, "only display summary at end")
	rootCmd.Flags().BoolVar(&Privileged, "privileged", false, "run as privileged user (needed on Windows)")
	rootCmd.Flags().BoolVar(&Timestamp, "timestamp", false, "prepend timestamps to output")
	rootCmd.Flags().BoolVar(&Dropped, "dropped", false, "log dropped packets")
	rootCmd.Flags().BoolVar(&NoRTT, "no-rtt", false, "do not record RTTs (reduces memory use for long sessions)")
	rootCmd.Flags().StringVarP(&Output, "output", "o", "", "write to the specified file as well as stdout")
	rootCmd.Flags().Lookup("output").NoOptDefVal = "<hostname>.log"
}

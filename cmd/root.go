/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"math"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var Count int
var Dropped bool
var ForceOverwrite bool
var Interval time.Duration
var Color bool
var NoRTT bool
var Output string
var Privileged bool
var Quiet bool
var Size int
var Timeout time.Duration
var Timestamp bool

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
	rootCmd.Flags().IntVarP(&Count, "count", "c", -1, "number of packets to send")
	rootCmd.Flags().BoolVarP(&Color, "color", "C", false, "enable colorized output")
	rootCmd.Flags().BoolVarP(&Dropped, "dropped", "d", false, "log dropped packets")
	rootCmd.Flags().BoolVarP(&ForceOverwrite, "force", "f", false, "overwrite log file without prompting")
	rootCmd.Flags().DurationVarP(&Interval, "interval", "i", time.Second, "time between packets")
	rootCmd.Flags().BoolVarP(&NoRTT, "no-rtt", "n", false, "do not record RTTs (reduces memory use for long sessions)")
	rootCmd.Flags().StringVarP(&Output, "output", "o", "", "write to the specified file as well as stdout")
	rootCmd.Flags().BoolVarP(&Privileged, "privileged", "p", false, "run as privileged user (always enabled on Windows)")
	rootCmd.Flags().BoolVarP(&Quiet, "quiet", "q", false, "only display summary at end")
	rootCmd.Flags().IntVarP(&Size, "size", "s", 56, "size of packets, in bytes")
	rootCmd.Flags().DurationVarP(&Timeout, "timeout", "w", time.Duration(math.MaxInt64), "connection timeout")
	rootCmd.Flags().BoolVarP(&Timestamp, "timestamp", "t", false, "prepend timestamps to output")

	rootCmd.Flags().Lookup("output").NoOptDefVal = "<hostname>.log"
}

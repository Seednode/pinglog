/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"log"
	"math"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var Beep bool
var Color bool
var Count uint64
var Dropped bool
var ForceOverwrite bool
var Interval time.Duration
var IPv4 bool
var IPv6 bool
var MaxRTT time.Duration
var Output string
var Quiet bool
var Size uint16
var Timeout time.Duration
var Timestamp bool
var TTL uint16

var rootCmd = &cobra.Command{
	Use:   "pinglog [flags] <host>",
	Short: "A more featureful ping tool.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := pingCmd(args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&Beep, "beep", "b", true, "enable audible bell for exceeded max-rtt")
	rootCmd.Flags().BoolVarP(&Color, "color", "C", true, "enable colorized output")
	rootCmd.Flags().Uint64VarP(&Count, "count", "c", 0, "number of pings to send")
	rootCmd.Flags().BoolVarP(&Dropped, "dropped", "d", true, "log dropped pings")
	rootCmd.Flags().BoolVarP(&ForceOverwrite, "force", "f", false, "overwrite log file without prompting")
	rootCmd.Flags().DurationVarP(&Interval, "interval", "i", time.Second, "time between pings")
	rootCmd.Flags().BoolVarP(&IPv4, "ipv4", "4", false, "force dns resolution to ipv4")
	rootCmd.Flags().BoolVarP(&IPv6, "ipv6", "6", false, "force dns resolution to ipv6")
	rootCmd.MarkFlagsMutuallyExclusive("ipv4", "ipv6")
	rootCmd.Flags().DurationVarP(&MaxRTT, "max-rtt", "m", time.Hour, "colorize pings over this rtt")
	rootCmd.Flags().StringVarP(&Output, "output", "o", "", "write to the specified file as well as stdout")
	rootCmd.Flags().Lookup("output").NoOptDefVal = "<hostname>.log"
	rootCmd.Flags().BoolVarP(&Quiet, "quiet", "q", false, "only display summary at end")
	rootCmd.Flags().Uint16VarP(&Size, "size", "s", 56, "size of payload, in bytes")
	rootCmd.Flags().DurationVarP(&Timeout, "timeout", "w", time.Duration(math.MaxInt64), "connection timeout")
	rootCmd.Flags().BoolVarP(&Timestamp, "timestamp", "t", true, "prepend timestamps to output")
	rootCmd.Flags().Uint16VarP(&TTL, "ttl", "T", 128, "maximum time-to-live")
}

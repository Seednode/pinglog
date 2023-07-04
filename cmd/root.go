/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	Version string = "0.20.0"
)

var beep bool
var colorize bool
var count uint64
var dropped bool
var forceOverwrite bool
var interval time.Duration
var ipv4 bool
var ipv6 bool
var maxRTT time.Duration
var output string
var quiet bool
var size uint16
var timeout time.Duration
var timestamp bool
var ttl uint16
var version bool

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
	rootCmd.Flags().BoolVarP(&beep, "beep", "b", true, "enable audible bell for exceeded max-rtt")
	rootCmd.Flags().BoolVarP(&colorize, "colorize", "C", true, "enable colorized output")
	rootCmd.Flags().Uint64VarP(&count, "count", "c", 0, "number of pings to send")
	rootCmd.Flags().BoolVarP(&dropped, "dropped", "d", true, "log dropped pings")
	rootCmd.Flags().BoolVarP(&forceOverwrite, "force", "f", false, "overwrite log file without prompting")
	rootCmd.Flags().DurationVarP(&interval, "interval", "i", time.Second, "time between pings")
	rootCmd.Flags().BoolVarP(&ipv4, "ipv4", "4", false, "force dns resolution to ipv4")
	rootCmd.Flags().BoolVarP(&ipv6, "ipv6", "6", false, "force dns resolution to ipv6")
	rootCmd.MarkFlagsMutuallyExclusive("ipv4", "ipv6")
	rootCmd.Flags().DurationVarP(&maxRTT, "max-rtt", "m", time.Hour, "colorize pings over this rtt")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "write to the specified file as well as stdout")
	rootCmd.Flags().Lookup("output").NoOptDefVal = "<hostname>.log"
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "only display summary at end")
	rootCmd.Flags().Uint16VarP(&size, "size", "s", 56, "size of payload, in bytes")
	rootCmd.Flags().DurationVarP(&timeout, "timeout", "w", 10*time.Second, "connection timeout")
	rootCmd.Flags().BoolVarP(&timestamp, "timestamp", "t", true, "prepend timestamps to output")
	rootCmd.Flags().Uint16VarP(&ttl, "ttl", "T", 128, "maximum time-to-live")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "display version and exit")

	rootCmd.SetVersionTemplate("pinglog v{{.Version}}\n")
	rootCmd.Version = Version
}

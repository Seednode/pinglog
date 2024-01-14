/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/spf13/cobra"
)

const (
	ReleaseVersion string = "0.24.0"
)

var (
	ErrInvalidCount = errors.New("count must be a positive integer")
	ErrInvalidSize  = errors.New("size must be a positive integer between 1 and 65527 bytes inclusive")
	ErrInvalidTtl   = errors.New("ttl must be a positive integer no higher than 255")
)

var beep bool
var colorize bool
var count int
var dropped bool
var interval time.Duration
var ipv4 bool
var ipv6 bool
var maxRtt time.Duration
var quiet bool
var size int
var timeout time.Duration
var timestamp bool
var ttl int
var version bool

var rootCmd = &cobra.Command{
	Use:   "pinglog [flags] <host>",
	Short: "A more featureful ping tool.",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		switch {
		case count < 0:
			return ErrInvalidCount
		case size < 1 || size > 65527:
			return ErrInvalidSize
		case ttl < 1 || ttl > 255:
			return ErrInvalidTtl
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := pingCmd(args)
		if err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&beep, "beep", "b", true, "enable audible bell for exceeded max-rtt")
	rootCmd.Flags().BoolVarP(&colorize, "color", "C", true, "enable colorized output")
	rootCmd.Flags().IntVarP(&count, "count", "c", 0, "number of pings to send")
	rootCmd.Flags().BoolVarP(&dropped, "dropped", "d", true, "log dropped pings")
	rootCmd.Flags().DurationVarP(&interval, "interval", "i", time.Second, "time between pings")
	rootCmd.Flags().BoolVarP(&ipv4, "ipv4", "4", false, "force dns resolution to ipv4")
	rootCmd.Flags().BoolVarP(&ipv6, "ipv6", "6", false, "force dns resolution to ipv6")
	rootCmd.MarkFlagsMutuallyExclusive("ipv4", "ipv6")
	rootCmd.Flags().DurationVarP(&maxRtt, "max-rtt", "m", time.Hour, "colorize pings over this rtt")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "only display summary at end")
	rootCmd.Flags().IntVarP(&size, "size", "s", 56, "size of payload, in bytes")
	rootCmd.Flags().DurationVarP(&timeout, "timeout", "w", time.Duration(math.MaxInt64), "timeout before ping exits, regardless of number of packets sent or received")
	rootCmd.Flags().BoolVarP(&timestamp, "timestamp", "t", true, "prepend timestamps to output")
	rootCmd.Flags().IntVarP(&ttl, "ttl", "T", 128, "maximum time-to-live")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "display version and exit")

	rootCmd.Flags().SetInterspersed(true)

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.SilenceErrors = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	rootCmd.SetVersionTemplate("pinglog v{{.Version}}\n")
	rootCmd.Version = ReleaseVersion
}

/*
Copyright © 2025 Seednode <seednode@seedno.de>
*/

package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ReleaseVersion string = "1.2.0"
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

func main() {
	cmd := &cobra.Command{
		Use:   "pinglog [flags] <host>",
		Short: "A more featureful ping tool.",
		Args:  cobra.ExactArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initializeConfig(cmd)
		},
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

	lossCmd := &cobra.Command{
		Use:   "loss <file1> [file2]...",
		Short: "Calculate periods of packet loss from log file(s)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for file := range args {
				err := calculateLoss(args[file])
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.AddCommand(lossCmd)

	var stripCmd = &cobra.Command{
		Use:   "strip <file>",
		Short: "Strip ANSI color codes from log file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := StripColors(args)
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(stripCmd)

	cmd.Flags().BoolVarP(&beep, "beep", "b", true, "enable audible bell for exceeded max-rtt")
	cmd.Flags().BoolVarP(&colorize, "color", "C", true, "enable colorized output")
	cmd.Flags().IntVarP(&count, "count", "c", 0, "number of pings to send")
	cmd.Flags().BoolVarP(&dropped, "dropped", "d", true, "log dropped pings")
	cmd.Flags().DurationVarP(&interval, "interval", "i", time.Second, "time between pings")
	cmd.Flags().BoolVarP(&ipv4, "ipv4", "4", false, "force dns resolution to ipv4")
	cmd.Flags().BoolVarP(&ipv6, "ipv6", "6", false, "force dns resolution to ipv6")
	cmd.MarkFlagsMutuallyExclusive("ipv4", "ipv6")
	cmd.Flags().DurationVarP(&maxRtt, "max-rtt", "m", time.Hour, "colorize pings over this rtt")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "only display summary at end")
	cmd.Flags().IntVarP(&size, "size", "s", 56, "size of payload, in bytes")
	cmd.Flags().DurationVarP(&timeout, "timeout", "w", time.Duration(math.MaxInt64), "timeout before ping exits, regardless of number of packets sent or received")
	cmd.Flags().BoolVarP(&timestamp, "timestamp", "t", true, "prepend timestamps to output")
	cmd.Flags().IntVarP(&ttl, "ttl", "T", 128, "maximum time-to-live")
	cmd.Flags().BoolVarP(&version, "version", "V", false, "display version and exit")

	cmd.CompletionOptions.HiddenDefaultCmd = true

	cmd.Flags().SetInterspersed(true)

	cmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	cmd.SetVersionTemplate("pinglog v{{.Version}}\n")

	cmd.SilenceErrors = true

	cmd.Version = ReleaseVersion

	log.SetFlags(0)

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func initializeConfig(cmd *cobra.Command) {
	v := viper.New()

	v.SetEnvPrefix("pinglog")

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	v.AutomaticEnv()

	bindFlags(cmd, v)
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", "_")

		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

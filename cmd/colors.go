package cmd

import (
	"fmt"

	"github.com/fatih/color"
)

func HighlightPacketLoss(packetLoss float64) string {
	red := color.New(color.FgRed).Add(color.Bold)

	if packetLoss != 0.0 {
		return red.Sprintf("%.3f%% packet loss", packetLoss)
	}

	return fmt.Sprintf("%.3f%% packet loss", packetLoss)
}

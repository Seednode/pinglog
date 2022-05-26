package cmd

import (
	"github.com/fatih/color"
)

func HighlightPacketLoss(packetLoss float64) string {
	blue := color.New(color.FgBlue)
	red := color.New(color.FgRed).Add(color.Bold)

	if packetLoss != 0.0 {
		return red.Sprintf("%.3f%%", packetLoss)
	}

	return blue.Sprintf("%.3f%%", packetLoss)
}

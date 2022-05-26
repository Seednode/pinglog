package cmd

import (
	"fmt"

	"github.com/fatih/color"
)

func Yellow() func(a ...interface{}) string {
	return color.New(color.FgYellow).SprintFunc()
}

func HighlightPacketLoss(packetLoss float64) string {
	if packetLoss != 0.0 {
		return color.RedString(fmt.Sprintf("%.3f%%", packetLoss))
	}

	return fmt.Sprintf("%.3f%%", packetLoss)
}

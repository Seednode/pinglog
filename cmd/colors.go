/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Colors struct {
	Blue  *color.Color
	Green *color.Color
	Grey  *color.Color
	Red   *color.Color
}

func highlightPacketLoss(packetLoss float64, colors *Colors) string {
	if packetLoss != 0.0 {
		return colors.Red.Sprintf("%.3f%%", packetLoss)
	} else {
		return colors.Blue.Sprintf("%.3f%%", packetLoss)
	}
}

func highlightLongRTT(packetRTT time.Duration, colors *Colors, isEnding bool) string {
	switch {
	case packetRTT > maxRtt && beep && !isEnding:
		fmt.Print("\a")
		return colors.Red.Sprintf("%v", packetRTT)
	case packetRTT > maxRtt:
		return colors.Red.Sprintf("%v", packetRTT)
	default:
		return colors.Blue.Sprintf("%v", packetRTT)
	}
}

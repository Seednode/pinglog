/*
Copyright © 2025 Seednode <seednode@seedno.de>
*/

package main

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
		fmt.Println("\a")

		return colors.Red.Sprintf("%s", packetRTT)
	case packetRTT > maxRtt:
		return colors.Red.Sprintf("%s", packetRTT)
	default:
		return colors.Blue.Sprintf("%s", packetRTT)
	}
}

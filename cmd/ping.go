/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
)

const DATE string = "2006-01-02 15:04:05.000 MST"

func highlightPacketLoss(packetLoss float64) string {
	if packetLoss != 0.0 {
		red := color.New(color.FgRed).Add(color.Bold)

		return red.Sprintf("%.3f%%", packetLoss)
	}

	blue := color.New(color.FgBlue)

	return blue.Sprintf("%.3f%%", packetLoss)
}

func highlightLongRTT(packetRTT time.Duration, isEnding bool) string {
	blue := color.New(color.FgBlue)
	red := color.New(color.FgRed)

	if packetRTT > MaxRTT {
		if Beep && !isEnding {
			fmt.Print("\a")
		}

		return red.Sprintf("%v", packetRTT)
	} else {
		return blue.Sprintf("%v", packetRTT)
	}
}

func humanReadableSize(bytes int) string {
	const unit = 1000

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0

	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB",
		float64(bytes)/float64(div), "kMGTPE"[exp])
}

func pingCmd(arguments []string) {
	host := arguments[0]

	var wasInterrupted = false

	myPing, err := ping.NewPinger(host)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			<-c
			wasInterrupted = true
			myPing.Stop()
		}
	}()

	myPing.Count = Count
	myPing.Size = Size
	myPing.Interval = Interval
	myPing.Timeout = Timeout
	myPing.TTL = TTL
	myPing.RecordRtts = RTT

	// privileged is required on Windows
	switch runtime.GOOS {
	case "windows":
		myPing.SetPrivileged(true)
	default:
		myPing.SetPrivileged(false)
	}

	switch {
	case IPv4:
		myPing.SetNetwork("ip4")
	case IPv6:
		myPing.SetNetwork("ip6")
	default:
		myPing.SetNetwork("ip")
	}
	myPing.Resolve()

	switch Color {
	case true:
		color.NoColor = false
	default:
		color.NoColor = true
	}

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	grey := color.New(color.FgHiBlack)
	red := color.New(color.FgRed)

	var expectedPacket = 0
	var currentPacket = 0

	myPing.OnRecv = func(pkt *ping.Packet) {
		currentPacket = pkt.Seq

		if Dropped && Timestamp && (expectedPacket != currentPacket) {
			for c := expectedPacket; c < currentPacket; c++ {
				timeStamp := time.Now().Format(DATE)

				fmt.Printf("%v | %v", grey.Sprintf(timeStamp), red.Sprintf("Packet %v lost or arrived out of order.\n", c))
			}
			expectedPacket = currentPacket + 1
		} else if Dropped && (expectedPacket != currentPacket) {
			for c := expectedPacket; c < currentPacket; c++ {
				fmt.Printf("%v", red.Sprintf("Packet %v lost or arrived out of order.\n", c))
			}
			expectedPacket = currentPacket + 1
		} else if Dropped {
			expectedPacket = currentPacket + 1
		}

		if Quiet {
			return
		}

		if Timestamp {
			timeStamp := time.Now().Format(DATE)
			fmt.Printf("%v | %v from %v: icmp_seq=%v ttl=%v time=%v\n",
				grey.Sprintf(timeStamp),
				blue.Sprintf("%v bytes", pkt.Nbytes-8),
				blue.Sprintf("%v", pkt.IPAddr),
				blue.Sprintf("%v", pkt.Seq),
				blue.Sprintf("%v", pkt.Ttl),
				highlightLongRTT(pkt.Rtt.Truncate(time.Microsecond), false))
		} else {
			fmt.Printf("%v from %v: icmp_seq=%v ttl=%v time=%v\n",
				blue.Sprintf("%v bytes", pkt.Nbytes-8),
				blue.Sprintf("%v", pkt.IPAddr),
				blue.Sprintf("%v", pkt.Seq),
				blue.Sprintf("%v", pkt.Ttl),
				highlightLongRTT(pkt.Rtt.Truncate(time.Microsecond), false))
		}

		if currentPacket == (Count - 1) {
			myPing.Stop()
		}
	}

	myPing.OnDuplicateRecv = func(pkt *ping.Packet) {
		if Timestamp {
			timeStamp := time.Now().Format(DATE)

			fmt.Printf("%v | %v from %v: icmp_seq=%v ttl=%v time=%v %v\n",
				grey.Sprintf(timeStamp),
				blue.Sprintf("%v bytes", pkt.Nbytes-8),
				blue.Sprintf("%v", pkt.IPAddr),
				blue.Sprintf("%v", pkt.Seq),
				blue.Sprintf("%v", pkt.Ttl),
				highlightLongRTT(pkt.Rtt.Truncate(time.Microsecond), false),
				red.Sprintf("(DUP!)"))

			return
		}

		fmt.Printf("%v from %v: icmp_seq=%v ttl=%v time=%v %v\n",
			blue.Sprintf("%v bytes", pkt.Nbytes-8),
			blue.Sprintf("%v", pkt.IPAddr),
			blue.Sprintf("%v", pkt.Seq),
			blue.Sprintf("%v", pkt.Ttl),
			highlightLongRTT(pkt.Rtt.Truncate(time.Microsecond), false),
			red.Sprintf("(DUP!)"))
	}

	var startTime time.Time

	myPing.OnFinish = func(stats *ping.Statistics) {
		runTime := time.Since(startTime)

		if !wasInterrupted && Dropped && (Count != -1) && (currentPacket != (Count - 1)) {
			for c := currentPacket + 1; c < Count; c++ {
				fmt.Printf("%v", red.Sprintf("Packet %v lost or arrived out of order.\n", c))
			}
		}

		fmt.Printf("\n--- %v ping statistics ---\n", green.Sprintf(stats.Addr))

		fmt.Printf("%v packets transmitted (%v), %v received (%v), %v packet loss, time %v\n",
			blue.Sprintf("%v", stats.PacketsSent),
			blue.Sprintf(humanReadableSize(stats.PacketsSent*myPing.Size)),
			blue.Sprintf("%v", stats.PacketsRecv),
			blue.Sprintf(humanReadableSize(stats.PacketsRecv*myPing.Size)),
			highlightPacketLoss(stats.PacketLoss),
			blue.Sprintf("%v", runTime.Truncate(time.Millisecond)))

		fmt.Printf("rtt min/avg/max/mdev = %v/%v/%v/%v\n\n",
			highlightLongRTT(stats.MinRtt.Truncate(time.Microsecond), true),
			highlightLongRTT(stats.AvgRtt.Truncate(time.Microsecond), true),
			highlightLongRTT(stats.MaxRtt.Truncate(time.Microsecond), true),
			blue.Sprintf("%v", stats.StdDevRtt.Truncate(time.Microsecond)))
	}

	if Output == "<hostname>.log" {
		endLogging, err := logOutput(host + ".log")
		if err != nil {
			log.Fatal(err)
		}
		defer endLogging()
	} else if Output != "" {
		endLogging, err := logOutput(Output)
		if err != nil {
			log.Fatal(err)
		}
		defer endLogging()
	}

	startTime = time.Now()

	fmt.Printf("PING %v (%v) %v(%v) bytes of data.\n",
		green.Sprintf("%v", myPing.Addr()),
		blue.Sprintf("%v", myPing.IPAddr()),
		blue.Sprintf("%v", Size),
		blue.Sprintf("%v", Size+28))

	err = myPing.Run()
	if err != nil {
		log.Fatal(err)
	}
}

/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
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

type Colors struct {
	Blue  *color.Color
	Green *color.Color
	Grey  *color.Color
	Red   *color.Color
}

type Packets struct {
	Expected int
	Current  int
}

func initializeColors() *Colors {
	return &Colors{
		Blue:  color.New(color.FgBlue),
		Green: color.New(color.FgGreen),
		Grey:  color.New(color.FgHiBlack),
		Red:   color.New(color.FgRed),
	}
}

func initializeCounters() *Packets {
	return &Packets{
		Expected: 0,
		Current:  0,
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

func highlightPacketLoss(packetLoss float64, colors *Colors) string {
	if packetLoss != 0.0 {
		return colors.Red.Sprintf("%.3f%%", packetLoss)
	} else {
		return colors.Blue.Sprintf("%.3f%%", packetLoss)
	}
}

func highlightLongRTT(packetRTT time.Duration, colors *Colors, isEnding bool) string {
	switch {
	case packetRTT > MaxRTT && Beep && !isEnding:
		fmt.Print("\a")
		return colors.Red.Sprintf("%v", packetRTT)
	case packetRTT > MaxRTT:
		return colors.Red.Sprintf("%v", packetRTT)
	default:
		return colors.Blue.Sprintf("%v", packetRTT)
	}
}

func configurePinger(myPing *ping.Pinger) error {
	myPing.Count = int(Count)
	myPing.Size = int(Size)
	myPing.Interval = Interval
	myPing.Timeout = Timeout
	myPing.TTL = int(TTL)
	myPing.RecordRtts = false

	// Running in privileged mode is required on Windows hosts
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
	err := myPing.Resolve()
	if err != nil {
		return err
	}

	color.NoColor = !Color

	return nil
}

func showReceived(pkt *ping.Packet, myPing *ping.Pinger, packets *Packets, colors *Colors) error {
	packets.Current = pkt.Seq

	switch {
	case Dropped && Timestamp && (packets.Expected != packets.Current):
		for c := packets.Expected; c < packets.Current; c++ {
			timeStamp := time.Now().Format(DATE)

			_, err := fmt.Printf("%v | %v", colors.Grey.Sprintf(timeStamp), colors.Red.Sprintf("Packet %v lost or arrived out of order.\n", c))
			if err != nil {
				return err
			}
		}
		packets.Expected = packets.Current + 1
	case Dropped && (packets.Expected != packets.Current):
		for c := packets.Expected; c < packets.Current; c++ {
			_, err := fmt.Printf("%v", colors.Red.Sprintf("Packet %v lost or arrived out of order.\n", c))
			if err != nil {
				return err
			}
		}
		packets.Expected = packets.Current + 1
	case Dropped:
		packets.Expected = packets.Current + 1
	}

	if Quiet {
		return nil
	}

	if Timestamp {
		timeStamp := time.Now().Format(DATE)
		_, err := fmt.Printf("%v | %v from %v: icmp_seq=%v ttl=%v time=%v\n",
			colors.Grey.Sprintf(timeStamp),
			colors.Blue.Sprintf("%v bytes", pkt.Nbytes-8),
			colors.Blue.Sprintf("%v", pkt.IPAddr),
			colors.Blue.Sprintf("%v", pkt.Seq),
			colors.Blue.Sprintf("%v", pkt.Ttl),
			highlightLongRTT(pkt.Rtt.Truncate(time.Microsecond), colors, false))
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Printf("%v from %v: icmp_seq=%v ttl=%v time=%v\n",
			colors.Blue.Sprintf("%v bytes", pkt.Nbytes-8),
			colors.Blue.Sprintf("%v", pkt.IPAddr),
			colors.Blue.Sprintf("%v", pkt.Seq),
			colors.Blue.Sprintf("%v", pkt.Ttl),
			highlightLongRTT(pkt.Rtt.Truncate(time.Microsecond), colors, false))
		if err != nil {
			return err
		}
	}

	if packets.Current == (int(Count) - 1) {
		myPing.Stop()
	}

	return nil
}

func showDuplicate(pkt *ping.Packet, colors *Colors) error {
	if Timestamp {
		timeStamp := time.Now().Format(DATE)

		_, err := fmt.Printf("%v | %v from %v: icmp_seq=%v ttl=%v time=%v %v\n",
			colors.Grey.Sprintf(timeStamp),
			colors.Blue.Sprintf("%v bytes", pkt.Nbytes-8),
			colors.Blue.Sprintf("%v", pkt.IPAddr),
			colors.Blue.Sprintf("%v", pkt.Seq),
			colors.Blue.Sprintf("%v", pkt.Ttl),
			highlightLongRTT(pkt.Rtt.Truncate(time.Microsecond), colors, false),
			colors.Red.Sprintf("(DUP!)"))
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Printf("%v from %v: icmp_seq=%v ttl=%v time=%v %v\n",
			colors.Blue.Sprintf("%v bytes", pkt.Nbytes-8),
			colors.Blue.Sprintf("%v", pkt.IPAddr),
			colors.Blue.Sprintf("%v", pkt.Seq),
			colors.Blue.Sprintf("%v", pkt.Ttl),
			highlightLongRTT(pkt.Rtt.Truncate(time.Microsecond), colors, false),
			colors.Red.Sprintf("(DUP!)"))
		if err != nil {
			return err
		}
	}

	return nil
}

func showStatistics(stats *ping.Statistics, myPing *ping.Pinger, packets *Packets, colors *Colors, startTime time.Time, wasInterrupted bool, isEnding bool) string {
	var s string

	runTime := time.Since(startTime)

	if isEnding && !wasInterrupted && Dropped && (Count != 0) && (packets.Current != (int(Count) - 1)) {
		for c := packets.Current + 1; c < int(Count); c++ {
			s += fmt.Sprintf("%v", colors.Red.Sprintf("Packet %v lost or arrived out of order.\n", c))
		}
	}

	s += fmt.Sprintf("--- %v ping statistics ---\n", colors.Green.Sprintf(stats.Addr))

	s += fmt.Sprintf("%v packets transmitted (%v), %v received (%v), %v packet loss, time %v\n",
		colors.Blue.Sprintf("%v", stats.PacketsSent),
		colors.Blue.Sprintf(humanReadableSize(stats.PacketsSent*myPing.Size)),
		colors.Blue.Sprintf("%v", stats.PacketsRecv),
		colors.Blue.Sprintf(humanReadableSize(stats.PacketsRecv*myPing.Size)),
		highlightPacketLoss(stats.PacketLoss, colors),
		colors.Blue.Sprintf("%v", runTime.Truncate(time.Millisecond)))

	s += fmt.Sprintf("rtt min/avg/max/mdev = %v/%v/%v/%v\n\n",
		highlightLongRTT(stats.MinRtt.Truncate(time.Microsecond), colors, true),
		highlightLongRTT(stats.AvgRtt.Truncate(time.Microsecond), colors, true),
		highlightLongRTT(stats.MaxRtt.Truncate(time.Microsecond), colors, true),
		colors.Blue.Sprintf("%v", stats.StdDevRtt.Truncate(time.Microsecond)))

	return s
}

func showStart(myPing *ping.Pinger, colors *Colors) error {
	_, err := fmt.Printf("PING %v (%v) %v(%v) bytes of data.\n",
		colors.Green.Sprintf("%v", myPing.Addr()),
		colors.Blue.Sprintf("%v", myPing.IPAddr()),
		colors.Blue.Sprintf("%v", Size),
		colors.Blue.Sprintf("%v", Size+28))
	if err != nil {
		return err
	}

	return nil
}

func pingCmd(arguments []string) error {
	host := arguments[0]

	var startTime = time.Time{}
	var wasInterrupted = false

	myPing, err := ping.NewPinger(host)
	if err != nil {
		return err
	}

	err = configurePinger(myPing)
	if err != nil {
		return err
	}

	colors := initializeColors()
	packets := initializeCounters()

	myPing.OnRecv = func(pkt *ping.Packet) {
		err := showReceived(pkt, myPing, packets, colors)
		if err != nil {
			log.Fatal(err)
		}
	}

	myPing.OnDuplicateRecv = func(pkt *ping.Packet) {
		err := showDuplicate(pkt, colors)
		if err != nil {
			log.Fatal(err)
		}
	}

	myPing.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n%v", showStatistics(stats, myPing, packets, colors, startTime, wasInterrupted, true))
	}

	switch {
	case Output == "<hostname>.log":
		endLogging, err := logOutput(host + ".log")
		if err != nil {
			return err
		}
		defer endLogging()
	case Output != "":
		endLogging, err := logOutput(Output)
		if err != nil {
			return err
		}
		defer endLogging()
	}

	showStart(myPing, colors)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			<-c
			wasInterrupted = true
			myPing.Stop()
		}
	}()

	go func() {
		consoleReader := bufio.NewReaderSize(os.Stdin, 1)
		for {
			input, _, err := consoleReader.ReadRune()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}

			if string(input) == "\n" {
				fmt.Fprint(os.Stderr, showStatistics(myPing.Statistics(), myPing, packets, colors, startTime, false, false))
			}
		}
	}()

	startTime = time.Now()

	err = myPing.Run()
	if err != nil {
		return err
	}

	return nil
}

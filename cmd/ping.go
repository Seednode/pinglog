/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
)

const DATE string = "2006-01-02 15:04:05.000 MST"

type Packets struct {
	Expected int
	Current  int
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

func configurePinger(pinger *ping.Pinger) error {
	pinger.Count = count
	pinger.Size = size
	pinger.Interval = interval
	pinger.Timeout = timeout
	pinger.TTL = ttl
	pinger.RecordRtts = false

	// Running in privileged mode is required on Windows hosts
	switch runtime.GOOS {
	case "windows":
		pinger.SetPrivileged(true)
	default:
		pinger.SetPrivileged(false)
	}

	switch {
	case ipv4:
		pinger.SetNetwork("ip4")
	case ipv6:
		pinger.SetNetwork("ip6")
	default:
		pinger.SetNetwork("ip")
	}

	err := pinger.Resolve()
	if err != nil {
		return err
	}

	color.NoColor = !colorize

	return nil
}

func showReceived(pkt *ping.Packet, pinger *ping.Pinger, packets *Packets, colors *Colors) error {
	packets.Current = pkt.Seq

	switch {
	case dropped && timestamp && (packets.Expected != packets.Current):
		for c := packets.Expected; c < packets.Current; c++ {
			_, err := fmt.Printf("%v | %v", colors.Grey.Sprintf(time.Now().Format(DATE)), colors.Red.Sprintf("Packet %v lost or arrived out of order.\n", c))
			if err != nil {
				return err
			}
		}
		packets.Expected = packets.Current + 1
	case dropped && (packets.Expected != packets.Current):
		for c := packets.Expected; c < packets.Current; c++ {
			_, err := fmt.Printf("%v", colors.Red.Sprintf("Packet %v lost or arrived out of order.\n", c))
			if err != nil {
				return err
			}
		}
		packets.Expected = packets.Current + 1
	case dropped:
		packets.Expected = packets.Current + 1
	}

	if timestamp && !quiet {
		_, err := fmt.Printf("%v | %v from %v: icmp_seq=%v ttl=%v time=%v\n",
			colors.Grey.Sprintf(time.Now().Format(DATE)),
			colors.Blue.Sprintf("%v bytes", pkt.Nbytes-8),
			colors.Blue.Sprintf("%v", pkt.IPAddr),
			colors.Blue.Sprintf("%v", pkt.Seq),
			colors.Blue.Sprintf("%v", pkt.Ttl),
			highlightLongRTT(pkt.Rtt.Truncate(time.Microsecond), colors, false))
		if err != nil {
			return err
		}
	} else if !quiet {
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

	if packets.Current == (count - 1) {
		pinger.Stop()
	}

	return nil
}

func showDuplicate(pkt *ping.Packet, colors *Colors) error {
	if timestamp {
		_, err := fmt.Printf("%v | %v from %v: icmp_seq=%v ttl=%v time=%v %v\n",
			colors.Grey.Sprintf(time.Now().Format(DATE)),
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

func showStatistics(stats *ping.Statistics, pinger *ping.Pinger, packets *Packets, colors *Colors, startTime time.Time, wasInterrupted bool, isEnding bool) string {
	var s strings.Builder

	if isEnding && !wasInterrupted && dropped && count != 0 && (packets.Current != (count - 1)) {
		for c := packets.Current + 1; c < count; c++ {
			s.WriteString(fmt.Sprintf("%v", colors.Red.Sprintf("Packet %v lost or arrived out of order.\n", c)))
		}
	}

	s.WriteString(fmt.Sprintf("--- %v ping statistics ---\n", colors.Green.Sprintf(stats.Addr)))

	s.WriteString(fmt.Sprintf("%v packets transmitted (%v), %v received (%v), %v packet loss, time %v\n",
		colors.Blue.Sprintf("%v", stats.PacketsSent),
		colors.Blue.Sprintf(humanReadableSize(stats.PacketsSent*pinger.Size)),
		colors.Blue.Sprintf("%v", stats.PacketsRecv),
		colors.Blue.Sprintf(humanReadableSize(stats.PacketsRecv*pinger.Size)),
		highlightPacketLoss(stats.PacketLoss, colors),
		colors.Blue.Sprintf("%v", time.Since(startTime).Truncate(time.Millisecond))))

	s.WriteString(fmt.Sprintf("rtt min/avg/max/mdev = %v/%v/%v/%v\n\n",
		highlightLongRTT(stats.MinRtt.Truncate(time.Microsecond), colors, true),
		highlightLongRTT(stats.AvgRtt.Truncate(time.Microsecond), colors, true),
		highlightLongRTT(stats.MaxRtt.Truncate(time.Microsecond), colors, true),
		colors.Blue.Sprintf("%v", stats.StdDevRtt.Truncate(time.Microsecond))))

	return s.String()
}

func showStart(pinger *ping.Pinger, colors *Colors) error {
	_, err := fmt.Printf("PING %v (%v) %v(%v) bytes of data.\n",
		colors.Green.Sprintf("%v", pinger.Addr()),
		colors.Blue.Sprintf("%v", pinger.IPAddr()),
		colors.Blue.Sprintf("%v", size),
		colors.Blue.Sprintf("%v", size+28))
	if err != nil {
		return err
	}

	return nil
}

func pingCmd(arguments []string) error {
	host := arguments[0]

	var startTime = time.Time{}
	var wasInterrupted = false

	pinger, err := ping.NewPinger(host)
	if err != nil {
		return err
	}

	err = configurePinger(pinger)
	if err != nil {
		return err
	}

	colors := &Colors{
		Blue:  color.New(color.FgBlue),
		Green: color.New(color.FgGreen),
		Grey:  color.New(color.FgHiBlack),
		Red:   color.New(color.FgRed),
	}

	packets := &Packets{
		Expected: 0,
		Current:  0,
	}

	errorChannel := make(chan error)
	done := make(chan bool, 1)

	pinger.OnRecv = func(pkt *ping.Packet) {
		err := showReceived(pkt, pinger, packets, colors)
		if err != nil {
			errorChannel <- err
		}
	}

	pinger.OnDuplicateRecv = func(pkt *ping.Packet) {
		err := showDuplicate(pkt, colors)
		if err != nil {
			errorChannel <- err
		}
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n%v", showStatistics(stats, pinger, packets, colors, startTime, wasInterrupted, true))

		done <- true
	}

	showStart(pinger, colors)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		for {
			<-interrupt
			wasInterrupted = true
			pinger.Stop()
		}
	}()

	go func() {
		consoleReader := bufio.NewReaderSize(os.Stdin, 1)
		for {
			input, _, err := consoleReader.ReadRune()
			if err != nil {
				errorChannel <- err
			}

			if string(input) == "\n" {
				fmt.Fprint(os.Stderr, showStatistics(pinger.Statistics(), pinger, packets, colors, startTime, false, false))
			}
		}
	}()

	startTime = time.Now()

	go func() {
		err = pinger.Run()
		if err != nil {
			errorChannel <- err
		}
	}()

Poll:
	for {
		select {
		case err := <-errorChannel:
			pinger.Stop()

			return err
		case <-done:
			break Poll
		}
	}

	return nil
}

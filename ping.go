/*
Copyright © 2025 Seednode <seednode@seedno.de>
*/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/fatih/color"
	ping "github.com/prometheus-community/pro-bing"
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

	// Running in privileged mode is required to send ICMP pings instead of UDP "pings"
	pinger.SetPrivileged(true)

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
			_, err := fmt.Printf("%s | %s", colors.Grey.Sprint(time.Now().Format(DATE)), colors.Red.Sprintf("Packet %d lost or arrived out of order.\n", c))
			if err != nil {
				return err
			}
		}
		packets.Expected = packets.Current + 1
	case dropped && (packets.Expected != packets.Current):
		for c := packets.Expected; c < packets.Current; c++ {
			colors.Red.Sprintf("Packet %d lost or arrived out of order.\n", c)
		}
		packets.Expected = packets.Current + 1
	case dropped:
		packets.Expected = packets.Current + 1
	}

	if timestamp && !quiet {
		_, err := fmt.Printf("%s | %s from %s: icmp_seq=%s ttl=%s time=%s\n",
			colors.Grey.Sprint(time.Now().Format(DATE)),
			colors.Blue.Sprintf("%d bytes", pkt.Nbytes-8),
			colors.Blue.Sprintf("%s", pkt.IPAddr),
			colors.Blue.Sprintf("%d", pkt.Seq),
			colors.Blue.Sprintf("%d", pkt.TTL),
			highlightLongRTT(pkt.Rtt.Round(time.Microsecond), colors, false))
		if err != nil {
			return err
		}
	} else if !quiet {
		_, err := fmt.Printf("%s from %s: icmp_seq=%s ttl=%s time=%s\n",
			colors.Blue.Sprintf("%d bytes", pkt.Nbytes-8),
			colors.Blue.Sprintf("%s", pkt.IPAddr),
			colors.Blue.Sprintf("%d", pkt.Seq),
			colors.Blue.Sprintf("%d", pkt.TTL),
			highlightLongRTT(pkt.Rtt.Round(time.Microsecond), colors, false))
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
		_, err := fmt.Printf("%s | %s from %s: icmp_seq=%s ttl=%s time=%s %s\n",
			colors.Grey.Sprint(time.Now().Format(DATE)),
			colors.Blue.Sprintf("%d bytes", pkt.Nbytes-8),
			colors.Blue.Sprintf("%s", pkt.IPAddr),
			colors.Blue.Sprintf("%d", pkt.Seq),
			colors.Blue.Sprintf("%d", pkt.TTL),
			highlightLongRTT(pkt.Rtt.Round(time.Microsecond), colors, false),
			colors.Red.Sprintf("(DUP!)"))
		if err != nil {
			return err
		}
	} else {
		_, err := fmt.Printf("%s from %s: icmp_seq=%s ttl=%s time=%s %s\n",
			colors.Blue.Sprintf("%d bytes", pkt.Nbytes-8),
			colors.Blue.Sprintf("%s", pkt.IPAddr),
			colors.Blue.Sprintf("%d", pkt.Seq),
			colors.Blue.Sprintf("%d", pkt.TTL),
			highlightLongRTT(pkt.Rtt.Round(time.Microsecond), colors, false),
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
			s.WriteString(colors.Red.Sprintf("Packet %d lost or arrived out of order.\n", c))
		}
	}

	s.WriteString(fmt.Sprintf("--- %v ping statistics ---\n", colors.Green.Sprint(stats.Addr)))

	s.WriteString(fmt.Sprintf("%s packets transmitted (%s), %s packets received (%s), %s packet loss, time %s\n",
		colors.Blue.Sprintf("%d", stats.PacketsSent),
		colors.Blue.Sprint(humanReadableSize(stats.PacketsSent*pinger.Size)),
		colors.Blue.Sprintf("%d", stats.PacketsRecv),
		colors.Blue.Sprint(humanReadableSize(stats.PacketsRecv*pinger.Size)),
		highlightPacketLoss(stats.PacketLoss, colors),
		colors.Blue.Sprintf("%s", time.Since(startTime).Round(time.Millisecond))))

	s.WriteString(fmt.Sprintf("round-trip min/avg/max/stddev = %s/%s/%s/%s\n\n",
		highlightLongRTT(stats.MinRtt.Round(time.Microsecond), colors, true),
		highlightLongRTT(stats.AvgRtt.Round(time.Microsecond), colors, true),
		highlightLongRTT(stats.MaxRtt.Round(time.Microsecond), colors, true),
		colors.Blue.Sprintf("%v", stats.StdDevRtt.Round(time.Microsecond))))

	return s.String()
}

func showStart(pinger *ping.Pinger, colors *Colors) error {
	_, err := fmt.Printf("PING %s (%s) %s(%s) bytes of data.\n",
		colors.Green.Sprintf("%s", pinger.Addr()),
		colors.Blue.Sprintf("%s", pinger.IPAddr()),
		colors.Blue.Sprintf("%d", size),
		colors.Blue.Sprintf("%d", size+28))
	if err != nil {
		return err
	}

	return nil
}

func pingCmd(arguments []string) error {
	timeZone := os.Getenv("TZ")
	if timeZone != "" {
		var err error

		time.Local, err = time.LoadLocation(timeZone)
		if err != nil {
			return err
		}
	}

	var host string = ""

	if net.ParseIP(arguments[0]) != nil {
		host = arguments[0]
	} else {
		url := regexp.MustCompile(`^([^:\/]+://)?([^:\/]+)(.+)?$`)

		asUrl := url.FindAllStringSubmatch(arguments[0], -1)
		if len(asUrl) > 0 && len(asUrl[0]) > 1 {
			host = asUrl[0][2]
		}
	}

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
		fmt.Printf("\n%s", showStatistics(stats, pinger, packets, colors, startTime, wasInterrupted, true))

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

				continue
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

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
)

const DATE string = "2006-01-02 15:04:05.000 MST"

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
	myPing, err := ping.NewPinger(host)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			<-c
			myPing.Stop()
		}
	}()

	myPing.Count = Count
	myPing.Size = Size
	myPing.Interval = Interval
	myPing.Timeout = Timeout
	myPing.RecordRtts = !(NoRTT)
	myPing.SetPrivileged(Privileged)

	if !(Color) {
		color.NoColor = true
	}

	var lastReceivedPacket = 0

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)

	myPing.OnRecv = func(pkt *ping.Packet) {
		currentPacket := pkt.Seq

		if Dropped && (lastReceivedPacket != currentPacket) {
			for c := lastReceivedPacket + 1; c < currentPacket; c++ {
				fmt.Printf("Packet %v dropped or arrived out of order.\n", c)
			}
			lastReceivedPacket = currentPacket
		}

		if Quiet {
			return
		}

		if Timestamp {
			timeStamp := time.Now().Format(DATE)
			fmt.Printf("%s | %d bytes from %s: icmp_seq=%d ttl=%v time=%v\n",
				timeStamp, pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Ttl, pkt.Rtt)
			return
		}

		fmt.Printf("%v from %v: icmp_seq=%v ttl=%v time=%v\n",
			blue.Sprintf("%v bytes", strconv.Itoa(pkt.Nbytes)), blue.Sprintf(pkt.IPAddr.String()), blue.Sprintf(strconv.Itoa(pkt.Seq)), blue.Sprintf(strconv.Itoa(pkt.Ttl)), blue.Sprintf(fmt.Sprintf("%v", pkt.Rtt)))
	}

	myPing.OnDuplicateRecv = func(pkt *ping.Packet) {
		if Timestamp {
			timeStamp := time.Now().Format(DATE)

			fmt.Printf("%s | %d bytes from %s: icmp_seq=%d ttl=%v time=%v %v\n",
				timeStamp, pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Ttl, pkt.Rtt, color.RedString("(DUP!)"))

			return
		}

		fmt.Printf("%d bytes from %s: icmp_seq=%d ttl=%v time=%v %v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Ttl, pkt.Rtt, color.RedString("(DUP!)"))
	}

	var startTime time.Time
	myPing.OnFinish = func(stats *ping.Statistics) {
		runTime := time.Since(startTime)

		fmt.Printf("\n--- %s ping statistics ---\n", green.Sprintf(stats.Addr))

		fmt.Printf("%v packets transmitted, %v received, %v, time %v\n",
			blue.Sprintf("%v", stats.PacketsSent), blue.Sprintf("%v", stats.PacketsRecv), HighlightPacketLoss(stats.PacketLoss), blue.Sprintf("%vms", runTime.Milliseconds()))

		fmt.Printf("rtt min/avg/max/mdev = %v/%v/%v/%v\n",
			blue.Sprintf("%v", stats.MinRtt), blue.Sprintf("%v", stats.AvgRtt), blue.Sprintf("%v", stats.MaxRtt), blue.Sprintf("%v", stats.StdDevRtt))

		sentBytes := int((float64(stats.PacketsSent) * (100 - stats.PacketLoss) * float64(myPing.Size)) / 100)
		receivedBytes := int((float64(stats.PacketsRecv) * (100 - stats.PacketLoss) * float64(myPing.Size)) / 100)

		fmt.Printf(
			"\n%s%v\n%s%v\n",
			"Sent = ", blue.Sprintf(humanReadableSize(sentBytes)),
			"Recv = ", blue.Sprintf(humanReadableSize(receivedBytes)),
		)
	}

	if Output == "<hostname>.log" {
		endLogging := logOutput(host + ".log")
		defer endLogging()
	} else if Output != "" {
		endLogging := logOutput(Output)
		defer endLogging()
	}

	startTime = time.Now()

	fmt.Printf("PING %s (%s):\n", green.Sprintf("%v", myPing.Addr()), blue.Sprintf("%v", myPing.IPAddr()))

	err = myPing.Run()

	if err != nil {
		panic(err)
	}
}

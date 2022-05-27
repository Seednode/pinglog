package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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
		log.Fatal(err)
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

	if NoColor {
		color.NoColor = true
	}

	blue := color.New(color.FgBlue)
	green := color.New(color.FgGreen)
	grey := color.New(color.FgHiBlack)
	red := color.New(color.FgRed)

	var lastReceivedPacket = 0

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
			fmt.Printf("%v | %v from %v: icmp_seq=%v ttl=%v time=%v\n",
				grey.Sprintf(timeStamp),
				blue.Sprintf("%v bytes", (pkt.Nbytes-8)),
				blue.Sprintf("%v", pkt.IPAddr),
				blue.Sprintf("%v", pkt.Seq),
				blue.Sprintf("%v", pkt.Ttl),
				blue.Sprintf("%v", pkt.Rtt.Truncate(time.Microsecond)))
			return
		}

		fmt.Printf("%v from %v: icmp_seq=%v ttl=%v time=%v\n",
			blue.Sprintf("%v bytes", (pkt.Nbytes-8)),
			blue.Sprintf("%v", pkt.IPAddr),
			blue.Sprintf("%v", pkt.Seq),
			blue.Sprintf("%v", pkt.Ttl),
			blue.Sprintf("%v", pkt.Rtt.Truncate(time.Microsecond)))
	}

	myPing.OnDuplicateRecv = func(pkt *ping.Packet) {
		if Timestamp {
			timeStamp := time.Now().Format(DATE)

			fmt.Printf("%v | %v from %v: icmp_seq=%v ttl=%v time=%v %v\n",
				grey.Sprintf(timeStamp),
				blue.Sprintf("%v bytes", (pkt.Nbytes-8)),
				blue.Sprintf("%v", pkt.IPAddr),
				blue.Sprintf("%v", pkt.Seq),
				blue.Sprintf("%v", pkt.Ttl),
				blue.Sprintf("%v", pkt.Rtt.Truncate(time.Microsecond)),
				red.Sprintf("(DUP!)"))

			return
		}

		fmt.Printf("%v from %v: icmp_seq=%v ttl=%v time=%v %v\n",
			blue.Sprintf("%v bytes", (pkt.Nbytes-8)),
			blue.Sprintf("%v", pkt.IPAddr),
			blue.Sprintf("%v", pkt.Seq),
			blue.Sprintf("%v", pkt.Ttl),
			blue.Sprintf("%v", pkt.Rtt.Truncate(time.Microsecond)),
			red.Sprintf("(DUP!)"))
	}

	var startTime time.Time

	myPing.OnFinish = func(stats *ping.Statistics) {
		runTime := time.Since(startTime)

		fmt.Printf("\n--- %v ping statistics ---\n", green.Sprintf(stats.Addr))

		fmt.Printf("%v packets transmitted, %v received, %v packet loss, time %v\n",
			blue.Sprintf("%v", stats.PacketsSent),
			blue.Sprintf("%v", stats.PacketsRecv),
			highlightPacketLoss(stats.PacketLoss),
			blue.Sprintf("%v", runTime.Truncate(time.Millisecond)))

		fmt.Printf("rtt min/avg/max/mdev = %v/%v/%v/%v\n",
			blue.Sprintf("%v", stats.MinRtt.Truncate(time.Microsecond)),
			blue.Sprintf("%v", stats.AvgRtt.Truncate(time.Microsecond)),
			blue.Sprintf("%v", stats.MaxRtt.Truncate(time.Microsecond)),
			blue.Sprintf("%v", stats.StdDevRtt.Truncate(time.Microsecond)))

		sentBytes := int(stats.PacketsSent * myPing.Size)
		receivedBytes := int(stats.PacketsRecv * myPing.Size)

		fmt.Printf("\n%v%v\n%v%v\n",
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

	fmt.Printf("PING %v (%v):\n",
		green.Sprintf("%v", myPing.Addr()),
		blue.Sprintf("%v", myPing.IPAddr()))

	err = myPing.Run()
	if err != nil {
		log.Fatal(err)
	}
}

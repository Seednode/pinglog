/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	packetsLost = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pinglog_packets_lost",
		Help: "The total number of lost packets",
	})
	packetsSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pinglog_packets_sent",
		Help: "The total number of sent packets",
	})
	packetsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pinglog_packets_received",
		Help: "The total number of received packets",
	})
)

func serveMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

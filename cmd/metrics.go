/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"net/http"

	"github.com/go-ping/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsCollector struct {
	Pinger      *ping.Pinger
	PacketsSent *prometheus.Desc
	PacketsRecv *prometheus.Desc
	RttMin      *prometheus.Desc
	RttMax      *prometheus.Desc
	RttAvg      *prometheus.Desc
	RttStdDev   *prometheus.Desc
}

var MinRTTs prometheus.Histogram
var AvgRTTs prometheus.Histogram
var MaxRTTs prometheus.Histogram
var StdDevRTTs prometheus.Histogram

func newMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		PacketsSent: prometheus.NewDesc("packets_sent",
			"Number of packets sent",
			nil, nil,
		),
		PacketsRecv: prometheus.NewDesc("packets_received",
			"Number of packets received",
			nil, nil,
		),
		RttMin: prometheus.NewDesc("rtt_min",
			"Minimum RTT (in microseconds)",
			nil, nil,
		),
		RttMax: prometheus.NewDesc("rtt_max",
			"Maximum RTT (in microseconds)",
			nil, nil,
		),
		RttAvg: prometheus.NewDesc("rtt_avg",
			"Average RTT (in microseconds)",
			nil, nil,
		),
		RttStdDev: prometheus.NewDesc("rtt_stddev",
			"Standard deviation of RTT (in microseconds)",
			nil, nil,
		),
	}
}

func (collector *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.PacketsSent
	ch <- collector.PacketsRecv
	ch <- collector.RttMin
	ch <- collector.RttAvg
	ch <- collector.RttMax
	ch <- collector.RttStdDev
}

func (collector *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(collector.PacketsSent, prometheus.CounterValue, float64(collector.Pinger.PacketsSent))
	ch <- prometheus.MustNewConstMetric(collector.PacketsRecv, prometheus.CounterValue, float64(collector.Pinger.PacketsRecv))
	MinRTTs.Observe(float64(collector.Pinger.Statistics().MinRtt.Microseconds()))
	ch <- MinRTTs
	AvgRTTs.Observe(float64(collector.Pinger.Statistics().AvgRtt.Microseconds()))
	ch <- AvgRTTs
	MaxRTTs.Observe(float64(collector.Pinger.Statistics().MaxRtt.Microseconds()))
	ch <- MaxRTTs
	StdDevRTTs.Observe(float64(collector.Pinger.Statistics().StdDevRtt.Microseconds()))
	ch <- StdDevRTTs
}

func initMetrics(registry *prometheus.Registry) {
	MinRTTs = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "min_rtts",
		Help:    "The minimum RTT.",
		Buckets: prometheus.ExponentialBuckets(1000, 2, 11),
	})
	AvgRTTs = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "avg_rtts",
		Help:    "The average RTT.",
		Buckets: prometheus.ExponentialBuckets(1000, 2, 11),
	})
	MaxRTTs = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "max_rtts",
		Help:    "The maximum RTT.",
		Buckets: prometheus.ExponentialBuckets(1000, 2, 11),
	})
	StdDevRTTs = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "stddev_rtts",
		Help:    "The standard deviation of RTTs.",
		Buckets: prometheus.ExponentialBuckets(1000, 2, 11),
	})
}

func serveMetrics(myPing *ping.Pinger) error {
	registry := prometheus.NewRegistry()
	initMetrics(registry)

	collector := newMetricsCollector()
	registry.MustRegister(collector)

	collector.Pinger = myPing

	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	err := http.ListenAndServe(fmt.Sprintf(":%v", Port), nil)
	if err != nil {
		return err
	}

	return nil
}

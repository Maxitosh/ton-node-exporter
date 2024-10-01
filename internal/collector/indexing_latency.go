package collector

import (
	"time"
	"ton-node-exporter/internal/fetcher"

	"github.com/prometheus/client_golang/prometheus"
)

const ElectorAddress = "Ef8zMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzM0vF"

var (
	tonLastElectorTxTimeDesc = prometheus.NewDesc(
		"ton_node_last_elector_tx_time",
		"Last Elector transaction time",
		[]string{}, nil,
	)
	tonIndexingLatencyDesc = prometheus.NewDesc(
		"ton_node_indexing_latency",
		"Time lag between the last Elector transaction and the current time",
		nil, nil,
	)
)

var TimeNow = time.Now

type TonIndexingLatencyCollector struct {
	fetcher fetcher.Fetcher
}

func NewTonIndexingLatencyCollector(
	fetcher fetcher.Fetcher,
) *TonIndexingLatencyCollector {
	return &TonIndexingLatencyCollector{
		fetcher: fetcher,
	}
}

func (collector *TonIndexingLatencyCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- tonLastElectorTxTimeDesc
	ch <- tonIndexingLatencyDesc
}

func (collector *TonIndexingLatencyCollector) Collect(ch chan<- prometheus.Metric) {
	lastElectorTxTime, err := collector.fetcher.FetchAddressLastTransactionTime(ElectorAddress)
	if err == nil {
		ch <- prometheus.MustNewConstMetric(
			tonLastElectorTxTimeDesc,
			prometheus.GaugeValue,
			float64(lastElectorTxTime),
		)
		ch <- prometheus.MustNewConstMetric(
			tonIndexingLatencyDesc,
			prometheus.GaugeValue,
			float64(TimeNow().Unix()-int64(lastElectorTxTime)),
		)
	}
}

package collector

import (
	"math"
	"sync"
	"ton-node-exporter/internal/fetcher"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	tonMasterChainBlockNumberDesc = prometheus.NewDesc(
		"ton_node_master_chain_block_number",
		"Master chain block number",
		[]string{"env"}, nil,
	)
	tonHeadLagDesc = prometheus.NewDesc(
		"ton_node_head_lag",
		"Head block lag",
		nil, nil,
	)
)

type TonBlockNumberCollector struct {
	localFetcher  fetcher.Fetcher
	globalFetcher fetcher.Fetcher
}

func NewTonBlockNumberCollector(
	localFetcher fetcher.Fetcher,
	globalFetcher fetcher.Fetcher,
) *TonBlockNumberCollector {
	return &TonBlockNumberCollector{
		localFetcher:  localFetcher,
		globalFetcher: globalFetcher,
	}
}

func (collector *TonBlockNumberCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- tonMasterChainBlockNumberDesc
	ch <- tonHeadLagDesc
}

func (collector *TonBlockNumberCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(2)
	var localBlockNumber, globalBlockNumber float64
	var localErr, globalErr error

	// Fetch local and global block numbers concurrently
	go func() {
		defer wg.Done()
		localBlockNumber, localErr = collector.localFetcher.FetchMasterChainBlockNumber()
	}()
	go func() {
		defer wg.Done()
		globalBlockNumber, globalErr = collector.globalFetcher.FetchMasterChainBlockNumber()
	}()
	wg.Wait()

	reportBlockNumber(ch, tonMasterChainBlockNumberDesc, localBlockNumber, localErr, "local")
	reportBlockNumber(ch, tonMasterChainBlockNumberDesc, globalBlockNumber, globalErr, "global")

	if localErr == nil && globalErr == nil {
		headLag := math.Max(0, globalBlockNumber-localBlockNumber)
		ch <- prometheus.MustNewConstMetric(
			tonHeadLagDesc,
			prometheus.GaugeValue,
			headLag,
		)
	}
}

// reportBlockNumber reports the block number metric if no error occurred
func reportBlockNumber(ch chan<- prometheus.Metric, desc *prometheus.Desc, blockNumber float64, err error, env string) {
	if err == nil {
		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			blockNumber,
			env,
		)
	}
}

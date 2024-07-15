package main

import (
	"context"
	"log"
	"net/http"
	"ton-node-exporter/internal/collector"
	"ton-node-exporter/internal/config"
	"ton-node-exporter/internal/fetcher"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed to load configuration: ", err.Error())
		return
	}

	// Initialize local lite client connection
	localClient := liteclient.NewConnectionPool()
	err = localClient.AddConnection(context.Background(), cfg.LiteServerAddr, cfg.LiteServerKey)
	if err != nil {
		log.Fatalln("Local client connection err: ", err.Error())
		return
	}
	localAPI := ton.NewAPIClient(localClient)

	// Initialize global lite client connection
	globalClient := liteclient.NewConnectionPool()
	err = globalClient.AddConnectionsFromConfigUrl(context.Background(), cfg.GlobalLiteClientConfig)
	if err != nil {
		log.Fatalln("Global client connection err: ", err.Error())
		return
	}
	globalAPI := ton.NewAPIClient(globalClient)

	tonBlockNumberCollector := collector.NewTonBlockNumberCollector(
		fetcher.NewADNLFetcher(localAPI),
		fetcher.NewADNLFetcher(globalAPI),
	)
	prometheus.MustRegister(tonBlockNumberCollector)

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Beginning to serve on port %s\n", cfg.ExporterPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ExporterPort, nil))
}

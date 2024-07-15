package collector

import "github.com/prometheus/client_golang/prometheus"

type Collector interface {
	prometheus.Collector
}

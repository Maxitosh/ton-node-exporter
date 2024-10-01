package collector

import (
	"strings"
	"testing"
	"time"
	mock_fetcher "ton-node-exporter/internal/fetcher/mocks"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTonIndexingLatencyCollector_Describe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fetcher := mock_fetcher.NewMockFetcher(ctrl)
	collector := NewTonIndexingLatencyCollector(fetcher)

	ch := make(chan *prometheus.Desc, 2)
	collector.Describe(ch)
	close(ch)

	var descriptions []*prometheus.Desc
	for desc := range ch {
		descriptions = append(descriptions, desc)
	}

	assert.Contains(t, descriptions, tonLastElectorTxTimeDesc)
	assert.Contains(t, descriptions, tonIndexingLatencyDesc)
}

func TestTonIndexingLatencyCollector_Collect(t *testing.T) {
	tests := []struct {
		name                string
		setupMocks          func(fetcher *mock_fetcher.MockFetcher)
		expectedMetrics     string
		expectedMetricCount int
	}{
		{
			name: "Successful fetch",
			setupMocks: func(fetcher *mock_fetcher.MockFetcher) {
				fetcher.EXPECT().FetchAddressLastTransactionTime(ElectorAddress).Return(uint32(1727767355), nil).AnyTimes()
			},
			expectedMetrics: `
# HELP ton_node_indexing_latency Time lag between the last Elector transaction and the current time
# TYPE ton_node_indexing_latency gauge
ton_node_indexing_latency 5
# HELP ton_node_last_elector_tx_time Last Elector transaction time
# TYPE ton_node_last_elector_tx_time gauge
ton_node_last_elector_tx_time 1.727767355e+09
`,
			expectedMetricCount: 2,
		},
		{
			name: "Error in Indexing latency fetcher",
			setupMocks: func(fetcher *mock_fetcher.MockFetcher) {
				fetcher.EXPECT().FetchAddressLastTransactionTime(ElectorAddress).Return(uint32(0), assert.AnError).AnyTimes()
			},
			expectedMetrics:     ``,
			expectedMetricCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 1727767355
			fixedTime := time.Date(2024, 10, 1, 7, 22, 40, 0, time.UTC)
			TimeNow = func() time.Time {
				return fixedTime
			}
			defer func() { TimeNow = time.Now }() // Reset after the test.

			fetcher := mock_fetcher.NewMockFetcher(ctrl)
			tt.setupMocks(fetcher)

			collector := NewTonIndexingLatencyCollector(fetcher)
			reg := prometheus.NewPedanticRegistry()
			reg.MustRegister(collector)

			metrics, err := reg.Gather()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMetricCount, len(metrics))

			if tt.expectedMetricCount > 0 {
				err := testutil.GatherAndCompare(reg, strings.NewReader(tt.expectedMetrics))
				assert.NoError(t, err)
			}
		})
	}
}

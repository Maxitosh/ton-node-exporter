package collector

import (
	"strings"
	"testing"
	mock_fetcher "ton-node-exporter/internal/fetcher/mocks"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTonBlockNumberCollector_Describe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	localFetcher := mock_fetcher.NewMockFetcher(ctrl)
	globalFetcher := mock_fetcher.NewMockFetcher(ctrl)
	collector := NewTonBlockNumberCollector(localFetcher, globalFetcher)

	ch := make(chan *prometheus.Desc, 2)
	collector.Describe(ch)
	close(ch)

	var descriptions []*prometheus.Desc
	for desc := range ch {
		descriptions = append(descriptions, desc)
	}

	assert.Contains(t, descriptions, tonMasterChainBlockNumberDesc)
	assert.Contains(t, descriptions, tonHeadLagDesc)
}

func TestTonBlockNumberCollector_Collect(t *testing.T) {
	tests := []struct {
		name                string
		setupMocks          func(localFetcher, globalFetcher *mock_fetcher.MockFetcher)
		expectedMetrics     string
		expectedMetricCount int
	}{
		{
			name: "Successful fetch",
			setupMocks: func(localFetcher, globalFetcher *mock_fetcher.MockFetcher) {
				localFetcher.EXPECT().FetchMasterChainBlockNumber().Return(10.0, nil).AnyTimes()
				globalFetcher.EXPECT().FetchMasterChainBlockNumber().Return(20.0, nil).AnyTimes()
			},
			expectedMetrics: `
# HELP ton_node_head_lag Head block lag
# TYPE ton_node_head_lag gauge
ton_node_head_lag 10
# HELP ton_node_master_chain_block_number Master chain block number
# TYPE ton_node_master_chain_block_number gauge
ton_node_master_chain_block_number{env="local"} 10
ton_node_master_chain_block_number{env="global"} 20
`,
			expectedMetricCount: 2,
		},
		{
			name: "Error in local fetcher",
			setupMocks: func(localFetcher, globalFetcher *mock_fetcher.MockFetcher) {
				localFetcher.EXPECT().FetchMasterChainBlockNumber().Return(0.0, assert.AnError).AnyTimes()
				globalFetcher.EXPECT().FetchMasterChainBlockNumber().Return(20.0, nil).AnyTimes()
			},
			expectedMetrics: `
# HELP ton_node_master_chain_block_number Master chain block number
# TYPE ton_node_master_chain_block_number gauge
ton_node_master_chain_block_number{env="global"} 20
`,
			expectedMetricCount: 1,
		},
		{
			name: "Error in global fetcher",
			setupMocks: func(localFetcher, globalFetcher *mock_fetcher.MockFetcher) {
				localFetcher.EXPECT().FetchMasterChainBlockNumber().Return(10.0, nil).AnyTimes()
				globalFetcher.EXPECT().FetchMasterChainBlockNumber().Return(0.0, assert.AnError).AnyTimes()
			},
			expectedMetrics: `
# HELP ton_node_master_chain_block_number Master chain block number
# TYPE ton_node_master_chain_block_number gauge
ton_node_master_chain_block_number{env="local"} 10
`,
			expectedMetricCount: 1,
		},
		{
			name: "Errors in both fetchers",
			setupMocks: func(localFetcher, globalFetcher *mock_fetcher.MockFetcher) {
				localFetcher.EXPECT().FetchMasterChainBlockNumber().Return(0.0, assert.AnError).AnyTimes()
				globalFetcher.EXPECT().FetchMasterChainBlockNumber().Return(0.0, assert.AnError).AnyTimes()
			},
			expectedMetrics:     ``,
			expectedMetricCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			localFetcher := mock_fetcher.NewMockFetcher(ctrl)
			globalFetcher := mock_fetcher.NewMockFetcher(ctrl)
			tt.setupMocks(localFetcher, globalFetcher)

			collector := NewTonBlockNumberCollector(localFetcher, globalFetcher)
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

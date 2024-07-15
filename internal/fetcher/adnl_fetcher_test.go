package fetcher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xssnick/tonutils-go/ton"
)

type MockAPIClient struct {
	mock.Mock
	MockGetMasterchainInfo func(ctx context.Context) (*ton.BlockIDExt, error)
}

func (m *MockAPIClient) GetMasterchainInfo(ctx context.Context) (*ton.BlockIDExt, error) {
	return m.MockGetMasterchainInfo(ctx)
}

func TestADNLFetcher_FetchMasterChainBlockNumber(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse *ton.BlockIDExt
		mockError    error
		expected     float64
		expectError  bool
	}{
		{
			name: "Successful fetch",
			mockResponse: &ton.BlockIDExt{
				SeqNo: 42,
			},
			expected:    42,
			expectError: false,
		},
		{
			name:        "Error fetching",
			mockError:   assert.AnError,
			expected:    -1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockAPIClient{
				MockGetMasterchainInfo: func(ctx context.Context) (*ton.BlockIDExt, error) {
					return tt.mockResponse, tt.mockError
				},
			}
			fetcher := NewADNLFetcher(mockClient)

			result, err := fetcher.FetchMasterChainBlockNumber()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

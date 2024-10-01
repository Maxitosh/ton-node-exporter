package fetcher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
)

type MockAPIClient struct {
	mock.Mock
	MockGetMasterchainInfo func(ctx context.Context) (*ton.BlockIDExt, error)
	MockGetAccount         func(ctx context.Context, block *ton.BlockIDExt, addr *address.Address) (*tlb.Account, error)
	MockListTransactions   func(
		ctx context.Context,
		addr *address.Address,
		limit uint32,
		lt uint64,
		txHash []byte,
	) ([]*tlb.Transaction, error)
}

func (m *MockAPIClient) GetMasterchainInfo(ctx context.Context) (*ton.BlockIDExt, error) {
	return m.MockGetMasterchainInfo(ctx)
}

func (m *MockAPIClient) GetAccount(
	ctx context.Context,
	block *ton.BlockIDExt,
	addr *address.Address,
) (*tlb.Account, error) {
	return m.MockGetAccount(ctx, block, addr)
}

func (m *MockAPIClient) ListTransactions(
	ctx context.Context,
	addr *address.Address,
	limit uint32,
	lt uint64,
	txHash []byte,
) ([]*tlb.Transaction, error) {
	return m.MockListTransactions(ctx, addr, limit, lt, txHash)
}

// Helper function to create a mock client.
func createMockClient(
	mockInfo func(ctx context.Context) (*ton.BlockIDExt, error),
	mockAccount func(ctx context.Context, block *ton.BlockIDExt, addr *address.Address) (*tlb.Account, error),
	mockTransactions func(
		ctx context.Context,
		addr *address.Address,
		limit uint32,
		lt uint64,
		txHash []byte,
	) ([]*tlb.Transaction, error),
) *MockAPIClient {
	return &MockAPIClient{
		MockGetMasterchainInfo: mockInfo,
		MockGetAccount:         mockAccount,
		MockListTransactions:   mockTransactions,
	}
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
			mockClient := createMockClient(func(ctx context.Context) (*ton.BlockIDExt, error) {
				return tt.mockResponse, tt.mockError
			}, nil, nil)

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

func TestADNLFetcher_FetchAddressLastTransactionTime(t *testing.T) {
	tests := []struct {
		name        string
		mockInfo    func(ctx context.Context) (*ton.BlockIDExt, error)
		mockAccount func(ctx context.Context, block *ton.BlockIDExt, addr *address.Address) (*tlb.Account, error)
		mockTxs     func(
			ctx context.Context,
			addr *address.Address,
			limit uint32,
			lt uint64,
			txHash []byte,
		) ([]*tlb.Transaction, error)
		expected    uint32
		expectError bool
	}{
		{
			name: "Successful fetch",
			mockInfo: func(ctx context.Context) (*ton.BlockIDExt, error) {
				return &ton.BlockIDExt{SeqNo: 42}, nil
			},
			mockAccount: func(ctx context.Context, block *ton.BlockIDExt, addr *address.Address) (*tlb.Account, error) {
				return &tlb.Account{LastTxLT: 123, LastTxHash: []byte("hash")}, nil
			},
			mockTxs: func(
				ctx context.Context,
				addr *address.Address,
				limit uint32,
				lt uint64,
				txHash []byte,
			) ([]*tlb.Transaction, error) {
				return []*tlb.Transaction{{Now: 456}}, nil
			},
			expected:    456,
			expectError: false,
		},
		{
			name: "Error fetching master block",
			mockInfo: func(ctx context.Context) (*ton.BlockIDExt, error) {
				return nil, assert.AnError
			},
			expected:    0,
			expectError: true,
		},
		{
			name: "Error fetching account",
			mockInfo: func(ctx context.Context) (*ton.BlockIDExt, error) {
				return &ton.BlockIDExt{SeqNo: 42}, nil
			},
			mockAccount: func(ctx context.Context, block *ton.BlockIDExt, addr *address.Address) (*tlb.Account, error) {
				return nil, assert.AnError
			},
			expected:    0,
			expectError: true,
		},
		{
			name: "Error fetching transaction",
			mockInfo: func(ctx context.Context) (*ton.BlockIDExt, error) {
				return &ton.BlockIDExt{SeqNo: 42}, nil
			},
			mockAccount: func(ctx context.Context, block *ton.BlockIDExt, addr *address.Address) (*tlb.Account, error) {
				return &tlb.Account{LastTxLT: 123, LastTxHash: []byte("hash")}, nil
			},
			mockTxs: func(
				ctx context.Context,
				addr *address.Address,
				limit uint32,
				lt uint64,
				txHash []byte,
			) ([]*tlb.Transaction, error) {
				return nil, assert.AnError
			},
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := createMockClient(tt.mockInfo, tt.mockAccount, tt.mockTxs)
			fetcher := NewADNLFetcher(mockClient)

			result, err := fetcher.FetchAddressLastTransactionTime("Ef8zMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzMzM0vF")

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

package fetcher

import (
	"context"
	"time"

	"github.com/xssnick/tonutils-go/ton"
)

type TonAPIClient interface {
	GetMasterchainInfo(ctx context.Context) (*ton.BlockIDExt, error)
}

type ADNLFetcher struct {
	client TonAPIClient
}

func NewADNLFetcher(client TonAPIClient) *ADNLFetcher {
	return &ADNLFetcher{
		client: client,
	}
}

func (fetcher *ADNLFetcher) FetchMasterChainBlockNumber() (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	masterChainInfo, err := fetcher.client.GetMasterchainInfo(ctx)
	if err != nil {
		return -1, err
	}

	return float64(masterChainInfo.SeqNo), nil
}

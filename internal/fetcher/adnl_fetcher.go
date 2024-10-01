package fetcher

import (
	"context"
	"time"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
)

type TonAPIClient interface {
	GetMasterchainInfo(ctx context.Context) (*ton.BlockIDExt, error)
	GetAccount(ctx context.Context, block *ton.BlockIDExt, addr *address.Address) (*tlb.Account, error)
	ListTransactions(
		ctx context.Context,
		addr *address.Address,
		num uint32,
		lt uint64,
		txHash []byte,
	) ([]*tlb.Transaction, error)
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

func (fetcher *ADNLFetcher) FetchAddressLastTransactionTime(addr string) (uint32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// we need fresh master block info to run get methods
	masterChainInfo, err := fetcher.client.GetMasterchainInfo(ctx)
	if err != nil {
		return 0, err
	}

	// fetch account info to get last tx hash and lt
	addrParsed := address.MustParseAddr(addr)
	account, err := fetcher.client.GetAccount(ctx, masterChainInfo, addrParsed)
	if err != nil {
		return 0, err
	}

	// fetch last transaction
	txs, err := fetcher.client.ListTransactions(ctx, address.MustParseAddr(addr), 1, account.LastTxLT, account.LastTxHash)
	if err != nil {
		return 0, err
	}

	electorTx := txs[0]
	return electorTx.Now, nil
}

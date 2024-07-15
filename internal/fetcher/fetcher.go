package fetcher

//go:generate mockgen -source=fetcher.go -destination=mocks/fetcher.go -package=mock_fetcher Fetcher
type Fetcher interface {
	FetchMasterChainBlockNumber() (float64, error)
}

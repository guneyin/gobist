package quote

import (
	"context"

	"github.com/guneyin/gobist/store"
)

type API struct {
	fetcher *fetcher
}

func NewAPI(store store.Store) *API {
	return &API{
		fetcher: newFetcher(store),
	}
}

func (a *API) GetQuote(_ context.Context, code string, opts ...OptionFunc) (*Quote, error) {
	return a.fetcher.GetQuote(code, opts...)
}

func (a *API) GetQuoteList(_ context.Context, symbols []string, opts ...OptionFunc) (*List, error) {
	return a.fetcher.GetQuoteList(symbols, opts...)
}

func (a *API) GetSymbolList(ctx context.Context) (*SymbolList, error) {
	return a.fetcher.GetSymbolList(ctx)
}

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

func (a *API) GetQuote(ctx context.Context, code string, opts ...OptionFunc) (*Quote, error) {
	return a.fetcher.GetQuote(ctx, code, opts...)
}

func (a *API) GetQuoteList(ctx context.Context, symbols []string, opts ...OptionFunc) (*List, error) {
	return a.fetcher.GetQuoteList(ctx, symbols, opts...)
}

func (a *API) GetSymbolList(ctx context.Context) (*SymbolList, error) {
	return a.fetcher.GetSymbolList(ctx)
}

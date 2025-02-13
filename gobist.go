package gobist

import (
	"context"

	"github.com/guneyin/gobist/kap"

	"github.com/guneyin/gobist/quote"
	"github.com/guneyin/gobist/store"
)

type Bist struct {
	quote *quote.API
	kap   *kap.API

	store store.Store
}

func New() *Bist {
	s := store.NewMemoryStore()
	return &Bist{
		quote: quote.NewAPI(s),
		kap:   kap.NewAPI(s),
	}
}

func (b *Bist) WithStore(store store.Store) *Bist {
	b.store = store
	return b
}

func (b *Bist) GetSymbolList(ctx context.Context) (*quote.SymbolList, error) {
	return b.quote.GetSymbolList(ctx)
}

func (b *Bist) GetQuote(ctx context.Context, symbols string, opts ...quote.OptionFunc) (*quote.Quote, error) {
	return b.quote.GetQuote(ctx, symbols, opts...)
}

func (b *Bist) GetQuoteList(ctx context.Context, symbols []string, opts ...quote.OptionFunc) (*quote.List, error) {
	return b.quote.GetQuoteList(ctx, symbols, opts...)
}

func (b *Bist) GetCompany(ctx context.Context, code string) (*kap.Company, error) {
	return b.kap.GetCompany(ctx, code)
}

func (b *Bist) GetCompanyWithShares(ctx context.Context, code string) (*kap.Company, error) {
	return b.kap.GetCompany(ctx, code, kap.WithShares())
}

package gobist

import (
	"context"

	"github.com/guneyin/gobist/quote"
	"github.com/guneyin/gobist/store"
)

type Bist struct {
	qc    *quote.API
	store store.Store
}

func New() *Bist {
	s := store.NewMemoryStore()
	return &Bist{
		qc: quote.NewAPI(s),
	}
}

func (b *Bist) WithStore(store store.Store) *Bist {
	b.store = store
	return b
}

func (b *Bist) GetSymbolList(ctx context.Context) (*quote.SymbolList, error) {
	return b.qc.GetSymbolList(ctx)
}

func (b *Bist) GetQuote(ctx context.Context, symbols string, opts ...quote.OptionFunc) (*quote.Quote, error) {
	return b.qc.GetQuote(ctx, symbols, opts...)
}

func (b *Bist) GetQuoteList(ctx context.Context, symbols []string, opts ...quote.OptionFunc) (*quote.List, error) {
	return b.qc.GetQuoteList(ctx, symbols, opts...)
}

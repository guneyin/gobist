package gobist

import (
	"github.com/guneyin/gobist/quote"
	"github.com/guneyin/gobist/store"
)

type Bist struct {
	qc    *quote.Client
	store store.Store
}

func New() *Bist {
	s := store.NewMemoryStore()
	return &Bist{
		qc: quote.NewClient(s),
	}
}

func (b *Bist) WithStore(store store.Store) *Bist {
	b.store = store
	return b
}

func (b *Bist) GetSymbolList() (*quote.SymbolList, error) {
	return b.qc.Fetcher().GetSymbolList()
}

func (b *Bist) GetQuote(symbols string, opts ...quote.OptionFunc) (*quote.Quote, error) {
	return b.qc.Fetcher().GetQuote(symbols, opts...)
}

func (b *Bist) GetQuoteList(symbols []string, opts ...quote.OptionFunc) (*quote.List, error) {
	return b.qc.Fetcher().GetQuoteList(symbols, opts...)
}

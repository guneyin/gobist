package gobist

import "time"

type Bist struct {
	api   *api
	store Store
}

func New() *Bist {
	store := newMemoryStore()
	return &Bist{
		api: newApi(store),
	}
}

func (b *Bist) WithStore(store Store) *Bist {
	b.store = store
	return b
}

func (b *Bist) GetSymbolList() (*SymbolList, error) {
	return b.api.getSymbolList()
}

func (b *Bist) GetQuote(symbols []string) (*QuoteList, error) {
	return b.api.getQuote(symbols)
}

func (b *Bist) GetQuoteWithHistory(symbols []string, period ...time.Time) (*QuoteList, error) {
	return b.api.getQuote(symbols, period...)
}

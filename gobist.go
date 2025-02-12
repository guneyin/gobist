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

func (b *Bist) GetQuote(symbols string, period ...time.Time) (*Quote, error) {
	return b.api.getQuote(symbols, period...)
}
func (b *Bist) GetQuoteList(symbols []string, period ...time.Time) (*QuoteList, error) {
	return b.api.getQuoteList(symbols, period...)
}

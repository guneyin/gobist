package gobist

import "time"

type Bist struct {
	api   *api
	store Store
}

func New() *Bist {
	return &Bist{
		api:   newApi(),
		store: newMemoryStore(),
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
	return b.api.getQuote(symbols, nil)
}

func (b *Bist) GetQuoteWithHistory(symbols []string, date ...time.Time) (*QuoteList, error) {
	var dt *time.Time

	if len(date) > 0 {
		dt = &date[0]
	}

	return b.api.getQuote(symbols, dt)
}

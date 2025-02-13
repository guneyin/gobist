package gobist

type Bist struct {
	api   *api
	store Store
}

func New() *Bist {
	store := newMemoryStore()
	return &Bist{
		api: newAPI(store),
	}
}

func (b *Bist) WithStore(store Store) *Bist {
	b.store = store
	return b
}

func (b *Bist) GetSymbolList() (*SymbolList, error) {
	return b.api.qf().GetSymbolList()
}

func (b *Bist) GetQuote(symbols string, opts ...QuoteOptionFunc) (*Quote, error) {
	return b.api.qf().GetQuote(symbols, opts...)
}

func (b *Bist) GetQuoteList(symbols []string, opts ...QuoteOptionFunc) (*QuoteList, error) {
	return b.api.qf().GetQuoteList(symbols, opts...)
}

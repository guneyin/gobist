package gobist

import "time"

type Bist struct {
	api *yahooApi
}

func New() (*Bist, error) {
	api, err := newApi()
	if err != nil {
		return nil, err
	}

	return &Bist{api: api}, nil
}

func (b *Bist) GetQuote(symbol string) (*Quote, error) {
	return b.api.getQuote(symbol, nil)
}

func (b *Bist) GetQuoteWithHistory(symbol string, date ...time.Time) (*Quote, error) {
	var dt *time.Time

	if len(date) > 0 {
		dt = &date[0]
	}

	return b.api.getQuote(symbol, dt)
}

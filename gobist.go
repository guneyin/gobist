package gobist

import (
	"time"
)

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
	return b.api.getQuote(symbol)
}

func (b *Bist) GetQuoteWithHistory(symbol string, date time.Time) (*Quote, error) {
	return b.api.getQuoteWithHistory(symbol, date)
}

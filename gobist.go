package gobist

import (
	"time"

	"github.com/guneyin/gobist/internal"
)

type Bist struct {
	api *internal.Api
}

func New() (*Bist, error) {
	api, err := internal.NewApi()
	if err != nil {
		return nil, err
	}

	return &Bist{api: api}, nil
}

func (b *Bist) GetQuote(symbol string) (*internal.Quote, error) {
	return b.api.GetQuote(symbol)
}

func (b *Bist) GetQuoteWithHistory(symbol string, date time.Time) (*internal.Quote, error) {
	return b.api.GetQuoteWithHistory(symbol, date)
}

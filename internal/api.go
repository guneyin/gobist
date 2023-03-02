package internal

import (
	"fmt"
	"time"
)

type Api struct {
	c *Client
}

func NewApi() (*Api, error) {
	c, err := newClient()
	if err != nil {
		return nil, err
	}

	return &Api{c: c}, nil
}

func (a *Api) GetQuote(symbol string) (*Quote, error) {
	r, err := a.c.Get(fmt.Sprintf(quotePath, symbol))
	if err != nil {
		return nil, err
	}

	y := new(yahooQuote)
	y.Unmarshall(r)

	oc := y.OptionChain.Result
	if len(oc) == 0 {
		return nil, ErrInvalidPayload
	}

	q := oc[0].Quote

	return &Quote{
		Symbol:  symbol,
		Name:    q.ShortName,
		Price:   q.RegularMarketPrice,
		History: nil,
	}, nil
}

func (a *Api) GetQuoteWithHistory(symbol string, date time.Time) (*Quote, error) {
	q, err := a.GetQuote(symbol)
	if err != nil {
		return nil, err
	}

	r, err := a.c.Get(fmt.Sprintf(historyPath, symbol, date.Unix(), date.Unix()))
	if err != nil {
		return nil, err
	}

	y := new(yahooHistory)
	y.Parse(r)

	q.History = &History{
		Date:  date,
		Price: y.Close,
	}

	return q, nil
}

package quote

import (
	"errors"

	"github.com/guneyin/gobist/store"
)

const (
	yahooCrumbPath = "/v1/test/getcrumb"
	yahooChartPath = "/v8/finance/chart/%s.IS?includeAdjustedClose=true&interval=1d&period1=%d&period2=%d"

	twSymbolListURL = "https://scanner.tradingview.com/turkey/scan"
)

var (
	errNoDataFound         = errors.New("no data found")
	errHistoryDataNotFound = errors.New("history data not found")
)

type Client struct {
	c *client
}

func NewClient(store store.Store) *Client {
	c := newClient(store)

	return &Client{
		c: c,
	}
}

func (c *Client) Fetcher() *Fetcher {
	return newQuoteFetcher(c.c)
}

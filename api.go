package gobist

import (
	"errors"
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

type api struct {
	c *client
}

func newAPI(store Store) *api {
	c := newClient(store)

	return &api{
		c: c,
	}
}

func (a *api) qf() *quoteFetcher {
	return newQuoteFetcher(a.c)
}

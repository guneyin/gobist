package quote

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/guneyin/gobist/store"

	"github.com/imroc/req/v3"
	"github.com/shopspring/decimal"
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

type fetcher struct {
	c   *client
	opt *Option
}

func newFetcher(store store.Store) *fetcher {
	return &fetcher{
		c:   newClient(store),
		opt: newDefaultOptions(),
	}
}

func (f *fetcher) applyOptions(opts ...OptionFunc) {
	for _, opt := range opts {
		opt(f)
	}
}

func (f *fetcher) GetSymbolList(ctx context.Context) (*SymbolList, error) {
	resp := symbolListResponse{}
	body := `
		{
			"columns": [
				"name",
				"description",
				"logoid"
			],
			"options": {
				"lang": "tr"
			},
			"sort": {
				"sortBy": "name",
				"sortOrder": "asc",
				"nullsFirst": false
			},
			"preset": "all_stocks"
		}`

	r, err := f.c.general.R().
		SetContext(ctx).
		SetBodyJsonString(body).
		SetSuccessResult(&resp).
		Post(twSymbolListURL)
	if err != nil {
		return nil, err
	}

	if r.IsErrorState() {
		return nil, r.Err
	}

	res := new(SymbolList)

	return res.FromDTO(&resp), nil
}

func (f *fetcher) GetQuote(ctx context.Context, symbol string, opts ...OptionFunc) (*Quote, error) {
	list, err := f.GetQuoteList(ctx, []string{symbol}, opts...)
	if err != nil {
		return nil, err
	}

	if list.Count == 0 {
		return nil, errHistoryDataNotFound
	}

	return &list.Items[0], nil
}

func (f *fetcher) GetQuoteList(ctx context.Context, symbols []string, opts ...OptionFunc) (*List, error) {
	f.applyOptions(opts...)

	quoteList := &List{
		Count: len(symbols),
		Items: make([]Quote, len(symbols)),
	}

	wg := sync.WaitGroup{}
	wg.Add(len(symbols))

	for i, symbol := range symbols {
		item := &quoteList.Items[i]
		item.Symbol = symbol

		go f.syncQuote(ctx, item, &wg)
	}

	wg.Wait()

	return quoteList, nil
}

func (f *fetcher) syncQuote(ctx context.Context, q *Quote, wg *sync.WaitGroup) {
	defer wg.Done()

	data, err := f.fetchYahooChart(ctx, q.Symbol, f.opt.period.begin.Unix())
	if err != nil {
		q.SetError(err)
		return
	}

	q.Name = data.Chart.Result[0].Meta.ShortName
	q.Price = decimal.NewFromFloat(data.Chart.Result[0].Meta.RegularMarketPrice).Truncate(2).String()

	if !f.opt.period.isSingleDay() {
		h := History{}

		if !data.adjCloseCheck() {
			q.SetError(errHistoryDataNotFound)
			return
		}

		dt, cp := data.getClosePrice()
		h.SetBegin(dt, cp)

		data, err = f.fetchYahooChart(ctx, q.Symbol, f.opt.period.end.Unix())
		if err != nil {
			q.SetError(err)
			return
		}

		if !data.adjCloseCheck() {
			q.SetError(errHistoryDataNotFound)
			return
		}
		dt, cp = data.getClosePrice()
		h.SetEnd(dt, cp)

		if h.IsValid() {
			bp, _ := decimal.NewFromString(h.Begin.Price)
			ep, _ := decimal.NewFromString(h.End.Price)
			ratio := ep.Sub(bp).Mul(decimal.NewFromInt(100)).Div(ep)
			h.Change = HistoryChange{
				ByRatio:  ratio.Truncate(2).String(),
				ByAmount: ep.Sub(bp).Truncate(2).String(),
			}

			q.History = h
		}
	}
}

func (f *fetcher) fetchYahooChart(ctx context.Context, symbol string, ts int64) (*quoteDTO, error) {
	data := &quoteDTO{}

	rq := f.c.yahoo.R().
		SetContext(ctx).
		SetRetryCount(1).
		SetRetryFixedInterval(1 * time.Second).
		AddRetryHook(func(_ *req.Response, _ error) {
			_, err := setYahooCrumb()
			if err != nil {
				log.Printf("failed to set yahoo crumb: %v", err)
			}
		}).
		AddRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				log.Printf("failed to set yahoo crumb: %v", err)
			}

			return resp.StatusCode == http.StatusUnauthorized
		}).
		SetSuccessResult(&data)

	tsBegin := time.Unix(ts, 0).AddDate(0, 0, -15).Unix()
	url := fmt.Sprintf(yahooChartPath, symbol, tsBegin, ts)
	r, err := rq.Get(url)
	if err != nil {
		return nil, err
	}

	if r.IsErrorState() {
		return nil, errors.New(r.Status)
	}

	if len(data.Chart.Result) == 0 {
		return nil, errNoDataFound
	}

	return data, nil
}

func setYahooCrumb() (string, error) {
	res, err := req.C().R().Get(yahooCrumbPath)
	if err != nil {
		return "", fmt.Errorf("crumb error: %w", err)
	}

	return res.String(), nil
}

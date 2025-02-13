package quote

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"github.com/shopspring/decimal"
)

type Fetcher struct {
	c   *client
	opt *Option
}

func newQuoteFetcher(client *client) *Fetcher {
	return &Fetcher{
		c:   client,
		opt: NewDefaultOptions(),
	}
}

func (f *Fetcher) applyOptions(opts ...OptionFunc) {
	for _, opt := range opts {
		opt(f)
	}
}

func (f *Fetcher) GetSymbolList() (*SymbolList, error) {
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

func (f *Fetcher) GetQuote(symbol string, opts ...OptionFunc) (*Quote, error) {
	list, err := f.GetQuoteList([]string{symbol}, opts...)
	if err != nil {
		return nil, err
	}

	if list.Count == 0 {
		return nil, errHistoryDataNotFound
	}

	return &list.Items[0], nil
}

func (f *Fetcher) GetQuoteList(symbols []string, opts ...OptionFunc) (*List, error) {
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

		go f.syncQuote(item, &wg)
	}

	wg.Wait()

	return quoteList, nil
}

func (f *Fetcher) syncQuote(q *Quote, wg *sync.WaitGroup) {
	defer wg.Done()

	data, err := f.fetchYahooChart(q.Symbol, f.opt.period.begin.Unix())
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

		data, err = f.fetchYahooChart(q.Symbol, f.opt.period.end.Unix())
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

func (f *Fetcher) fetchYahooChart(symbol string, ts int64) (*quoteDTO, error) {
	data := &quoteDTO{}

	rq := f.c.yahoo.R().
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

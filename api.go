package gobist

import (
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	fp "github.com/nikolaydubina/fpmoney"
	"net/http"
	"sync"
	"time"
)

const (
	yahooCrumbPath = "/v1/test/getcrumb"
	yahooChartPath = "/v8/finance/chart/%s.IS?includeAdjustedClose=true&interval=1d&period1=%d&period2=%d"

	twSymbolListUrl = "https://scanner.tradingview.com/turkey/scan"
)

var (
	errNoDataFound         = errors.New("no data found")
	errHistoryDataNotFound = errors.New("history data not found")
)

type api struct {
	c *client
}

func newApi(store Store) *api {
	return &api{c: newClient(store)}
}

func (a *api) getSymbolList() (*SymbolList, error) {
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

	r, err := a.c.general.R().
		SetBodyJsonString(body).
		SetSuccessResult(&resp).
		Post(twSymbolListUrl)
	if err != nil {
		return nil, err
	}

	if r.IsErrorState() {
		return nil, r.Err
	}

	res := new(SymbolList)

	return res.fromDTO(&resp), nil
}

type pdate time.Time

func (pd *pdate) String() string {
	return time.Time(*pd).Format("2006-01-02")
}

func (pd *pdate) Unix() int64 {
	return time.Time(*pd).Unix()
}

func isToday(d time.Time) bool {
	return d.Format(time.DateOnly) == time.Now().Format(time.DateOnly)
}

func (pd *pdate) Set(dt time.Time) {
	if isToday(dt) {
		dt = time.Now()
	} else {
		y, m, d := dt.Date()
		dt = time.Date(y, m, d, 11, 0, 0, 0, dt.Location())
	}

	*pd = pdate(dt)
}

type period struct {
	begin, end pdate
}

func (p *period) Begin() *pdate {
	return &p.begin
}

func (p *period) End() *pdate {
	return &p.end
}

func (p *period) IsSingleDay() bool {
	return p.begin.String() == p.end.String()
}

func (a *api) getQuote(symbols []string, dates ...time.Time) (*QuoteList, error) {
	var p period

	switch len(dates) {
	case 0:
		dtToday := time.Now()
		p.Begin().Set(dtToday)
		p.End().Set(dtToday)
	case 1:
		p.Begin().Set(dates[0])
		p.End().Set(dates[0])
	default:
		p.Begin().Set(dates[0])
		p.End().Set(dates[1])
	}

	quoteList := &QuoteList{
		Count: len(symbols),
		Items: make([]Quote, len(symbols)),
	}

	wg := sync.WaitGroup{}
	wg.Add(len(symbols))

	for i, symbol := range symbols {
		go func(q *Quote) {
			defer wg.Done()

			q.Symbol = symbol

			data, err := a.fetchYahooChart(symbol, p.Begin().Unix())
			if err != nil {
				q.SetError(err)
				return
			}

			q.Name = data.Chart.Result[0].Meta.ShortName
			q.Price = fp.FromFloat(data.Chart.Result[0].Meta.RegularMarketPrice, fp.TRY)

			if !p.IsSingleDay() {
				h := History{}

				if !data.adjCloseCheck() {
					q.SetError(errHistoryDataNotFound)
					return
				}

				dt, cp := data.getClosePrice()
				h.SetBegin(dt, cp)

				data, err = a.fetchYahooChart(symbol, p.End().Unix())
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
					ratio := h.End.Price.Sub(h.Begin.Price).Mul(100).Float64() / h.End.Price.Float64()
					h.Change = HistoryChange{
						ByRatio:  fp.FromFloat(ratio, fp.TRY),
						ByAmount: h.End.Price.Sub(h.Begin.Price),
					}

					q.History = h
				}
			}
		}(&quoteList.Items[i])
	}

	wg.Wait()

	return quoteList, nil
}

func (a *api) fetchYahooChart(symbol string, ts int64) (*quoteDTO, error) {
	data := &quoteDTO{}

	rq := a.c.yahoo.R().
		SetRetryCount(1).
		SetRetryFixedInterval(1 * time.Second).
		AddRetryHook(func(resp *req.Response, err error) {
			_, err = setYahooCrumb()
		}).
		AddRetryCondition(func(resp *req.Response, err error) bool {
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
		return "", fmt.Errorf("crumb error: %v", err)
	}

	return res.String(), nil
}

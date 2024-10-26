package gobist

import (
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
	"strconv"
	"time"
)

const (
	yahooCrumbPath = "/v1/test/getcrumb"
	yahooChartPath = "/v8/finance/chart/%s.IS?includeAdjustedClose=true&interval=1d&period1=%v&period2=%v"

	twSymbolListUrl = "https://scanner.tradingview.com/turkey/scan"
)

var (
	ErrNoDataFound = errors.New("NO_DATA_FOUND")
)

type api struct {
	c *client
}

func newApi() *api {
	return &api{c: newClient()}
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

type TPeriod [2]string

func (p *TPeriod) SetBegin(s string) {
	p[0] = s
}

func (p *TPeriod) SetEnd(s string) {
	p[1] = s
}

func (p *TPeriod) Begin() string {
	return p[0]
}

func (p *TPeriod) End() string {
	return p[1]
}

func (p *TPeriod) IsSingleDay() bool {
	return p.Begin() == p.End()
}

func (a *api) getQuote(symbols []string, dates ...time.Time) (*QuoteList, error) {
	var period TPeriod

	dtToday := adjustDate(time.Now())

	switch len(dates) {
	case 0:
		period.SetBegin(dtToday)
		period.SetEnd(dtToday)
	case 1:
		period.SetBegin(adjustDate(dates[0]))
		period.SetEnd(dtToday)
	default:
		period.SetBegin(adjustDate(dates[0]))
		period.SetEnd(adjustDate(dates[1]))
	}

	quoteList := &QuoteList{
		Count: len(symbols),
		Items: make([]Quote, len(symbols)),
	}

	for i, symbol := range symbols {
		q := &quoteList.Items[i]
		q.Symbol = symbol

		data, err := a.fetchYahooChart(symbol, period.Begin())
		if err != nil {
			q.SetError(err.Error())
			continue
		}

		if len(data.Chart.Result) == 0 {
			q.SetError("no data found")
			continue
		}

		q.Name = data.Chart.Result[0].Meta.ShortName
		q.Price = data.Chart.Result[0].Meta.RegularMarketPrice

		if !period.IsSingleDay() {
			h := History{}

			if len(data.Chart.Result[0].Indicators.Adjclose) == 0 {
				q.SetError("close price not found")
				continue
			}

			h.SetBegin(period.Begin(), data.Chart.Result[0].Indicators.Adjclose[0].Adjclose[0])

			data, err = a.fetchYahooChart(symbol, period.End())
			if err != nil {
				q.SetError(err.Error())
				continue
			}

			if len(data.Chart.Result[0].Indicators.Adjclose[0].Adjclose) == 0 {
				q.SetError("close price not found")
				continue
			}
			h.SetEnd(period.Begin(), data.Chart.Result[0].Indicators.Adjclose[0].Adjclose[0])

			if h.IsValid() {
				h.Change = HistoryChange{
					ByRatio:  (h.End.Price - h.Begin.Price) * (100 / h.End.Price),
					ByAmount: h.End.Price - h.Begin.Price,
				}

				q.History = h
			}
		}
	}

	return quoteList, nil
}

func (a *api) fetchYahooChart(symbol, dt string) (*quote, error) {
	data := &quote{}

	rq := a.c.yahoo.R().
		SetRetryCount(1).
		SetRetryFixedInterval(1 * time.Second).
		AddRetryHook(func(resp *req.Response, err error) {
			setYahooCrumb()
		}).
		AddRetryCondition(func(resp *req.Response, err error) bool {
			return resp.StatusCode == http.StatusUnauthorized
		}).
		SetSuccessResult(&data)

	url := fmt.Sprintf(yahooChartPath, symbol, dt, dt)
	r, err := rq.Get(url)
	if err != nil {
		return nil, err
	}

	if r.IsErrorState() {
		return nil, errors.New(r.Status)
	}

	if len(data.Chart.Result) == 0 {
		return nil, ErrNoDataFound
	}

	return data, nil
}

func setYahooCrumb() string {
	res, err := req.C().R().Get(yahooCrumbPath)
	if err != nil {
		fmt.Printf("crumb error: %v\n", err)
	}

	crumb := res.String()

	fmt.Printf("crumb has been set: %v\n", crumb)
	return crumb
}

func adjustDate(d time.Time) string {
	dta, _ := time.Parse(time.DateOnly, d.Format("2006-01-02"))
	dta = dta.Add(time.Hour * 20)

	return strconv.Itoa(int(dta.Unix()))
}

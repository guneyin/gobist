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

func (a *api) getQuote(symbols []string, date *time.Time) (*QuoteList, error) {
	data := quote{}

	var (
		dtAdjusted time.Time
		dt         string
	)

	if date != nil {
		dtAdjusted, _ = time.Parse(time.DateOnly, date.Format("2006-01-02"))
		dtAdjusted = dtAdjusted.Add(time.Hour * 20)

		dt = strconv.Itoa(int(dtAdjusted.Unix()))
	}

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

	quoteList := &QuoteList{
		Count: len(symbols),
		Items: make([]Quote, len(symbols)),
	}
	for i, symbol := range symbols {
		url := fmt.Sprintf(yahooChartPath, symbol, dt, dt)
		r, err := rq.Get(url)
		if err != nil {
			return nil, err
		}

		if r.IsErrorState() {
			return nil, err
		}

		if len(data.Chart.Result) == 0 {
			return nil, ErrNoDataFound
		}

		q := data.Chart.Result[0]

		h := &History{}
		if date != nil {
			if len(q.Indicators.Adjclose[0].Adjclose) > 0 {
				h.Date = &dtAdjusted
				h.Price = q.Indicators.Adjclose[0].Adjclose[0]
				h.Change = &Change{
					ByRatio:  (q.Meta.RegularMarketPrice - h.Price) * (100 / q.Meta.RegularMarketPrice),
					ByAmount: q.Meta.RegularMarketPrice - h.Price,
				}
			}
		}

		quoteList.Items[i] = Quote{
			Symbol:  symbol,
			Name:    q.Meta.ShortName,
			Price:   q.Meta.RegularMarketPrice,
			History: h,
		}
	}

	return quoteList, nil
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

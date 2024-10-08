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
	baseURL   = "https://query2.finance.yahoo.com"
	crumbPath = "/v1/test/getcrumb"
	chartPath = "/v8/finance/chart/%s.IS?includeAdjustedClose=true&interval=1d&period1=%v&period2=%v"
)

var (
	ErrNoDataFound = errors.New("NO_DATA_FOUND")
)

type yahooApi struct {
	c *req.Client
}

func newApi() (*yahooApi, error) {
	return &yahooApi{c: newClient()}, nil
}

func (y *yahooApi) getQuote(symbol string, date *time.Time) (*Quote, error) {
	quote := yahooQuote{}

	var (
		dtAdjusted time.Time
		dt         string
	)

	if date != nil {
		dtAdjusted, _ = time.Parse(time.DateOnly, date.Format("2006-01-02"))
		dtAdjusted = dtAdjusted.Add(time.Hour * 20)

		dt = strconv.Itoa(int(dtAdjusted.Unix()))
	}

	url := fmt.Sprintf(chartPath, symbol, dt, dt)

	r, err := y.c.
		//DevMode().
		R().
		SetRetryCount(1).
		SetRetryFixedInterval(1 * time.Second).
		AddRetryHook(func(resp *req.Response, err error) {
			setYahooCrumb()
		}).
		AddRetryCondition(func(resp *req.Response, err error) bool {
			return resp.StatusCode == http.StatusUnauthorized
		}).
		SetSuccessResult(&quote).
		Get(url)
	if err != nil {
		return nil, err
	}

	if r.IsErrorState() {
		return nil, r.Err
	}

	if len(quote.Chart.Result) == 0 {
		return nil, ErrNoDataFound
	}

	q := quote.Chart.Result[0]

	h := &History{}
	if date != nil {
		if len(q.Indicators.Adjclose[0].Adjclose) > 0 {
			h.Date = &dtAdjusted
			h.Price = q.Indicators.Adjclose[0].Adjclose[0]
		}
	}

	return &Quote{
		Symbol:  symbol,
		Name:    q.Meta.ShortName,
		Price:   q.Meta.RegularMarketPrice,
		History: h,
	}, nil
}

package gobist

import (
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"net/http"
	"time"
)

const (
	baseURL     = "https://query2.finance.yahoo.com"
	crumbPath   = "/v1/test/getcrumb"
	quotePath   = "/v7/finance/options/%s.IS?crumb=%s"
	historyPath = "/v8/finance/chart/%s.IS?period1=%v&period2=%v&interval=1m"
)

var (
	ErrInvalidPayload = errors.New("INVALID_PAYLOAD")
)

type yahooApi struct {
	c *req.Client
}

func newApi() (*yahooApi, error) {
	return &yahooApi{c: newClient()}, nil
}

func (y *yahooApi) getQuote(symbol string) (*Quote, error) {
	quote := yahooQuote{}

	r, err := y.c.
		//DevMode().
		R().
		SetRetryCount(1).
		SetRetryFixedInterval(1 * time.Second).
		AddRetryHook(func(resp *req.Response, err error) {
			crumb = getYahooCrumb()
			resp.Request.SetURL(fmt.Sprintf(quotePath, symbol, crumb))
		}).
		AddRetryCondition(func(resp *req.Response, err error) bool {
			return resp.StatusCode == http.StatusUnauthorized
		}).
		SetSuccessResult(&quote).
		Get(fmt.Sprintf(quotePath, symbol, crumb))
	if err != nil {
		return nil, err
	}

	if r.IsErrorState() {
		return nil, r.Err
	}

	oc := quote.OptionChain.Result
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

func (y *yahooApi) getQuoteWithHistory(symbol string, date time.Time) (*Quote, error) {
	q, err := y.getQuote(symbol)
	if err != nil {
		return nil, err
	}

	history := yahooHistory{}

	r, err := y.c.R().
		SetSuccessResult(&history).
		Get(fmt.Sprintf(historyPath, symbol, date.Unix(), date.Unix()))
	if err != nil {
		return nil, err
	}

	if r.IsErrorState() {
		return nil, r.Err
	}

	q.History = &History{
		Date:  date,
		Price: history.Close,
	}

	return q, nil
}

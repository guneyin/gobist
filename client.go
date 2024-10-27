package gobist

import (
	"github.com/imroc/req/v3"
)

const (
	yahooBaseURL = "https://query2.finance.yahoo.com"
)

type client struct {
	yahoo   *req.Client
	general *req.Client
}

func newClient(store Store) *client {
	return &client{
		yahoo:   newYahooClient(store),
		general: req.NewClient(),
	}
}

func newYahooClient(store Store) *req.Client {
	jar := newCookieJar(store)

	headers := make(map[string]string)
	headers["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8"
	headers["Host"] = "query2.finance.yahoo.com"

	return req.NewClient().
		SetCookieJar(jar).
		SetBaseURL(yahooBaseURL).
		SetCommonHeaders(headers).
		SetUserAgent("Mozilla/5.0 (Windows NT 10.0; Windows; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")
}

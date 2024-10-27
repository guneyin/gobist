package gobist

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	cookiesKey = "cookies"
)

type CookieJar struct {
	jar   http.CookieJar
	store Store

	cookies []*http.Cookie
}

var _ http.CookieJar = (*CookieJar)(nil)

func newCookieJar(store Store) *CookieJar {
	jar, _ := cookiejar.New(nil)

	return &CookieJar{
		jar:   jar,
		store: store,
	}
}

func (c CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if err := c.save(); err != nil {
		c.jar.SetCookies(u, cookies)
	}
}

func (c CookieJar) Cookies(u *url.URL) []*http.Cookie {
	return c.jar.Cookies(u)
}

func (c CookieJar) save() error {
	data, err := json.Marshal(c.cookies)
	if err != nil {
		return err
	}

	return c.store.Set(cookiesKey, string(data))
}

func (c CookieJar) load() []*http.Cookie {
	data, err := c.store.Get(cookiesKey)
	if err != nil {
		return nil
	}

	res := make([]*http.Cookie, 0)
	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		return nil
	}

	return res
}

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

type cookieJar struct {
	jar   http.CookieJar
	store Store

	cookies []*http.Cookie
}

var _ http.CookieJar = (*cookieJar)(nil)

func newCookieJar(store Store) *cookieJar {
	jar, _ := cookiejar.New(nil)

	return &cookieJar{
		jar:   jar,
		store: store,
	}
}

func (c cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if err := c.save(); err != nil {
		c.jar.SetCookies(u, cookies)
	}
}

func (c cookieJar) Cookies(u *url.URL) []*http.Cookie {
	return c.jar.Cookies(u)
}

func (c cookieJar) save() error {
	data, err := json.Marshal(c.cookies)
	if err != nil {
		return err
	}

	return c.store.Set(cookiesKey, string(data))
}

func (c cookieJar) load() []*http.Cookie {
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

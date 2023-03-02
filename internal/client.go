package internal

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type method string

const (
	get  method = "get"
	post method = "post"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func newClient() (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		baseURL:    u,
		httpClient: http.DefaultClient,
	}, nil
}

func (c *Client) call(url string, method method, body ...[]byte) ([]byte, error) {
	u := c.baseURL.String() + url

	var (
		r   *http.Response
		err error
	)

	switch method {
	case get:
		r, err = c.httpClient.Get(u)
		if err != nil {
			return nil, err
		}
	case post:
		var b *bytes.Buffer

		if body[0] != nil {
			b = bytes.NewBuffer(body[0])
		}

		r, err = c.httpClient.Post(u, "application/json", b)
		defer r.Body.Close()
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrInvalidHttpMethod
	}

	defer r.Body.Close()

	if r.StatusCode != 200 {
		return nil, ErrHttpError(r)
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *Client) Get(url string) ([]byte, error) {
	return c.call(url, get)
}

func (c *Client) Post(url string, body []byte) ([]byte, error) {
	return c.call(url, post, body)
}

package kap

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/imroc/req/v3"
)

const (
	baseURL = "https://www.kap.org.tr"
)

type client struct {
	c *req.Client
}

type Result struct {
	Error error
	Body  io.Reader
}

func newClient() *client {
	c := req.NewClient().
		SetBaseURL(baseURL).
		EnableForceHTTP1()
	return &client{
		c: c,
	}
}

func newResult(err error, body io.Reader) *Result {
	return &Result{
		Error: err,
		Body:  body,
	}
}

func (s *client) fetch(ctx context.Context, method, url string, reqBody, resBody any) *Result {
	res, err := s.c.R().
		SetContext(ctx).
		SetContentType("application/json").
		SetBody(reqBody).
		SetSuccessResult(resBody).
		Send(method, url)
	if err != nil {
		return newResult(err, nil)
	}
	if res.StatusCode != http.StatusOK {
		return newResult(errors.New(res.Status), nil)
	}

	return newResult(nil, res.Body)
}

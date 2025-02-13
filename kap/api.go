package kap

import (
	"context"

	"github.com/guneyin/gobist/store"
)

type API struct {
	scraper *scraper
}

func NewAPI(_ store.Store) *API {
	return &API{
		scraper: newScraper(),
	}
}

func (a *API) GetCompany(ctx context.Context, code string, opts ...OptionFunc) (*Company, error) {
	cmp := &Company{Code: code}
	err := a.scraper.syncCompany(ctx, cmp)
	if err != nil {
		return nil, err
	}

	err = a.scraper.syncCompanyDetail(ctx, cmp, opts...)
	if err != nil {
		return nil, err
	}

	return cmp, nil
}

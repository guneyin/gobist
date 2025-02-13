package kap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type scraper struct {
	c   *client
	opt *Option
}

func newScraper() *scraper {
	return &scraper{
		c:   newClient(),
		opt: &Option{},
	}
}

func (s *scraper) applyOptions(opts ...OptionFunc) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *scraper) syncCompany(ctx context.Context, cmp *Company) error {
	req := SymbolRequest{
		Keyword:   cmp.Code,
		DiscClass: "ALL",
		Lang:      "tr",
		Channel:   "WEB",
	}
	res := make([]SymbolResponse, 0)
	sri := SymbolResultItem{}
	keys := []string{"combined", "smart"}

loop:
	for _, key := range keys {
		url := fmt.Sprintf("/kapsrc/%s", key)
		sr := s.c.fetch(ctx, http.MethodPost, url, req, &res)
		if sr.Error != nil {
			return sr.Error
		}

		for _, cr := range res {
			for _, result := range cr.Results {
				sliced := strings.Split(result.CmpOrFundCode, ",")
				for _, symbol := range sliced {
					if strings.EqualFold(symbol, cmp.Code) {
						sri = result
						break loop
					}
				}
			}
		}
	}

	if sri.MemberOrFundOid == "" {
		return errors.New("company not found")
	}

	cmp.MemberID = sri.MemberOrFundOid

	return nil
}

func (s *scraper) syncCompanyDetail(ctx context.Context, cmp *Company, opts ...OptionFunc) error {
	s.applyOptions(opts...)

	url := fmt.Sprintf("/tr/sirket-bilgileri/ozet/%s", cmp.MemberID)

	sr := s.c.fetch(ctx, http.MethodGet, url, nil, nil)
	if sr.Error != nil {
		return sr.Error
	}

	doc, err := goquery.NewDocumentFromReader(sr.Body)
	if err != nil {
		return err
	}

	selector := ".w-clearfix.w-inline-block.a-table-row.infoRow"
	list := doc.Find(selector)
	list.Each(func(i int, sel *goquery.Selection) {
		val := strings.TrimSpace(sel.Find("div:nth-child(2)").Text())
		switch i {
		case 0:
			cmp.Address = val
		case 1:
			cmp.Email = val
		case 2:
			cmp.Website = val
		case 5:
			cmp.Index = val
		case 6:
			cmp.Sector = val
		case 7:
			cmp.Market = val
		}

		if s.opt.withShares {
			err = s.syncCompanyShares(ctx, cmp)
			if err != nil {
				log.Println(err)
			}
		}
	})

	return nil
}

func (s *scraper) syncCompanyShares(ctx context.Context, cmp *Company) error {
	url := fmt.Sprintf("tr/infoHistory/kpy41_acc5_sermayede_dogrudan/%s", cmp.MemberID)

	sr := s.c.fetch(ctx, http.MethodGet, url, nil, nil)
	if sr.Error != nil {
		return sr.Error
	}

	doc, err := goquery.NewDocumentFromReader(sr.Body)
	if err != nil {
		return err
	}

	var shareDate time.Time
	selector := ".modal-info.my-modal-info > div > a"
	list := doc.Find(selector)
	list.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}

		if dt, ok := isDate(s); ok {
			shareDate = *dt
			return
		}

		diff := time.Now().Year() - shareDate.Year()
		if diff > 2 {
			return
		}

		if cs, ok := parseLineAsCompanyShareHolder(s); ok {
			cmp.AddShareHolder(shareDate, *cs)
		}
	})

	return nil
}

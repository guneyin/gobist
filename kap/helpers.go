package kap

import (
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func isDate(sel *goquery.Selection) (*time.Time, bool) {
	dtStr := asString(sel)
	dt, err := time.Parse("02/01/2006", dtStr)
	if err != nil {
		return nil, false
	}

	return &dt, true
}

func parseLineAsCompanyShareHolder(sel *goquery.Selection) (*CompanyShareHolder, bool) {
	cs := &CompanyShareHolder{}
	cs.Title = asString(sel.Find("div:nth-child(1)"))
	cs.CapitalByAmount = asFloat(sel.Find("div:nth-child(2)"))
	cs.CapitalByVolume = asFloat(sel.Find("div:nth-child(3)"))
	cs.VoteRight = asFloat(sel.Find("div:nth-child(4)"))

	return cs, true
}

func asString(sel *goquery.Selection) string {
	return strings.TrimSpace(sel.Text())
}

func asFloat(sel *goquery.Selection) float64 {
	s := asString(sel)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

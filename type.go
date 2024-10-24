package gobist

import (
	"encoding/json"
	"fmt"
	"time"
)

type Quote struct {
	Symbol  string   `json:"symbol"`
	Name    string   `json:"name"`
	Price   float64  `json:"price"`
	History *History `json:"history,omitempty"`
}

type History struct {
	Date   *time.Time `json:"date,omitempty"`
	Price  float64    `json:"price,omitempty"`
	Change *Change    `json:"change,omitempty"`
}

type Change struct {
	ByRatio  float64 `json:"byRatio"`
	ByAmount float64 `json:"byAmount"`
}

type Symbol struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type SymbolList struct {
	Count int      `json:"count"`
	Items []Symbol `json:"items"`
}

type QuoteList struct {
	Count int     `json:"count"`
	Items []Quote `json:"items"`
}

func (q Quote) ToJson() string {
	d, _ := json.MarshalIndent(q, "", "  ")
	return string(d)
}

func (s *SymbolList) fromDTO(d *symbolListResponse) *SymbolList {
	if d == nil {
		return s
	}

	s.Count = d.TotalCount
	s.Items = make([]Symbol, d.TotalCount)

	for i, v := range d.Data {
		s.Items[i] = parseSymbolData(i, v.D)
	}

	return s
}

func (ql QuoteList) ToJson() string {
	d, _ := json.MarshalIndent(ql, "", "  ")
	return string(d)
}

func parseSymbolData(i int, d []string) Symbol {
	s := Symbol{
		Id: i,
	}

	if len(d) != 3 {
		return s
	}

	imgUrl := ""
	if d[2] != "" {
		imgUrl = fmt.Sprintf("https://s3-symbol-logo.tradingview.com/%s.svg", d[2])
	}

	s.Code = d[0]
	s.Name = d[1]
	s.Icon = imgUrl

	return s
}

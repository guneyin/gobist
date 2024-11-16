package gobist

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type Quote struct {
	Symbol  string  `json:"symbol"`
	Name    string  `json:"name"`
	Price   string  `json:"price"`
	History History `json:"history,omitempty"`
	Error   string  `json:"error,omitempty"`
}

func (q *Quote) ToJson() string {
	d, _ := json.MarshalIndent(q, "", "  ")
	return string(d)
}

func (q *Quote) SetError(err error) {
	q.Error = err.Error()
}

type History struct {
	Begin  HistoryData   `json:"begin,omitempty"`
	End    HistoryData   `json:"end,omitempty"`
	Change HistoryChange `json:"change,omitempty"`
}

func (h *History) SetBegin(dt time.Time, price float64) {
	h.Begin = HistoryData{dt.Format(time.DateOnly), decimal.NewFromFloat(price).Truncate(2).String()}
}

func (h *History) SetEnd(dt time.Time, price float64) {
	h.End = HistoryData{dt.Format(time.DateOnly), decimal.NewFromFloat(price).Truncate(2).String()}
}

func (h *History) IsValid() bool {
	return h.Begin.Date != "" && h.End.Date != ""
}

type HistoryData struct {
	Date  string `json:"date,omitempty"`
	Price string `json:"price,omitempty"`
}

type HistoryChange struct {
	ByRatio  string `json:"byRatio"`
	ByAmount string `json:"byAmount"`
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

type QuoteList struct {
	Count int     `json:"count"`
	Items []Quote `json:"items"`
}

func (ql QuoteList) ToJson() string {
	d, _ := json.MarshalIndent(ql, "", "  ")
	return string(d)
}

func parseSymbolData(i int, d []string) Symbol {
	s := Symbol{
		Id: i + 1,
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

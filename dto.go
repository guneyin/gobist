package gobist

import "time"

type quoteDTO struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				ExchangeName         string  `json:"exchangeName"`
				FullExchangeName     string  `json:"fullExchangeName"`
				InstrumentType       string  `json:"instrumentType"`
				FirstTradeDate       int     `json:"firstTradeDate"`
				RegularMarketTime    int     `json:"regularMarketTime"`
				HasPrePostMarketData bool    `json:"hasPrePostMarketData"`
				Gmtoffset            int     `json:"gmtoffset"`
				Timezone             string  `json:"timezone"`
				ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				FiftyTwoWeekHigh     float64 `json:"fiftyTwoWeekHigh"`
				FiftyTwoWeekLow      float64 `json:"fiftyTwoWeekLow"`
				RegularMarketDayHigh float64 `json:"regularMarketDayHigh"`
				RegularMarketDayLow  float64 `json:"regularMarketDayLow"`
				RegularMarketVolume  int     `json:"regularMarketVolume"`
				LongName             string  `json:"longName"`
				ShortName            string  `json:"shortName"`
				ChartPreviousClose   float64 `json:"chartPreviousClose"`
				PriceHint            int     `json:"priceHint"`
				CurrentTradingPeriod struct {
					Pre struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"pre"`
					Regular struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"regular"`
					Post struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"post"`
				} `json:"currentTradingPeriod"`
				DataGranularity string   `json:"dataGranularity"`
				Range           string   `json:"range"`
				ValidRanges     []string `json:"validRanges"`
			} `json:"meta"`
			Timestamp  []int `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Volume []int     `json:"volume"`
					High   []float64 `json:"high"`
					Open   []float64 `json:"open"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
				} `json:"quoteDTO"`
				Adjclose []struct {
					Adjclose []float64 `json:"adjclose"`
				} `json:"adjclose"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

func (q quoteDTO) adjCloseCheck() bool {
	if len(q.Chart.Result[0].Indicators.Adjclose) == 0 {
		return false
	}

	if len(q.Chart.Result[0].Indicators.Adjclose[0].Adjclose) == 0 {
		return false
	}

	return true
}

func (q quoteDTO) getClosePrice() (time.Time, float64) {
	tsSlice := q.Chart.Result[0].Timestamp
	closeSlice := q.Chart.Result[0].Indicators.Adjclose[0].Adjclose

	ts := tsSlice[len(tsSlice)-1]
	cp := closeSlice[len(closeSlice)-1]

	return time.Unix(int64(ts), 0), cp
}

type symbolListResponse struct {
	TotalCount int `json:"totalCount"`
	Data       []struct {
		S string   `json:"s"`
		D []string `json:"d"`
	} `json:"data"`
	Params struct {
		Turkey struct {
			Symbols struct {
				Query struct {
					Types []string `json:"types"`
				} `json:"query"`
			} `json:"symbols"`
			Filter []struct {
				Left      string      `json:"left"`
				Operation string      `json:"operation"`
				Right     interface{} `json:"right"`
			} `json:"filter"`
			Sort struct {
				SortBy     string `json:"sortBy"`
				SortOrder  string `json:"sortOrder"`
				NullsFirst bool   `json:"nullsFirst"`
			} `json:"sort"`
			Options struct {
				Lang string `json:"lang"`
			} `json:"options"`
		} `json:"turkey"`
	} `json:"params"`
}

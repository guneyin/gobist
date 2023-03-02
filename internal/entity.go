package internal

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strconv"
	"time"
)

type yahooQuote struct {
	OptionChain struct {
		Result []struct {
			UnderlyingSymbol string        `json:"underlyingSymbol"`
			ExpirationDates  []interface{} `json:"expirationDates"`
			Strikes          []interface{} `json:"strikes"`
			HasMiniOptions   bool          `json:"hasMiniOptions"`
			Quote            struct {
				Language                          string  `json:"language"`
				Region                            string  `json:"region"`
				QuoteType                         string  `json:"quoteType"`
				TypeDisp                          string  `json:"typeDisp"`
				QuoteSourceName                   string  `json:"quoteSourceName"`
				Triggerable                       bool    `json:"triggerable"`
				CustomPriceAlertConfidence        string  `json:"customPriceAlertConfidence"`
				Currency                          string  `json:"currency"`
				RegularMarketChangePercent        float64 `json:"regularMarketChangePercent"`
				RegularMarketPrice                float64 `json:"regularMarketPrice"`
				FirstTradeDateMilliseconds        int64   `json:"firstTradeDateMilliseconds"`
				PriceHint                         int     `json:"priceHint"`
				MarketState                       string  `json:"marketState"`
				Tradeable                         bool    `json:"tradeable"`
				CryptoTradeable                   bool    `json:"cryptoTradeable"`
				Exchange                          string  `json:"exchange"`
				ShortName                         string  `json:"shortName"`
				LongName                          string  `json:"longName"`
				MessageBoardId                    string  `json:"messageBoardId"`
				ExchangeTimezoneName              string  `json:"exchangeTimezoneName"`
				ExchangeTimezoneShortName         string  `json:"exchangeTimezoneShortName"`
				GmtOffSetMilliseconds             int     `json:"gmtOffSetMilliseconds"`
				Market                            string  `json:"market"`
				EsgPopulated                      bool    `json:"esgPopulated"`
				RegularMarketChange               float64 `json:"regularMarketChange"`
				RegularMarketTime                 int     `json:"regularMarketTime"`
				RegularMarketDayHigh              float64 `json:"regularMarketDayHigh"`
				RegularMarketDayRange             string  `json:"regularMarketDayRange"`
				RegularMarketDayLow               float64 `json:"regularMarketDayLow"`
				RegularMarketVolume               int     `json:"regularMarketVolume"`
				RegularMarketPreviousClose        float64 `json:"regularMarketPreviousClose"`
				Bid                               float64 `json:"bid"`
				Ask                               float64 `json:"ask"`
				BidSize                           int     `json:"bidSize"`
				AskSize                           int     `json:"askSize"`
				FullExchangeName                  string  `json:"fullExchangeName"`
				FinancialCurrency                 string  `json:"financialCurrency"`
				RegularMarketOpen                 float64 `json:"regularMarketOpen"`
				AverageDailyVolume3Month          int     `json:"averageDailyVolume3Month"`
				AverageDailyVolume10Day           int     `json:"averageDailyVolume10Day"`
				FiftyTwoWeekLowChange             float64 `json:"fiftyTwoWeekLowChange"`
				FiftyTwoWeekLowChangePercent      float64 `json:"fiftyTwoWeekLowChangePercent"`
				FiftyTwoWeekRange                 string  `json:"fiftyTwoWeekRange"`
				FiftyTwoWeekHighChange            float64 `json:"fiftyTwoWeekHighChange"`
				FiftyTwoWeekHighChangePercent     float64 `json:"fiftyTwoWeekHighChangePercent"`
				FiftyTwoWeekLow                   float64 `json:"fiftyTwoWeekLow"`
				FiftyTwoWeekHigh                  float64 `json:"fiftyTwoWeekHigh"`
				EarningsTimestamp                 int     `json:"earningsTimestamp"`
				EarningsTimestampStart            int     `json:"earningsTimestampStart"`
				EarningsTimestampEnd              int     `json:"earningsTimestampEnd"`
				TrailingAnnualDividendRate        float64 `json:"trailingAnnualDividendRate"`
				TrailingPE                        float64 `json:"trailingPE"`
				TrailingAnnualDividendYield       float64 `json:"trailingAnnualDividendYield"`
				EpsTrailingTwelveMonths           float64 `json:"epsTrailingTwelveMonths"`
				EpsForward                        float64 `json:"epsForward"`
				EpsCurrentYear                    float64 `json:"epsCurrentYear"`
				PriceEpsCurrentYear               float64 `json:"priceEpsCurrentYear"`
				SharesOutstanding                 int     `json:"sharesOutstanding"`
				BookValue                         float64 `json:"bookValue"`
				FiftyDayAverage                   float64 `json:"fiftyDayAverage"`
				FiftyDayAverageChange             float64 `json:"fiftyDayAverageChange"`
				FiftyDayAverageChangePercent      float64 `json:"fiftyDayAverageChangePercent"`
				TwoHundredDayAverage              float64 `json:"twoHundredDayAverage"`
				TwoHundredDayAverageChange        float64 `json:"twoHundredDayAverageChange"`
				TwoHundredDayAverageChangePercent float64 `json:"twoHundredDayAverageChangePercent"`
				MarketCap                         int64   `json:"marketCap"`
				ForwardPE                         float64 `json:"forwardPE"`
				PriceToBook                       float64 `json:"priceToBook"`
				SourceInterval                    int     `json:"sourceInterval"`
				ExchangeDataDelayedBy             int     `json:"exchangeDataDelayedBy"`
				AverageAnalystRating              string  `json:"averageAnalystRating"`
				Symbol                            string  `json:"symbol"`
			} `json:"quote"`
			Options []interface{} `json:"options"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"optionChain"`
}

type yahooHistory struct {
	Date     time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
	AdjClose float64
	Volume   int64
}

func (q *yahooQuote) Unmarshall(d []byte) {
	_ = json.Unmarshal(d, q)
}

func (h *yahooHistory) Parse(d []byte) {
	r := csv.NewReader(bytes.NewReader(d))
	r.Comma = ','

	record, err := r.ReadAll()
	if err != nil {
		return
	}

	if len(record) > 1 {
		h.Date, _ = time.Parse("2006-01-02", record[1][0])
		h.Open, _ = strconv.ParseFloat(record[1][1], 64)
		h.High, _ = strconv.ParseFloat(record[1][2], 64)
		h.Low, _ = strconv.ParseFloat(record[1][3], 64)
		h.Close, _ = strconv.ParseFloat(record[1][4], 64)
		h.AdjClose, _ = strconv.ParseFloat(record[1][5], 64)
		h.Volume, _ = strconv.ParseInt(record[1][6], 10, 64)
	}
}

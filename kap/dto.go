package kap

type SymbolRequest struct {
	Keyword   string `json:"keyword"`
	DiscClass string `json:"discClass"`
	Lang      string `json:"lang"`
	Channel   string `json:"channel"`
}

type SymbolResponse struct {
	Category string             `json:"category"`
	Results  []SymbolResultItem `json:"results"`
}

type SymbolResultItem struct {
	SearchValue     string      `json:"searchValue"`
	SearchType      string      `json:"searchType"`
	ActionKey       string      `json:"actionKey"`
	MemberOrFundOid string      `json:"memberOrFundOid"`
	SubjectOid      interface{} `json:"subjectOid"`
	MarketOid       interface{} `json:"marketOid"`
	DiscType        interface{} `json:"discType"`
	Year            int         `json:"year"`
	Period          int         `json:"period"`
	CmpOrFundCode   string      `json:"cmpOrFundCode"`
	IndexList       interface{} `json:"indexList"`
}

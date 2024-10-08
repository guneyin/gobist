package gobist

import (
	"time"
)

type Quote struct {
	Symbol  string   `json:"symbol,omitempty"`
	Name    string   `json:"name,omitempty"`
	Price   float64  `json:"price,omitempty"`
	History *History `json:"history,omitempty"`
}

type History struct {
	Date  *time.Time `json:"date,omitempty"`
	Price float64    `json:"price,omitempty"`
}

package internal

import (
	"time"
)

type Quote struct {
	Symbol  string
	Name    string
	Price   float64
	History *History
}

type History struct {
	Date  time.Time
	Price float64
}

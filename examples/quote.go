package main

import (
	"fmt"
	"log"
	"time"

	"github.com/guneyin/gobist"
)

func main() {
	bist := gobist.New()

	t, _ := time.Parse(time.DateOnly, "2024-09-25")
	q, err := bist.GetQuoteWithHistory([]string{"TUPRS", "BIMAS", "VESBE", "THYAO"}, t)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range q.Items {
		fmt.Println(fmt.Sprintf("symbol=%s name=%s price=%f, history_price=%f, change=%f", item.Symbol, item.Name, item.Price, item.History.Price, item.History.Change.ByRatio))
	}
}

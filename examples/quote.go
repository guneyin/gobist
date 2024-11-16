package main

import (
	"fmt"
	"log"
	"time"

	"github.com/guneyin/gobist"
)

func main() {
	bist := gobist.New()

	tBegin, _ := time.Parse(time.DateOnly, "2023-10-16")
	tEnd, _ := time.Parse(time.DateOnly, "2024-10-15")

	q, err := bist.GetQuoteWithHistory([]string{"TUPRS", "BIMAS", "VESBE", "THYAO"}, tBegin, tEnd)
	if err != nil {
		log.Fatal(err)
	}

	const layout = "%-10s %-30s %-20s %-20s %-20s %-15s %-30s"
	fmt.Println(fmt.Sprintf(layout, "Symbol", "Name", "Current Price", "History Begin", "History End", "Change (%)", "Error"))
	for _, item := range q.Items {
		fmt.Println(fmt.Sprintf(layout,
			item.Symbol, item.Name, item.Price, item.History.Begin.Price,
			item.History.End.Price, item.History.Change.ByRatio, item.Error))
	}
}

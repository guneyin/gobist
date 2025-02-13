package main

import (
	"fmt"
	"log"
	"time"

	"github.com/guneyin/gobist"
)

func main() {
	bist := gobist.New()

	tBegin, _ := time.Parse(time.DateOnly, "2024-02-13")
	tEnd, _ := time.Parse(time.DateOnly, "2025-02-13")

	q, err := bist.GetQuote("TUPRS", gobist.WithPeriod(gobist.NewPeriod(tBegin, tEnd)))
	if err != nil {
		log.Fatal(err)
	}

	const layout = "%-10s %-30s %-20s %-20s %-20s %-15s %-30s"
	fmt.Println(fmt.Sprintf(layout, "Symbol", "Name", "Current Price", "History Begin", "History End", "Change (%)", "Error"))
	fmt.Println(fmt.Sprintf(layout,
		q.Symbol, q.Name, q.Price, q.History.Begin.Price,
		q.History.End.Price, q.History.Change.ByRatio, q.Error))
}

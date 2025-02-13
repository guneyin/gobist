package main

import (
	"context"
	"fmt"
	"github.com/guneyin/gobist/quote"
	"log"
	"time"

	"github.com/guneyin/gobist"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bist := gobist.New()

	tBegin, _ := time.Parse(time.DateOnly, "2024-02-13")
	tEnd, _ := time.Parse(time.DateOnly, "2025-02-13")

	q, err := bist.GetQuote(ctx, "TUPRS", quote.WithPeriod(quote.NewPeriod(tBegin, tEnd)))
	if err != nil {
		log.Fatal(err)
	}

	const layout = "%-10s %-30s %-20s %-20s %-20s %-15s %-30s"
	fmt.Println(fmt.Sprintf(layout, "Symbol", "Name", "Current Price", "History Begin", "History End", "Change (%)", "Error"))
	fmt.Println(fmt.Sprintf(layout,
		q.Symbol, q.Name, q.Price, q.History.Begin.Price,
		q.History.End.Price, q.History.Change.ByRatio, q.Error))
}

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

	fmt.Println(fmt.Sprintf("%-10s %-30s %-20s %-20s %-20s %-15s %-30s", "Symbol", "Name", "Current Price", "History Begin", "History End", "Change", "Error"))
	for _, item := range q.Items {
		fmt.Println(fmt.Sprintf("%-10s %-30s %-20f %-20f %-20f %-15f %-30s",
			item.Symbol, item.Name, item.Price.Float64(), item.History.Begin.Price.Float64(), item.History.End.Price.Float64(), item.History.Change.ByRatio.Float64(), item.Error))
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/guneyin/gobist"
)

func main() {
	bist, err := gobist.New()
	if err != nil {
		log.Fatal(err)
	}

	q, err := bist.GetQuote("TUPRS")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("symbol=%s name=%s price=%f", q.Symbol, q.Name, q.Price))
}

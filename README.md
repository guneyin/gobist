
# gobist - GO library for BIST (Borsa Istanbul)

This project aims to provide some useful tools to fetch stock data for BIST via Yahoo Finance API


## Installation

    $ go get github.com/guneyin/gobist

## Example
```go
func main() {
    bist := gobist.New()
    
    tBegin, _ := time.Parse(time.DateOnly, "2023-09-25")
    tEnd, _ := time.Parse(time.DateOnly, "2024-10-25")
    
    q, err := bist.GetQuoteWithHistory([]string{"TUPRS", "BIMAS", "VESBE", "THYAO"}, tBegin, tEnd)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(fmt.Sprintf("%-10s %-30s %-20s %-20s %-20s %-15s %-30s", "Symbol", "Name", "Current Price", "History Begin", "History End", "Change", "Error"))
    for _, item := range q.Items {
        fmt.Println(fmt.Sprintf("%-10s %-30s %-20f %-20f %-20f %-15f %-30s",
        item.Symbol, item.Name, item.Price, item.History.Begin.Price, item.History.End.Price, item.History.Change.ByRatio, item.Error))
    }
}
``` 

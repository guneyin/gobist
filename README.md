
# gobist - GO library for BIST (Borsa Istanbul)

This project aims to provide some useful tools to fetch stock data for BIST via Yahoo Finance API


## Installation

    $ go get github.com/guneyin/gobist

## Usage and Example

### Example
```go
func main() {
    bist := gobist.New()
    
    dt, _ := time.Parse(time.DateOnly, "2024-09-25")
    q, err := bist.GetQuoteWithHistory([]string{"TUPRS", "BIMAS", "VESBE", "THYAO"}, dt)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, item := range q.Items {
        fmt.Println(fmt.Sprintf("symbol=%s name=%s price=%f, history_price=%f, change=%f", item.Symbol, item.Name, item.Price, item.History.Price, item.History.Change.ByRatio))
    }
}
``` 

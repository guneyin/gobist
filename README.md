
# gobist - GO library for BIST (Borsa Istanbul)

This project aims to provide some useful tools to fetch stock data for BIST via Yahoo Finance API


## Installation

    $ go get github.com/guneyin/gobist

## Usage and Example

### Create Client
```go
bist, err := gobist.New()
if err != nil {
    log.Fatal(err)
}
```

### Get Quote
```go
q, err := bist.GetQuote("TUPRS")
if err != nil {
    log.Fatal(err)
}
```

### Example
```go
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
``` 
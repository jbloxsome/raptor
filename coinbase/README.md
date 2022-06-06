# Raptor Coinbase
Use this package to manage local level 2 orderbooks for Coinbase Pro.

## Usage
Import the package:
```go
import "github.com/jbloxsome/raptor/coinbase"
```

Create an orderbook:
```go
ob, err := coinbase.NewCoinbase("BTC-USD")
if err != nil {
    panic(err)
}
```

The the realtime orderbook is exposed over a channel:
```go
for orderbook := range ob.Orderbook {
    maxBid := orderbook.Bids.Max()
    minAsk := orderbook.Asks.Min()

    midMarket := (maxBid + minAsk) / 2
    log.Printf("BTC-USD: %f", midMarket)
}
```

To list the price levels:
```go
for orderbook := range ob.Orderbook {
    bids := orderbook.Bids.Keys(orderbook.Bids.Min(), orderbook.Bids.Max())
    asks := orderbook.Asks.Keys(orderbook.Asks.Min(), orderbook.Asks.Max())
    log.Printf("Bids: %f", bids)
    log.Printf("Asks: %f", asks)
}
```
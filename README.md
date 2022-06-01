# Raptor
Raptor provides a consistent Golang library for retrieving orderbook data from Crypto exchanges in real-time.

## Supported Exchanges
- Coinbase
- FTX

## Performance
The orderbook data is stored within a Red Black Tree. The code for the RBT can be found under /rbbst. Benchmarks for a tree with 1 million nodes resulted in the following:

- Insert - 210 ns/op
- Get    - 115 ns/op

## Example Usage
Orderbooks are streamed over a channel:

```go
package main

import (
	"os"
	"log"

	"github.com/jbloxsome/raptor/coinbase"
)

func main() {
	interrupt := make(chan os.Signal, 1)

	coinbase, err := coinbase.NewCoinbase("BTC-USD")
	if err != nil {
		panic(err)
	}

	go func() {
		for orderbook := range coinbase.Orderbook {
			maxBid := orderbook.Bids.Max()
			minAsk := orderbook.Asks.Min()
			midPrice := (maxBid + minAsk) / 2
			log.Printf("Coinbase BTC/USD mid price - %f", midPrice)
		}
	}()

	for {
		select {
		case <-interrupt:
			coinbase.Close()
			return
		}
	}
}
```

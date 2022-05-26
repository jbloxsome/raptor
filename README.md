# Raptor
Raptor provides a consistent Golang API for retrieving orderbook data from Crypto exchanges in real-time.

## Supported Exchanges
- Coinbase Pro

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
			log.Println(orderbook) 
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
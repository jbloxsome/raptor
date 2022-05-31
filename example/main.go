package main

import (
	"os"
	"log"

	"github.com/jbloxsome/raptor/coinbase"
	// "github.com/jbloxsome/raptor/ftx"
)

func main() {
	interrupt := make(chan os.Signal, 1)

	// btc_usd, err := ftx.NewFTX("BTC/USD")
	// if err != nil {
	// 	panic(err)
	// }

	btc_usd, err := coinbase.NewCoinbase("BTC-USD")
	if err != nil {
		panic(err)
	}

	// eth_usd, err := coinbase.NewCoinbase("ETH-USD")
	// if err != nil {
	// 	panic(err)
	// }

	// btc_eth, err := coinbase.NewCoinbase("ETH-BTC")
	// if err != nil {
	// 	panic(err)
	// }

	go func() {
		for orderbook := range btc_usd.Orderbook {
			// maxBid := orderbook.Bids.Max()
			// secondBid := orderbook.Bids.Size()
			// minAsk := orderbook.Asks.Min()
			prices := orderbook.Bids.Keys(31500, orderbook.Bids.Max())

			// midMarket := (maxBid + minAsk) / 2
			log.Printf("BTC-USD: %f", prices) 
			// log.Printf("First: %f", maxBid)
			// log.Printf("Second: %v", secondBid)
		}
	}()

	// go func() {
	// 	for orderbook := range eth_usd.Orderbook {
	// 		log.Println("ETH-USD: " + orderbook.Bids[0][0]) 
	// 	}
	// }()

	// go func() {
	// 	for orderbook := range btc_eth.Orderbook {
	// 		log.Println("ETH-BTC: " + orderbook.Bids[0][0]) 
	// 	}
	// }()

	for {
		select {
		case <-interrupt:
			btc_usd.Close()
			// eth_usd.Close()
			// btc_eth.Close()
			return
		}
	}
}
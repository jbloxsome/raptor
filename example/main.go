package main

import (
	"os"
	"log"
	"strconv"

	"github.com/jbloxsome/raptor/coinbase"
)

func main() {
	interrupt := make(chan os.Signal, 1)

	btc_usd, err := coinbase.NewCoinbase("BTC-USD")
	if err != nil {
		panic(err)
	}

	eth_usd, err := coinbase.NewCoinbase("ETH-USD")
	if err != nil {
		panic(err)
	}

	btc_eth, err := coinbase.NewCoinbase("ETH-BTC")
	if err != nil {
		panic(err)
	}

	go func() {
		for orderbook := range btc_usd.Orderbook {
			bestBid, err := strconv.ParseFloat(orderbook.Bids[0][0], 32)
			if err != nil {
				panic(err)
			}
			// bestAsk, err := strconv.ParseFloat(orderbook.Asks[0][0], 32)
			// if err != nil {
			// 	panic(err)
			// }
			// midMarket := (bestBid + bestAsk) / 2
			log.Printf("BTC-USD: %f", bestBid) 
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
			eth_usd.Close()
			btc_eth.Close()
			return
		}
	}
}
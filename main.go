package main

import (
	"os"
	"log"
)

func main() {
	interrupt := make(chan os.Signal, 1)

	coinbase, err := NewOrderbook("coinbase", "BTC-USD")
	if err != nil {
		panic(err)
	}

	go func() {
		for prices := range coinbase.Prices {
			log.Println(prices) 
		}
	}()

	go func() {
		for message := range coinbase.Message {
			log.Println(message)
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
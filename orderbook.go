package main

import (
	"log"
	"net/url"
	"errors"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string
	ProductId string
	Bids [][]string
	Asks [][]string
	Time string
	Changes [][]string
}

type Prices struct {
	Bids [][]string
	Asks [][]string
}

type Orderbook struct {
	Exchange string
	Pair string
	Prices chan *Prices
	Message chan string
	Connection *websocket.Conn
}

func NewOrderbook(exchange string, pair string) (*Orderbook, error) {
	if exchange != "coinbase" {
		return nil, errors.New("exchange not currently supported")
	}

	if pair != "BTC-USD" {
		return nil, errors.New("pair not currently supported for exchange")
	}

	// Open websocket connection
	u := url.URL{Scheme: "wss", Host: "ws-feed.exchange.coinbase.com", Path: "/",}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Subscribe to l2 feed
	err = c.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"subscribe\",\"product_ids\":[\"BTC-USD\"],\"channels\":[\"level2\"]}"))
	if err != nil {
		return nil, err
	}

	messages := make(chan string)
	prices := make(chan *Prices)
	var currentPrices Prices

	go func() {
		for {
			var asMessage Message 
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			err = json.Unmarshal(message, &asMessage)
			if err != nil {
				log.Println("read:", err)
				return
			}

			if asMessage.Type == "snapshot" {
				currentPrices.Bids = asMessage.Bids
				currentPrices.Asks = asMessage.Asks
				prices <-&currentPrices
			} else if asMessage.Type == "l2update" {
				
				for _, change := range asMessage.Changes {
					if change[0] == "buy" {
						for idx, bid := range currentPrices.Bids {
							if bid[0] == change[1] {
								currentPrices.Bids[idx] = []string{change[1], change[2]}
								break
							}
						}
					}

					if change[0] == "sell" {
						for idx, ask := range currentPrices.Asks {
							if ask[0] == change[1] {
								currentPrices.Asks[idx] = []string{change[1], change[2]}
								break
							}
						}
					}
				}

				prices <-&Prices{
					Bids: currentPrices.Bids,
					Asks: currentPrices.Asks,
				}
			} else {
				messages <-string(message)
			}	
		}
	}()

	return &Orderbook{
		Exchange: exchange,
		Pair: pair,
		Prices: prices,
		Message: messages,
		Connection: c,
	}, nil
}

func (o *Orderbook) Close() {
	o.Connection.Close()
	close(o.Prices)
	close(o.Message)
	return
}
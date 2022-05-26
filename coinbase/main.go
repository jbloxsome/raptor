package coinbase

import (
	"log"
	"sort"
	"net/url"
	"encoding/json"

	"github.com/gorilla/websocket"
	ob "github.com/jbloxsome/raptor/orderbook"
)

// Struct for deserialising messages received over the Coinbase websocket
type Message struct {
	Type string
	ProductId string
	Bids [][]string
	Asks [][]string
	Time string
	Changes [][]string
}

type Coinbase struct {
	Pair string
	Orderbook chan *ob.Orderbook
	Connection *websocket.Conn
}

func NewCoinbase(pair string) (*Coinbase, error) {
	// Open websocket connection
	u := url.URL{Scheme: "wss", Host: "ws-feed.exchange.coinbase.com", Path: "/",}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Subscribe to l2 feed
	err = c.WriteMessage(websocket.TextMessage, []byte("{\"type\":\"subscribe\",\"product_ids\":[\"" + pair + "\"],\"channels\":[\"level2\"]}"))
	if err != nil {
		return nil, err
	}

	orderbook := make(chan *ob.Orderbook)
	currentOrderbook := ob.NewOrderbook()

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
				for bid := range asMessage.Bids {
					currentOrderbook.AddBidLevel(bid[0], bid[1])
				}
				
				for ask := range asMessage.Asks {
					currentOrderbook.AddAskLevel(ask[0], ask[1])
				}

				orderbook <-&currentOrderbook
			} else if asMessage.Type == "l2update" {
				
				for _, change := range asMessage.Changes {
					if change[0] == "buy" {
						if change[2] != "0" {
							// Add or update bid level in orderbook
							level := currentOrderbook.GetBidLevel(change[1])
							if level != nil {
								currentOrderbook.RemoveBidLevel(change[1])
							}

							currentOrderbook.AddBidLevel(change[1], change[2])
						} else {
							// Remove bid level from orderbook
							currentOrderbook.RemoveBidLevel(change[1])
						}
					}

					if change[0] == "sell" {
						if change[2] != "0" {
							// Add or update ask level in orderbook
							level := currentOrderbook.GetAskLevel(change[1])
							if level != nil {
								currentOrderbook.RemoveAskLevel(change[1])
							}

							currentOrderbook.AddAskLevel(change[1], change[2])
						} else {
							// Remove ask level from orderbook
							currentOrderbook.RemoveAskLevel(change[1])
						}
					}
				}

				orderbook <-currentOrderbook
			}	
		}
	}()

	return &Coinbase{
		Pair: pair,
		Orderbook: orderbook,
		Connection: c,
	}, nil
}

func (c *Coinbase) Close() {
	c.Connection.Close()
	close(c.Orderbook)
	return
}
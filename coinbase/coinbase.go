package coinbase

import (
	"log"
	"strconv"
	"time"
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
				for _, bid := range asMessage.Bids {
					price, err := strconv.ParseFloat(bid[0], 64)
					if err != nil {
						log.Println("read:", err)
						return
					}

					volume, err := strconv.ParseFloat(bid[1], 64)
					if err != nil {
						log.Println("read:", err)
						return
					}

					currentOrderbook.AddBidLevel(price, volume)
				}
				
				for _, ask := range asMessage.Asks {
					price, err := strconv.ParseFloat(ask[0], 64)
					if err != nil {
						log.Println("read:", err)
						return
					}

					volume, err := strconv.ParseFloat(ask[1], 64)
					if err != nil {
						log.Println("read:", err)
						return
					}

					currentOrderbook.AddAskLevel(price, volume)
				}

				orderbook <-&currentOrderbook
			} else if asMessage.Type == "l2update" {
				start := time.Now()
				
				for _, change := range asMessage.Changes {

					price, err := strconv.ParseFloat(change[1], 64)
					if err != nil {
						log.Println("read:", err)
						return
					}

					volume, err := strconv.ParseFloat(change[2], 64)
					if err != nil {
						log.Println("read:", err)
						return
					}

					if change[0] == "buy" {
						if volume > 0.0 {
							// Add or update bid level in orderbook
							level := currentOrderbook.GetBidLevel(price)
							if level != nil {
								currentOrderbook.RemoveBidLevel(price)
							}

							currentOrderbook.AddBidLevel(price, volume)
						} else {
							// Remove bid level from orderbook
							currentOrderbook.RemoveBidLevel(price)
						}
					}

					if change[0] == "sell" {
						if volume > 0.0 {
							// Add or update ask level in orderbook
							level := currentOrderbook.GetAskLevel(price)
							if level != nil {
								currentOrderbook.RemoveAskLevel(price)
							}

							currentOrderbook.AddAskLevel(price, volume)
						} else {
							// Remove ask level from orderbook
							currentOrderbook.RemoveAskLevel(price)
						}
					}
				}

				log.Printf("orderbook update, execution time %s\n", time.Since(start))

				orderbook <-&currentOrderbook
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
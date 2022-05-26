package coinbase

import (
	"log"
	"net/url"
	"encoding/json"

	"github.com/gorilla/websocket"
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

type Orderbook struct {
	Bids [][]string
	Asks [][]string
}

type Coinbase struct {
	Pair string
	Orderbook chan *Orderbook
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

	orderbook := make(chan *Orderbook)
	var currentOrderbook Orderbook

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
				currentOrderbook.Bids = asMessage.Bids
				currentOrderbook.Asks = asMessage.Asks
				orderbook <-&currentOrderbook
			} else if asMessage.Type == "l2update" {
				
				for _, change := range asMessage.Changes {
					if change[0] == "buy" {
						if change[2] != "0" {
							found := false

							for idx, bid := range currentOrderbook.Bids {
								if bid[0] == change[1] {
									currentOrderbook.Bids[idx] = []string{change[1], change[2]}
									found = true
									break
								}
							}

							if found == false {
								currentOrderbook.Bids = append(currentOrderbook.Bids, []string{change[1], change[2]})
								break
							}
						} else {
							for idx, bid := range currentOrderbook.Bids {
								if bid[0] == change[1] {
									currentOrderbook.Bids[idx] = currentOrderbook.Bids[len(currentOrderbook.Bids)-1]
									currentOrderbook.Bids = currentOrderbook.Bids[:len(currentOrderbook.Bids)-1]
									break
								}
							}
						}
					}

					if change[0] == "sell" {
						if change[2] != "0" {
							found := false

							for idx, ask := range currentOrderbook.Asks {
								if ask[0] == change[1] {
									currentOrderbook.Asks[idx] = []string{change[1], change[2]}
									found = true
									break
								}
							}

							if found == false {
								currentOrderbook.Asks = append(currentOrderbook.Asks, []string{change[1], change[2]})
								break
							}
						} else {
							for idx, ask := range currentOrderbook.Asks {
								if ask[0] == change[1] {
									currentOrderbook.Asks[idx] = currentOrderbook.Asks[len(currentOrderbook.Asks)-1]
									currentOrderbook.Asks = currentOrderbook.Asks[:len(currentOrderbook.Asks)-1]
									break
								}
							}
						}
					}
				}

				orderbook <-&Orderbook{
					Bids: currentOrderbook.Bids,
					Asks: currentOrderbook.Asks,
				}
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
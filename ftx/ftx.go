package ftx

import (
	"log"
	"time"
	"net/url"
	"encoding/json"

	"github.com/gorilla/websocket"
	ob "github.com/jbloxsome/raptor/orderbook"
)

type MessageData struct {
	Action string
	Time float64
	Checksum float64
	Bids [][]float64
	Asks [][]float64
}

// Struct for deserialising messages received over the FTX websocket
type Message struct {
	Channel string
	Market string
	Type string
	Data MessageData
}

type FTX struct {
	Pair string
	Orderbook chan *ob.Orderbook
	Connection *websocket.Conn
}

func NewFTX(pair string) (*FTX, error) {
	// Open websocket connection
	u := url.URL{Scheme: "wss", Host: "ftx.com", Path: "/ws/",}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Subscribe to l2 orderbook feed
	err = c.WriteMessage(websocket.TextMessage, []byte("{\"op\":\"subscribe\",\"market\": \"" + pair + "\" ,\"channel\":\"orderbook\"}"))
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

			if asMessage.Data.Action == "partial" {
				for _, bid := range asMessage.Data.Bids {
					currentOrderbook.AddBidLevel(bid[0], bid[1])
				}
				
				for _, ask := range asMessage.Data.Asks {
					currentOrderbook.AddAskLevel(ask[0], ask[1])
				}

				orderbook <-&currentOrderbook
			} else if asMessage.Data.Action == "update" {
				start := time.Now()

				for _, bid := range asMessage.Data.Bids {
					if bid[1] > 0.0 {
						// Add or update bid level in orderbook
						level := currentOrderbook.GetBidLevel(bid[0])
						if level != nil {
							currentOrderbook.RemoveBidLevel(bid[0])
						}

						currentOrderbook.AddBidLevel(bid[0], bid[1])
					} else {
						// Remove bid level from orderbook
						currentOrderbook.RemoveBidLevel(bid[0])
					}
				}

				for _, ask := range asMessage.Data.Asks {
					if ask[1] > 0.0 {
						// Add or update ask level in orderbook
						level := currentOrderbook.GetAskLevel(ask[0])
						if level != nil {
							currentOrderbook.RemoveAskLevel(ask[0])
						}

						currentOrderbook.AddAskLevel(ask[0], ask[1])
					} else {
						// Remove ask level from orderbook
						currentOrderbook.RemoveAskLevel(ask[0])
					}
				}

				log.Printf("orderbook update, execution time %s\n", time.Since(start))

				orderbook <-&currentOrderbook
			}	
		}
	}()

	return &FTX{
		Pair: pair,
		Orderbook: orderbook,
		Connection: c,
	}, nil
}

func (c *FTX) Close() {
	c.Connection.Close()
	close(c.Orderbook)
	return
}
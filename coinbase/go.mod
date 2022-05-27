module github.com/jbloxsome/raptor/coinbase

go 1.18

require (
	github.com/gorilla/websocket v1.5.0
	github.com/jbloxsome/raptor/orderbook v0.0.0-00010101000000-000000000000
)

require github.com/jbloxsome/raptor/rbbst v0.0.0-00010101000000-000000000000 // indirect

replace github.com/jbloxsome/raptor/orderbook => ../orderbook

replace github.com/jbloxsome/raptor/rbbst => ../rbbst

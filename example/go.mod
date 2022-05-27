module github.com/jbloxsome/raptor/example

go 1.18

replace github.com/jbloxsome/raptor/coinbase => ../coinbase

require github.com/jbloxsome/raptor/coinbase v0.0.0-00010101000000-000000000000

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jbloxsome/raptor/orderbook v0.0.0-00010101000000-000000000000 // indirect
	github.com/jbloxsome/raptor/rbbst v0.0.0-00010101000000-000000000000 // indirect
)

replace github.com/jbloxsome/raptor/orderbook => ../orderbook

replace github.com/jbloxsome/raptor/rbbst => ../rbbst

replace github.com/jbloxsome/raptor/ftx => ../ftx

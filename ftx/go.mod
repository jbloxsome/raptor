module github.com/jbloxsome/raptor/ftx

go 1.18

replace github.com/jbloxsome/orderbook => ../orderbook

replace github.com/jbloxsome/raptor/rbbst => ../rbbst

require (
	github.com/gorilla/websocket v1.5.0
	github.com/jbloxsome/raptor/orderbook v0.0.0-20220527171206-0884d3234ae5
)

require github.com/jbloxsome/raptor/rbbst v0.0.0-00010101000000-000000000000 // indirect

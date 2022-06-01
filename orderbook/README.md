# Raptor Orderbook
Raptor's orderbook package provides a convenient wrapper around the rbbst package for maintaining a level 2 orderbook with high read/write performance.

## Usage
### Creating an orderbook
```go
ob := NewOrderbook()
```
### Inserting, Removing and Fetching Levels
To insert a bid level with price @ 34542.12 and volume @ 1.234
```go
ob.AddBidLevel(34542.12, 1.234)
```

To remove a bid level with price @ 34542.12
```go
ob.RemoveBidLevel(34542.12)
```

To fetch a bid level with price @ 34542.12 (returns a []float64, with level[0] being the price and level[1] being the volume).
```go
level := ob.GetBidLevel(34542.12)
```

## Running the tests and benchmarks
```
go test -v -bench=.
```

## Performance
Benchmarking with an orderbook with 10k bid levels and 10k ask levels resulted in the following:
- AddBidLevel: ~70 ns/op
- GetBidLevel: ~90 ns/op
- AddAskLevel: ~70 ns/op
- GetAskLevel: ~90 ns/op
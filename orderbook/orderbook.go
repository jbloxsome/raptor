package orderbook

import (
	"github.com/jbloxsome/raptor/rbbst"
)

type Orderbook struct {
	Bids *rbbst.RedBlackBST
	Asks *rbbst.RedBlackBST
}

func NewOrderbook() Orderbook {
	bids := rbbst.NewRedBlackBST()
	asks := rbbst.NewRedBlackBST()

	return Orderbook{
		Asks: &asks,
		Bids: &bids,
	}
}

func (ob *Orderbook) AddBidLevel(price float64, volume float64) {
	ob.Bids.Put(price, volume)
}

func (ob *Orderbook) RemoveBidLevel(price float64) {
	ob.Bids.Delete(price)
}

func (ob *Orderbook) GetBidLevel(price float64) []float64 {
	contains := ob.Bids.Contains(price)
	if !contains {
		return nil
	}
	
	volume := ob.Bids.Get(price)
	return []float64{price, volume}
}

func (ob *Orderbook) AddAskLevel(price float64, volume float64) {
	ob.Asks.Put(price, volume)
}

func (ob *Orderbook) RemoveAskLevel(price float64) {
	ob.Asks.Delete(price)
}

func (ob *Orderbook) GetAskLevel(price float64) []float64 {
	contains := ob.Asks.Contains(price)
	if !contains {
		return nil
	}

	volume :=  ob.Asks.Get(price)
	return []float64{price, volume}
}


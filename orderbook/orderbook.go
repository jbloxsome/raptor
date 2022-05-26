package orderbook

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type Orderbook struct {
	Bids *rbt.Tree
	Asks *rbt.Tree
}

func NewOrderbook() Orderbook {
	bids := rbt.NewWithStringComparator()
	asks := rbt.NewWithStringComparator()

	return Orderbook{
		Asks: &asks,
		Bids: &bids
	}
}

func (ob *Orderbook) AddBidLevel(price string, volume string) {
	ob.Bids.Put(price, volume)
}

func (ob *Orderbook) RemoveBidLevel(price string) {
	ob.Bids.Remove(price)
}

func (ob *Orderbook) GetBidLevel(price string) []string {
	volume, found :=  ob.Bids.Get(price)

	if found {
		return []string{price, volume}
	} else {
		return nil
	}
}

func (ob *Orderbook) AddAskLevel(price string, volume string) {
	ob.Ask.Put(price, volume)
}

func (ob *Orderbook) RemoveAskLevel(price string) {
	ob.Ask.Remove(price)
}

func (ob *Orderbook) GetAskLevel(price string) []string {
	volume, found :=  ob.Ask.Get(price)

	if found {
		return []string{price, volume}
	} else {
		return nil
	}
}


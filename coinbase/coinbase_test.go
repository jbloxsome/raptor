package coinbase

import (
	"testing"

	ob "github.com/jbloxsome/raptor/orderbook"
)

func TestHandleSnapshot(t *testing.T) {
	message := &Message{
		Type: "snapshot",
		ProductId: "BTC-USD",
		Bids: [][]string{[]string{"30000.00", "0.453"}, []string{"30050.00", "0.324"}},
		Asks: [][]string{[]string{"30100.00", "0.532"}, []string{"30150.00", "0.345"}},
	}


	ob := ob.NewOrderbook()

	handleSnapshot(&ob, message)

	level := ob.GetBidLevel(30000.00)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{30000.00, 0.453})
	}

	if level[0] != 30000.00 {
		t.Errorf("got %f, want %v", level, []float64{30000.00, 0.453})
	}

	if level[1] != 0.453 {
		t.Errorf("got %f, want %v", level, []float64{30000, 0.453})
	}

	level = ob.GetBidLevel(30050.00)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{30050.00, 0.324})
	}

	if level[0] != 30050.00 {
		t.Errorf("got %f, want %v", level, []float64{30050.00, 0.324})
	}

	if level[1] != 0.324 {
		t.Errorf("got %f, want %v", level, []float64{30050.00, 0.324})
	}
}

func TestHandleL2Update(t *testing.T) {
	// Create an orderbook and add 5000 bid levels and 5000 ask levels
	ob := ob.NewOrderbook()

	for i := 30000.00; i <= 35000.00; i++ {
		ob.AddBidLevel(i, 0.342)
	}

	for j := 35001.00; j <= 40000.00; j++ {
		ob.AddAskLevel(j, 0.343)
	}

	message := &Message{
		Type: "l2update",
		ProductId: "BTC-USD",
		Changes: [][]string{[]string{"buy", "35000.00", "0.100"}},
	}

	handleL2Update(&ob, message)

	// {"buy", "35000.00", "0.100"} should replace the 35000 bid level volume
	// with a new volume of 0.100
	level := ob.GetBidLevel(35000.00)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{30000.00, 0.100})
	}

	if level[0] != 35000.00 {
		t.Errorf("got %f, want %v", level, []float64{35000.00, 0.100})
	}

	if level[1] != 0.100 {
		t.Errorf("got %f, want %v", level, []float64{35000.00, 0.100})
	}
}

func BenchmarkL2Update(b *testing.B) {
	ob := ob.NewOrderbook()

	for i := 30000.00; i <= 35000.00; i++ {
		ob.AddBidLevel(i, 0.342)
	}

	for j := 35001.00; j <= 40000.00; j++ {
		ob.AddAskLevel(j, 0.343)
	}

	message := &Message{
		Type: "l2update",
		ProductId: "BTC-USD",
		Changes: [][]string{
			[]string{"buy", "35000.00", "0.100"},
			[]string{"buy", "35001.00", "0.321"},
			[]string{"sell", "35001.50", "3.424"},
		},
	}	

	for i := 0; i < b.N; i++ {
		handleL2Update(&ob, message)
	}
}

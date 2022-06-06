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


	orderbook := ob.NewOrderbook()

	handleSnapshot(&orderbook, message)

	level := orderbook.GetBidLevel(30000.00)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{30000.00, 0.453})
	}

	if level[0] != 30000.00 {
		t.Errorf("got %f, want %v", level, []float64{30000.00, 0.453})
	}

	if level[1] != 0.453 {
		t.Errorf("got %f, want %v", level, []float64{30000, 0.453})
	}

	level = orderbook.GetBidLevel(30050.00)

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

// TODO:
// Remove Ask Level
func TestHandleL2Update(t *testing.T) {
	// Create an orderbook and add 5000 bid levels and 5000 ask levels
	orderbook := ob.NewOrderbook()

	for i := 30000.00; i <= 35000.00; i++ {
		orderbook.AddBidLevel(i, 0.342)
	}

	for j := 35001.00; j <= 40000.00; j++ {
		orderbook.AddAskLevel(j, 0.343)
	}

	message := &Message{
		Type: "l2update",
		ProductId: "BTC-USD",
		Changes: [][]string{
			[]string{"buy", "29000.00", "0.400"},
			[]string{"buy", "35000.00", "0.100"},
			[]string{"buy", "34999.00", "0.000"},
			[]string{"sell", "35000.50", "0.100"},
			[]string{"sell", "35002.00", "0.250"},
			[]string{"sell", "35003.00", "0.000"},
		},
	}

	handleL2Update(&orderbook, message)

	// {"buy", "29000.00", "0.400"} should add a new bid level at 29000.00 with volume
	// of 0.400
	level := orderbook.GetBidLevel(29000.00)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{29000.00, 0.400})
	}

	if level[0] != 29000.00 {
		t.Errorf("got %f, want %v", level, []float64{29000.00, 0.400})
	}

	if level[1] != 0.400 {
		t.Errorf("got %f, want %v", level, []float64{29000.00, 0.400})
	}

	// {"buy", "35000.00", "0.100"} should replace the 35000 bid level volume
	// with a new volume of 0.100
	level = orderbook.GetBidLevel(35000.00)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{30000.00, 0.100})
	}

	if level[0] != 35000.00 {
		t.Errorf("got %f, want %v", level, []float64{35000.00, 0.100})
	}

	if level[1] != 0.100 {
		t.Errorf("got %f, want %v", level, []float64{35000.00, 0.100})
	}

	// {"buy", "34999.00", "0.000"} should remove the 34999 bid level
	level = orderbook.GetBidLevel(34999.00)

	if level != nil {
		t.Errorf("got %f, want nil", level)
	}

	// {"sell", "35000.50", "0.100"} should add a new ask level at 35000.50 with volume
	// of 0.100
	level = orderbook.GetAskLevel(35000.50)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{35000.50, 0.100})
	}

	if level[0] != 35000.50 {
		t.Errorf("got %f, want %v", level, []float64{35000.50, 0.100})
	}

	if level[1] != 0.100 {
		t.Errorf("got %f, want %v", level, []float64{35000.50, 0.100})
	}

	// {"sell", "35002.00", "0.250"} should replace the 35002 ask level volume
	// with a new volume of 0.250
	level = orderbook.GetAskLevel(35002.00)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{35002.00, 0.250})
	}

	if level[0] != 35002.00 {
		t.Errorf("got %f, want %v", level, []float64{35002.00, 0.250})
	}

	if level[1] != 0.250 {
		t.Errorf("got %f, want %v", level, []float64{35002.00, 0.250})
	}

	// {"sell", "35003.00", "0.000"} should remove the 35003 ask level
	level = orderbook.GetAskLevel(35003.00)

	if level != nil {
		t.Errorf("got %f, want nil", level)
	}
}

func BenchmarkL2Update(b *testing.B) {
	orderbook := ob.NewOrderbook()

	for i := 30000.00; i <= 35000.00; i++ {
		orderbook.AddBidLevel(i, 0.342)
	}

	for j := 35001.00; j <= 40000.00; j++ {
		orderbook.AddAskLevel(j, 0.343)
	}

	message := &Message{
		Type: "l2update",
		ProductId: "BTC-USD",
		Changes: [][]string{
			[]string{"buy", "29000.00", "0.400"},
			[]string{"buy", "35000.00", "0.100"},
			[]string{"buy", "34999.00", "0.421"},
			[]string{"sell", "35000.50", "0.100"},
			[]string{"sell", "35002.00", "0.250"},
			[]string{"sell", "35003.00", "0.400"},
		},
	}

	for i := 0; i < b.N; i++ {
		handleL2Update(&orderbook, message)
	}
}

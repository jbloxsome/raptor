package orderbook

import (
	"testing"
)

func TestAddBidLevel(t *testing.T) {
	ob := NewOrderbook()

	ob.AddBidLevel(34000.00, 1.23)

	level := ob.GetBidLevel(34000.00)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{34000.00, 1.23})
	}

	if level[0] != 34000.00 {
		t.Errorf("got %f, want %v", level, []float64{34000.00, 1.23})
	}

	if level[1] != 1.23 {
		t.Errorf("got %f, want %v", level, []float64{34000.00, 1.23})
	}
}

func TestRemoveBidLevel(t *testing.T) {
	ob := NewOrderbook()

	ob.AddBidLevel(34000.00, 1.23)

	ob.RemoveBidLevel(34000.00)

	level := ob.GetBidLevel(34000.00)

	if level != nil {
		t.Errorf("got %f, want nil", []float64{34000.00, 1.23})
	}
}

func TestGetBidLevel(t *testing.T) {
	ob := NewOrderbook()

	ob.AddBidLevel(34000.00, 1.23)

	level := ob.GetBidLevel(34000.00)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{34000.00, 1.23})
	}

	if level[0] != 34000.00 {
		t.Errorf("got %f, want %v", level, []float64{34000.00, 1.23})
	}

	if level[1] != 1.23 {
		t.Errorf("got %f, want %v", level, []float64{34000.00, 1.23})
	}
}

func TestAddAskLevel(t *testing.T) {
	ob := NewOrderbook()

	ob.AddAskLevel(34000.00, 1.23)

	level := ob.GetAskLevel(34000.00)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{34000.00, 1.23})
	}

	if level[0] != 34000.00 {
		t.Errorf("got %f, want %v", level, []float64{34000.00, 1.23})
	}

	if level[1] != 1.23 {
		t.Errorf("got %f, want %v", level, []float64{34000.00, 1.23})
	}
}

func TestRemoveAskLevel(t *testing.T) {
	ob := NewOrderbook()

	ob.AddAskLevel(34000.00, 1.23)

	ob.RemoveAskLevel(34000.00)

	level := ob.GetAskLevel(34000.00)

	if level != nil {
		t.Errorf("got %f, want nil", []float64{34000.00, 1.23})
	}
}

func TestGetAskLevel(t *testing.T) {
	ob := NewOrderbook()

	ob.AddAskLevel(34000.00, 1.23)

	level := ob.GetAskLevel(34000.00)

	if level == nil {
		t.Errorf("got nil, want %f", []float64{34000.00, 1.23})
	}

	if level[0] != 34000.00 {
		t.Errorf("got %f, want %v", level, []float64{34000.00, 1.23})
	}

	if level[1] != 1.23 {
		t.Errorf("got %f, want %v", level, []float64{34000.00, 1.23})
	}
}

func BenchmarkAddBidLevel(b *testing.B) {
	ob := NewOrderbook()

	for i := 0; i < 10000; i++ {
		ob.AddBidLevel(float64(i), float64(i))
		ob.AddAskLevel(float64(i), float64(i))
	}

	for i := 0; i < b.N; i++ {
		ob.AddBidLevel(5000.50, 1.345)
	}
}

func BenchmarkGetBidLevel(b *testing.B) {
	ob := NewOrderbook()

	for i := 0; i < 10000; i++ {
		ob.AddBidLevel(float64(i), float64(i))
		ob.AddAskLevel(float64(i), float64(i))
	}

	for i := 0; i < b.N; i++ {
		ob.GetBidLevel(5000.00)
	}
}

func BenchmarkAddAskLevel(b *testing.B) {
	ob := NewOrderbook()

	for i := 0; i < 10000; i++ {
		ob.AddBidLevel(float64(i), float64(i))
		ob.AddAskLevel(float64(i), float64(i))
	}

	for i := 0; i < b.N; i++ {
		ob.AddAskLevel(5000.50, 1.345)
	}
}

func BenchmarkGetAskLevel(b *testing.B) {
	ob := NewOrderbook()

	for i := 0; i < 10000; i++ {
		ob.AddBidLevel(float64(i), float64(i))
		ob.AddAskLevel(float64(i), float64(i))
	}

	for i := 0; i < b.N; i++ {
		ob.GetAskLevel(5000.00)
	}
}
package rbbst

import (
	"testing"
)

func TestSize(t *testing.T) {
	rbt := NewRedBlackBST()

	for i := 0; i < 10; i++ {
		rbt.Put(float64(i), float64(i + 1))
	}
	
	size := rbt.Size()

	if size != 10 {
		t.Errorf("got %v, want %v", size, 10)
	}
}

func TestIsEmpty(t *testing.T) {
	rbt := NewRedBlackBST()

	isEmpty := rbt.IsEmpty()

	if isEmpty != true {
		t.Errorf("got %t, want %t", isEmpty, true)
	}

	rbt.Put(32453.04, 1.0234)

	isEmpty = rbt.IsEmpty()

	if isEmpty != false {
		t.Errorf("got %t, want %t", isEmpty, false)
	}
}

func TestContains(t *testing.T) {
	rbt := NewRedBlackBST()

	contains := rbt.Contains(35250.12)

	if contains != false {
		t.Errorf("got %t, want %t", contains, false)
	}

	rbt.Put(35250.12, 1.342)
	contains = rbt.Contains(35250.12)

	if contains != true {
		t.Errorf("got %t, want %t", contains, true)
	}
}

func TestGet(t *testing.T) {
	rbt := NewRedBlackBST()

	rbt.Put(35250.12, 1.342)
	rbt.Put(35124.12, 0.445)
	rbt.Put(35563.12, 5.456)

	val := rbt.Get(35124.12)

	if val != 0.445 {
		t.Errorf("got %f, want %f", val, 0.445)
	}
}

func TestPut(t *testing.T) {
	rbt := NewRedBlackBST()

	rbt.Put(35250.12, 1.342)

	size := rbt.Size()

	if size != 1 {
		t.Errorf("got %v, want %v", size, 1)
	}

	val := rbt.Get(35250.12)

	if val != 1.342 {
		t.Errorf("got %f, want %f", val, 1.342)
	}
}

func TestHeight(t *testing.T) {
	rbt := NewRedBlackBST()

	for i := 0; i < 10; i++ {
		rbt.Put(float64(i), float64(i + 1))
	}

	height := rbt.Height()

	if height != 4 {
		t.Errorf("got %v, want %v", height, 4)
	}
}

func TestMin(t *testing.T) {
	rbt := NewRedBlackBST()

	for i := 0; i < 100; i++ {
		rbt.Put(float64(i), float64(i))
	}

	min := rbt.Min()

	if min != 0 {
		t.Errorf("got %v, want %v", min, 0)
	}
}

func TestMax(t *testing.T) {
	rbt := NewRedBlackBST()

	for i := 0; i < 100; i++ {
		rbt.Put(float64(i), float64(i))
	}

	max := rbt.Max()

	if max != 99 {
		t.Errorf("got %v, want %v", max, 99)
	}
}

func TestDelete(t *testing.T) {
	rbt := NewRedBlackBST()

	for i := 1; i < 101; i++ {
		rbt.Put(float64(i), float64(i))
	}

	rbt.Delete(50)

	size := rbt.Size()

	if size != 99 {
		t.Errorf("got %v, want %v", size, 99)
	}

	deleted := rbt.Contains(50)

	if deleted != false {
		t.Errorf("got %t, want %t", deleted, false)
	}
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []float64) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

func TestKeys(t *testing.T) {
	rbt := NewRedBlackBST()

	for i := 1; i < 5; i++ {
		rbt.Put(float64(i), float64(i))
	}

	keys := rbt.Keys(1, 4)

	expected := []float64{1, 2, 3, 4}

	if Equal(keys, expected) != true {
		t.Errorf("got %f, want %f", keys, expected)
	}
}

func BenchmarkPut(b *testing.B) {
	rbt := NewRedBlackBST()

	for i := 0; i < 10000; i++ {
		rbt.Put(float64(i), float64(i))
	}

	for i := 0; i < b.N; i++ {
		rbt.Put(50000.00, 50000.00)
	}
}

func BenchmarkGet(b *testing.B) {
	rbt := NewRedBlackBST()

	for i := 0; i < 10000; i++ {
		rbt.Put(float64(i), float64(i))
	}

	for i := 0; i < b.N; i++ {
		rbt.Get(5000.00)
	}
}
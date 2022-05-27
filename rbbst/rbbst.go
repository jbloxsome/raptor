package rbbst

import (
	"fmt"
)

type NodeRedBlack struct {
	Key float64
	Value float64
	Next *NodeRedBlack
	Prev *NodeRedBlack
	
	left *NodeRedBlack
	right *NodeRedBlack
	size int
	isRed bool
}

type RedBlackBST struct {
	root *NodeRedBlack
	minC *NodeRedBlack // cached min/max keys for O(1) access
	maxC *NodeRedBlack
}

func NewRedBlackBST() RedBlackBST {
	return RedBlackBST{}
}

func (t *RedBlackBST) Size() int {
	return t.size(t.root)
}

func (t *RedBlackBST) size(n *NodeRedBlack) int {
	if n == nil {
		return 0
	}

	return n.size
}

func (t *RedBlackBST) IsEmpty() bool {
	return t.size(t.root) == 0
}

func (t *RedBlackBST) panicIfEmpty() {
	if t.IsEmpty() {
		panic("Red Black BST is empty")
	}
}

func (t *RedBlackBST) Contains(key float64) bool {
	return t.get(t.root, key) != nil
}

func (t *RedBlackBST) Get(key float64) float64 {
	t.panicIfEmpty()

	x := t.get(t.root, key)
	if x == nil {
		panic(fmt.Sprintf("key %0.8f does not exist", key))
	}

	return x.Value
}

func (t *RedBlackBST) get(n *NodeRedBlack, key float64) *NodeRedBlack {
	if n == nil {
		return nil
	}

	if n.Key == key {
		return n
	}

	if n.Key > key {
		return t.get(n.left, key)
	} else {
		return t.get(n.right, key)
	}
}

func (t *RedBlackBST) isRed(n *NodeRedBlack) bool {
	if n == nil {
		// nil nodes are black by default
		return false
	}

	return n.isRed
}

func (t *RedBlackBST) flipColors(n *NodeRedBlack) {
	if n == nil {
		return
	}

	// inverse children colors
	if n.left != nil {
		n.left.isRed = !n.left.isRed
	}
	if n.right != nil {
		n.right.isRed = !n.right.isRed
	}

	// inverse node color
	n.isRed = !n.isRed
}

func (t *RedBlackBST) rotateLeft(n *NodeRedBlack) *NodeRedBlack {
	x := n.right
	n.right = x.left
	x.left = n

	x.isRed = n.isRed
	n.isRed = true

	// re-calculate sizes
	n.size = t.size(n.left) + 1 + t.size(n.right)
	x.size = t.size(x.left) + 1 + t.size(x.right)

	return x
}

func (t *RedBlackBST) rotateRight(n *NodeRedBlack) *NodeRedBlack {
	x := n.left
	n.left = x.right
	x.right = n

	x.isRed = n.isRed
	n.isRed = true

	// re-calculate sizes
	n.size = t.size(n.left) + 1 + t.size(n.right)
	x.size = t.size(x.left) + 1 + t.size(x.right)

	return x
}

func (t *RedBlackBST) Put(key float64, value float64) {
	t.root = t.put(t.root, key, value)

	// keeping root black
	t.root.isRed = false
}

func (t *RedBlackBST) put(n *NodeRedBlack, key float64, value float64) *NodeRedBlack {
	if n == nil {
		// search miss, creating a new node with a red link as a part of 3- or 4-node
		n := &NodeRedBlack{
			Value: value,
			Key: key,
			size: 1,
			isRed: true,
		}

		if t.minC == nil || key < t.minC.Key {
			// new min
			t.minC = n
		}
		if t.maxC == nil || key > t.maxC.Key {
			// new max
			t.maxC = n
		}

		return n
	}

	if n.Key == key {
		// search hit, updating the value
		n.Value = value
		return n
	}

	if n.Key > key {
		left := n.left
		n.left = t.put(n.left, key, value)
		if left == nil {
			// new node has been just inserted to the left
			prev := n.Prev
			if prev != nil {
				prev.Next = n.left	
			}
			n.left.Prev = prev
			n.left.Next = n
			n.Prev = n.left
		}
	} else {
		right := n.right
		n.right = t.put(n.right, key, value)
		if right == nil {
			// new node has been just inserted to the right
			next := n.Next
			if next != nil {
				next.Prev = n.right
			}
			n.right.Next = next
			n.right.Prev = n
			n.Next = n.right
		}
	}

	// balancing the tree
	if t.isRed(n.right) && !t.isRed(n.left) {
		// fixing right leaning red link case
		// this can lead to the next case in upper level
		n = t.rotateLeft(n)
	}
	if t.isRed(n.left) && t.isRed(n.left.left) {
		// making 4-node
		n = t.rotateRight(n)
	}
	if t.isRed(n.left) && t.isRed(n.right) {
		// convert 4-node into 3 2-nodes
		t.flipColors(n)
	}

	// re-calc size
	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *RedBlackBST) Height() int {
	if t.IsEmpty() {
		return 0
	}

	return t.height(t.root)
}

func (t *RedBlackBST) height(n *NodeRedBlack) int {
	if n == nil {
		return 0
	}

	lheight := t.height(n.left)
	rheight := t.height(n.right)
	
	height := lheight
	if rheight > lheight {
		height = rheight
	}

	return height + 1
}

func (t *RedBlackBST) IsRedBlack() bool {
	balanced, _ := t.isBalanced(t.root)
	return balanced && t.is23(t.root)
}

func (t *RedBlackBST) isBalanced(n *NodeRedBlack) (bool, int) {
	if n == nil {
		// nil node is black by default
		return true, 1
	}

	lb, l := t.isBalanced(n.left)
	rb, r := t.isBalanced(n.right)

	b := l
	if r > l {
		b = r
	}

	if !t.isRed(n) {
		b += 1
	}

	return lb && rb && l == r, b
}

func (t *RedBlackBST) is23(n *NodeRedBlack) bool {
	if n == nil {
		return true
	}

	if t.isRed(n.right) {
		// it should has only left leaning red links
		return false
	}

	if t.isRed(n) && t.isRed(n.left) {
		// no node should be connected by two red links
		return false
	}

	return t.is23(n.left) && t.is23(n.right)
}

func (t *RedBlackBST) Min() float64 {
	t.panicIfEmpty()
	return t.minC.Key
}

func (t *RedBlackBST) MinValue() float64 {
	t.panicIfEmpty()
	return t.minC.Value
}

func (t *RedBlackBST) MinPointer() *NodeRedBlack {
	t.panicIfEmpty()
	return t.minC
}

func (t *RedBlackBST) min(n *NodeRedBlack) *NodeRedBlack {
	if n.left == nil {
		return n
	}

	return t.min(n.left)
}

func (t *RedBlackBST) Max() float64 {
	t.panicIfEmpty()
	return t.maxC.Key
}

func (t *RedBlackBST) MaxValue() float64 {
	t.panicIfEmpty()
	return t.maxC.Value
}

func (t *RedBlackBST) MaxPointer() *NodeRedBlack {
	t.panicIfEmpty()
	return t.maxC
}

func (t *RedBlackBST) max(n *NodeRedBlack) *NodeRedBlack {
	if n.right == nil {
		return n
	}

	return t.max(n.right)
}

func (t *RedBlackBST) Floor(key float64) float64 {
	t.panicIfEmpty()

	floor := t.floor(t.root, key)
	if floor == nil {
		panic(fmt.Sprintf("there are no keys <= %0.8f", key))
	}

	return floor.Key
}

func (t *RedBlackBST) floor(n *NodeRedBlack, key float64) *NodeRedBlack {
	if n == nil {
		// search miss
		return nil
	}

	if n.Key == key {
		// search hit
		return n
	}

	if n.Key > key {
		// floor must be in the left sub-tree
		return t.floor(n.left, key)
	}

	// key could be in the right sub-tree, if not, using current root
	floor := t.floor(n.right, key)
	if floor != nil {
		return floor
	}

	return n
}

func (t *RedBlackBST) Ceiling(key float64) float64 {
	t.panicIfEmpty()

	ceiling := t.ceiling(t.root, key)
	if ceiling == nil {
		panic(fmt.Sprintf("there are no keys >= %0.8f", key))
	}

	return ceiling.Key
}

func (t *RedBlackBST) ceiling(n *NodeRedBlack, key float64) *NodeRedBlack {
	if n == nil {
		// search miss
		return nil
	}

	if n.Key == key {
		// search hit
		return n
	}

	if n.Key < key {
		// ceiling must be in the right sub-tree
		return t.ceiling(n.right, key)
	}

	// the key could be in the left sub-tree, if not, using current root
	ceiling := t.ceiling(n.left, key)
	if ceiling != nil {
		return ceiling
	}

	return n
}

func (t *RedBlackBST) Select(k int) float64 {
	if k < 0 || k >= t.Size() {
		panic("index out of range")
	}

	return t.selectNode(t.root, k).Key
}

func (t *RedBlackBST) selectNode(n *NodeRedBlack, k int) *NodeRedBlack {
	if t.size(n.left) == k {
		return n
	}

	if t.size(n.left) > k {
		return t.selectNode(n.left, k)
	}

	k = k - t.size(n.left) - 1
	return t.selectNode(n.right, k)
}

func (t *RedBlackBST) Rank(key float64) int {
	t.panicIfEmpty()
	return t.rank(t.root, key)
}

func (t *RedBlackBST) rank(n *NodeRedBlack, key float64) int {
	if n == nil {
		return 0
	}

	if n.Key == key {
		return t.size(n.left)
	}

	if n.Key > key {
		return t.rank(n.left, key)
	}

	return t.size(n.left) + 1 + t.rank(n.right, key)
}

func (t *RedBlackBST) moveRedLeft(n *NodeRedBlack) *NodeRedBlack {
	// assuming that n.left and n.left.left are black and n is red,
	// make h.left or one of its children red
	t.flipColors(n)
	// now n is black and both left and right are red

	// fixing red black invariat that no node can be connected with two red links
	if t.isRed(n.right.left) {
		n.right = t.rotateRight(n.right)
		// now n.right and n.right.right are red, fixing that by rotating n
		n = t.rotateLeft(n) 
		// now n.right, n.right.right and n.left are red

		t.flipColors(n)
		// now n.left and n.right are black, n.left.left is red
	}

	return n
}

func (t *RedBlackBST) DeleteMin() {
	t.panicIfEmpty()

	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		// making root red temporarily to fit invariant required for moveRedLeft method
		t.root.isRed = true
	}
	t.root = t.deleteMin(t.root)
	if !t.IsEmpty() {
		t.root.isRed = false
	}
}

func (t *RedBlackBST) deleteMin(n *NodeRedBlack) *NodeRedBlack {
	if n.left == nil {
		// we've reached the least leave of the tree
		next := n.Next
		prev := n.Prev
		if prev != nil {
			prev.Next = next
		}
		if next != nil {
			next.Prev = prev
		}
		n.Next = nil
		n.Prev = nil
		
		// updating global min
		if t.minC == n {
			t.minC = next
		}

		return n.right
	}

	// making current node a part of 3 or 4 node by moving red link to the left
	if !t.isRed(n.left) && !t.isRed(n.left.left) {
		n = t.moveRedLeft(n)
	}

	n.left = t.deleteMin(n.left)

	// we have to restore balance of the tree moving from bottom to top now
	if t.isRed(n.right) {
		n = t.rotateLeft(n)
	}
	if t.isRed(n.left) && t.isRed(n.left.left) {
		n = t.rotateRight(n)
	}
	if t.isRed(n.left) && t.isRed(n.right) {
		t.flipColors(n)
	}

	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *RedBlackBST) moveRedRight(n *NodeRedBlack) *NodeRedBlack {
	// assuming n is red, n.right and n.right.left are black,
	// make h.right or one of its children red
	t.flipColors(n)
	// now n is black and n.right is red

	if t.isRed(n.left.left) {
		// meaning n.left should be red now after fliping the colors of n
		n = t.rotateRight(n)
		// now n.left is red, n.right and n.right.right are red

		t.flipColors(n)
		// now n.left is black, n.right is black, n.right.right is red
	}
	return n
}

func (t *RedBlackBST) DeleteMax() {
	t.panicIfEmpty()

	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		t.root.isRed = true
	}
	t.root = t.deleteMax(t.root)
	if !t.IsEmpty() {
		t.root.isRed = false;
	}
}

func (t *RedBlackBST) deleteMax(n *NodeRedBlack) *NodeRedBlack {
	if t.isRed(n.left) {
		// making right red by rotating
		n = t.rotateRight(n)
	}
	if n.right == nil {
		// we've reached the largest key in the tree
		next := n.Next
		prev := n.Prev
		if prev != nil {
			prev.Next = next
		}
		if next != nil {
			next.Prev = prev
		}
		n.Next = nil
		n.Prev = nil

		// updating global max
		if t.maxC == n {
			t.maxC = prev
		}

		return n.left
	}

	// making right left on the way from top to bottom
	if !t.isRed(n.right) && !t.isRed(n.right.left) {
		n = t.moveRedRight(n)
	}

	n.right = t.deleteMax(n.right)
	
	// balancing back on the way from bottom to top
	if t.isRed(n.right) {
		n = t.rotateLeft(n)
	}
	if t.isRed(n.left) && t.isRed(n.left.left) {
		n = t.rotateRight(n)
	}
	if t.isRed(n.left) && t.isRed(n.right) {
		t.flipColors(n)
	}

	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *RedBlackBST) Delete(key float64) {
	t.panicIfEmpty()

	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		t.root.isRed = true
	}
	t.root = t.delete(t.root, key)
	if !t.IsEmpty() {
		t.root.isRed = false;
	}
}

func (t *RedBlackBST) delete(n *NodeRedBlack, key float64) *NodeRedBlack {
	if n.Key > key {
		if n.left == nil {
			// search miss
			return nil
		}

		// looking into the left sub-tree
		if !t.isRed(n.left) && !t.isRed(n.left.left) {
			n = t.moveRedLeft(n)
		}

		n.left = t.delete(n.left, key)

	} else {
		// checking current node and right sub-tree if required
		if t.isRed(n.left) {
			n = t.rotateRight(n)
		}
		if n.Key == key && n.right == nil {
			// search hit and we don't have right sub-tree
			
			// updating linked list
			next := n.Next
			prev := n.Prev
			if prev != nil {
				prev.Next = next
			}
			if next != nil {
				next.Prev = prev
			}
			n.Next = nil
			n.Prev = nil

			if t.maxC == n {
				t.maxC = prev
			}
			if t.minC == n {
				t.minC = next
			}

			return nil
		}

		if !t.isRed(n.right) && !t.isRed(n.right.left) {
			n = t.moveRedRight(n)
		}
		// h.right or one of its children red make

		if n.Key == key {
			// search hit, replacing the node with a successor
			rightMin := t.min(n.right)
			n.Key = rightMin.Key
			n.Value = rightMin.Value
			n.right = t.deleteMin(n.right)

			// global min will be updated automatically if requied, 
			// as we copy values from successor
		} else {
			if n.right == nil {
				// search miss
				return nil
			}

			n.right = t.delete(n.right, key)
		}
	}

	// balance
	if t.isRed(n.right) {
		n = t.rotateLeft(n)
	}
	if t.isRed(n.left) && t.isRed(n.left.left) {
		n = t.rotateRight(n)
	}
	if t.isRed(n.left) && t.isRed(n.right) {
		t.flipColors(n)
	}

	n.size = t.size(n.left) + 1 + t.size(n.right)
	return n
}

func (t *RedBlackBST) Keys(lo, hi float64) []float64 {
	if lo < t.Min() || hi > t.Max() {
		panic("keys out of range")
	}

	return t.keys(t.root, lo, hi)
}

func (t *RedBlackBST) keys(n *NodeRedBlack, lo, hi float64) []float64 {
	if n == nil {
		return nil
	}

	if n.Key < lo {
		return t.keys(n.right, lo, hi)
	} else if n.Key > hi {
		return t.keys(n.left, lo, hi)
	}

	l := t.keys(n.left, lo, hi)
	r := t.keys(n.right, lo, hi)
	
	keys := make([]float64, 0)
	if l != nil {
		keys = append(keys, l...)
	}
	keys = append(keys, n.Key)
	if r != nil {
		keys = append(keys, r...)
	}

	return keys
}

func (t *RedBlackBST) Print() {
	fmt.Println()
	t.print(t.root)
	fmt.Println()
}

func (t *RedBlackBST) print(n *NodeRedBlack) {
	if n == nil {
		return
	}

	if n.isRed {
		fmt.Printf("*")
	}
	fmt.Printf("%0.8f ", n.Key)

	t.print(n.left)
	t.print(n.right)
}
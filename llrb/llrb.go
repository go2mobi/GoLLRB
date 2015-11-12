// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A Left-Leaning Red-Black (LLRB) implementation of 2-3 balanced binary search trees,
// based on the following work:
//
//   http://www.cs.princeton.edu/~rs/talks/LLRB/08Penn.pdf
//   http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf
//   http://www.cs.princeton.edu/~rs/talks/LLRB/Java/RedBlackBST.java
//
//  2-3 trees (and the run-time equivalent 2-3-4 trees) are the de facto standard BST
//  algoritms found in implementations of Python, Java, and other libraries. The LLRB
//  implementation of 2-3 trees is a recent improvement on the traditional implementation,
//  observed and documented by Robert Sedgewick.
//
package llrb

// Tree is a Left-Leaning Red-Black (LLRB) implementation of 2-3 trees
type LLRB struct {
	count int
	root  *Node
}

type Node struct {
	Value       uint64
	Left, Right *Node // Pointers to left and right child nodes
	Black       bool  // If set, the color of the link (incoming from the parent) is black
	// In the LLRB, new nodes are always red, hence the zero-value for node
}

func less(x, y uint64) bool {
	return x < y
}

// New() allocates a new tree
func New() *LLRB {
	return &LLRB{}
}

// SetRoot sets the root node of the tree.
// It is intended to be used by functions that deserialize the tree.
func (t *LLRB) SetRoot(r *Node) {
	t.root = r
}

func (t *LLRB) SetLen(l int) {
	t.count = l
}

// Root returns the root node of the tree.
// It is intended to be used by functions that serialize the tree.
func (t *LLRB) Root() *Node {
	return t.root
}

// Len returns the number of nodes in the tree.
func (t *LLRB) Len() int { return t.count }

// Has returns true if the tree contains an element whose order is the same as that of key.
func (t *LLRB) Has(key uint64) bool {
	_, found := t.Get(key)
	return found
}

// Get retrieves an element from the tree whose order is the same as that of key.
func (t *LLRB) Get(key uint64) (uint64, bool) {
	h := t.root
	for h != nil {
		switch {
		case less(key, h.Value):
			h = h.Left
		case less(h.Value, key):
			h = h.Right
		default:
			return h.Value, true
		}
	}
	return 0, false
}

// Min returns the minimum element in the tree.
func (t *LLRB) Min() uint64 {
	h := t.root
	if h == nil {
		return 0
	}
	for h.Left != nil {
		h = h.Left
	}
	return h.Value
}

// Max returns the maximum element in the tree.
func (t *LLRB) Max() uint64 {
	h := t.root
	if h == nil {
		return 0
	}
	for h.Right != nil {
		h = h.Right
	}
	return h.Value
}

func (t *LLRB) ReplaceOrInsertBulk(items ...uint64) {
	for _, i := range items {
		t.ReplaceOrInsert(i)
	}
}

func (t *LLRB) InsertNoReplaceBulk(items ...uint64) {
	for _, i := range items {
		t.InsertNoReplace(i)
	}
}

// ReplaceOrInsert inserts item into the tree. If an existing
// element has the same order, it is removed from the tree and returned.
func (t *LLRB) ReplaceOrInsert(item uint64) bool {
	var replaced bool
	t.root, replaced = t.replaceOrInsert(t.root, item)
	t.root.Black = true
	if !replaced {
		t.count++
	}
	return replaced
}

func (t *LLRB) replaceOrInsert(h *Node, item uint64) (*Node, bool) {
	if h == nil {
		return newNode(item), false
	}

	h = walkDownRot23(h)

	var replaced bool
	if less(item, h.Value) { // BUG
		h.Left, replaced = t.replaceOrInsert(h.Left, item)
	} else if less(h.Value, item) {
		h.Right, replaced = t.replaceOrInsert(h.Right, item)
	} else {
		replaced, h.Value = true, item
	}

	h = walkUpRot23(h)

	return h, replaced
}

// InsertNoReplace inserts item into the tree. If an existing
// element has the same order, both elements remain in the tree.
func (t *LLRB) InsertNoReplace(item uint64) {
	t.root = t.insertNoReplace(t.root, item)
	t.root.Black = true
	t.count++
}

func (t *LLRB) insertNoReplace(h *Node, item uint64) *Node {
	if h == nil {
		return newNode(item)
	}

	h = walkDownRot23(h)

	if less(item, h.Value) {
		h.Left = t.insertNoReplace(h.Left, item)
	} else {
		h.Right = t.insertNoReplace(h.Right, item)
	}

	return walkUpRot23(h)
}

// Rotation driver routines for 2-3 algorithm

func walkDownRot23(h *Node) *Node { return h }

func walkUpRot23(h *Node) *Node {
	if isRed(h.Right) && !isRed(h.Left) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

// Rotation driver routines for 2-3-4 algorithm

func walkDownRot234(h *Node) *Node {
	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

func walkUpRot234(h *Node) *Node {
	if isRed(h.Right) && !isRed(h.Left) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	return h
}

// DeleteMin deletes the minimum element in the tree and returns the
// deleted item or nil otherwise.
func (t *LLRB) DeleteMin() (uint64, bool) {
	var deleted uint64
	var wasDeleted bool
	t.root, deleted, wasDeleted = deleteMin(t.root)
	if t.root != nil {
		t.root.Black = true
	}
	if wasDeleted {
		t.count--
	}
	return deleted, wasDeleted
}

// deleteMin code for LLRB 2-3 trees
func deleteMin(h *Node) (*Node, uint64, bool) {
	if h == nil {
		return nil, 0, false
	}
	if h.Left == nil {
		return nil, h.Value, true
	}

	if !isRed(h.Left) && !isRed(h.Left.Left) {
		h = moveRedLeft(h)
	}

	var deleted uint64
	var wasDeleted bool
	h.Left, deleted, wasDeleted = deleteMin(h.Left)

	return fixUp(h), deleted, wasDeleted
}

// DeleteMax deletes the maximum element in the tree and returns
// the deleted item or nil otherwise
func (t *LLRB) DeleteMax() (uint64, bool) {
	var deleted uint64
	var wasDeleted bool
	t.root, deleted, wasDeleted = deleteMax(t.root)
	if t.root != nil {
		t.root.Black = true
	}
	if wasDeleted {
		t.count--
	}
	return deleted, wasDeleted
}

func deleteMax(h *Node) (*Node, uint64, bool) {
	if h == nil {
		return nil, 0, false
	}
	if isRed(h.Left) {
		h = rotateRight(h)
	}
	if h.Right == nil {
		return nil, h.Value, true
	}
	if !isRed(h.Right) && !isRed(h.Right.Left) {
		h = moveRedRight(h)
	}
	var deleted uint64
	var wasDeleted bool
	h.Right, deleted, wasDeleted = deleteMax(h.Right)

	return fixUp(h), deleted, wasDeleted
}

// Delete deletes an item from the tree whose key equals key.
// The deleted item is return, otherwise nil is returned.
func (t *LLRB) Delete(key uint64) bool {
	var deleted bool

	t.root, deleted = t.delete(t.root, key)
	if t.root != nil {
		t.root.Black = true
	}
	if deleted {
		t.count--
	}
	return deleted
}

func (t *LLRB) delete(h *Node, item uint64) (*Node, bool) {
	var deleted bool

	if h == nil {
		return nil, false
	}
	if less(item, h.Value) {
		if h.Left == nil { // item not present. Nothing to delete
			return h, false
		}
		if !isRed(h.Left) && !isRed(h.Left.Left) {
			h = moveRedLeft(h)
		}
		h.Left, deleted = t.delete(h.Left, item)
	} else {
		if isRed(h.Left) {
			h = rotateRight(h)
		}
		// If @item equals @h.Value and no right children at @h
		if !less(h.Value, item) && h.Right == nil {
			return nil, true
		}
		// PETAR: Added 'h.Right != nil' below
		if h.Right != nil && !isRed(h.Right) && !isRed(h.Right.Left) {
			h = moveRedRight(h)
		}
		// If @item equals @h.Value, and (from above) 'h.Right != nil'
		if !less(h.Value, item) {
			h.Right, h.Value, deleted = deleteMin(h.Right)
			if !deleted {
				panic("logic")
			}
		} else { // Else, @item is bigger than @h.Value
			h.Right, deleted = t.delete(h.Right, item)
		}
	}

	return fixUp(h), deleted
}

// Internal node manipulation routines

func newNode(item uint64) *Node { return &Node{Value: item} }

func isRed(h *Node) bool {
	if h == nil {
		return false
	}
	return !h.Black
}

func rotateLeft(h *Node) *Node {
	x := h.Right
	if x.Black {
		panic("rotating a black link")
	}
	h.Right = x.Left
	x.Left = h
	x.Black = h.Black
	h.Black = false
	return x
}

func rotateRight(h *Node) *Node {
	x := h.Left
	if x.Black {
		panic("rotating a black link")
	}
	h.Left = x.Right
	x.Right = h
	x.Black = h.Black
	h.Black = false
	return x
}

// REQUIRE: Left and Right children must be present
func flip(h *Node) {
	h.Black = !h.Black
	h.Left.Black = !h.Left.Black
	h.Right.Black = !h.Right.Black
}

// REQUIRE: Left and Right children must be present
func moveRedLeft(h *Node) *Node {
	flip(h)
	if isRed(h.Right.Left) {
		h.Right = rotateRight(h.Right)
		h = rotateLeft(h)
		flip(h)
	}
	return h
}

// REQUIRE: Left and Right children must be present
func moveRedRight(h *Node) *Node {
	flip(h)
	if isRed(h.Left.Left) {
		h = rotateRight(h)
		flip(h)
	}
	return h
}

func fixUp(h *Node) *Node {
	if isRed(h.Right) {
		h = rotateLeft(h)
	}

	if isRed(h.Left) && isRed(h.Left.Left) {
		h = rotateRight(h)
	}

	if isRed(h.Left) && isRed(h.Right) {
		flip(h)
	}

	return h
}

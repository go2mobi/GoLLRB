// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

// GetHeight() returns an item in the tree with key @key, and it's height in the tree
func (t *LLRB) GetHeight(key uint64) (result uint64, depth int) {
	return t.getHeight(t.root, key)
}

func (t *LLRB) getHeight(h *Node, item uint64) (uint64, int) {
	if h == nil {
		return 0, 0
	}
	if less(item, h.Value) {
		result, depth := t.getHeight(h.Left, item)
		return result, depth + 1
	}
	if less(h.Value, item) {
		result, depth := t.getHeight(h.Right, item)
		return result, depth + 1
	}
	return h.Value, 0
}

// HeightStats() returns the average and standard deviation of the height
// of elements in the tree
func (t *LLRB) HeightStats() (avg, stddev float64) {
	av := &avgVar{}
	heightStats(t.root, 0, av)
	return av.GetAvg(), av.GetStdDev()
}

func heightStats(h *Node, d int, av *avgVar) {
	if h == nil {
		return
	}
	av.Add(float64(d))
	if h.Left != nil {
		heightStats(h.Left, d+1, av)
	}
	if h.Right != nil {
		heightStats(h.Right, d+1, av)
	}
}

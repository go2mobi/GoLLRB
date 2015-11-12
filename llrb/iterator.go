package llrb

type ItemIterator func(i uint64) bool

//func (t *Tree) Ascend(iterator ItemIterator) {
//	t.AscendGreaterOrEqual(Inf(-1), iterator)
//}

func (t *LLRB) AscendRange(greaterOrEqual, lessThan uint64, iterator ItemIterator) {
	t.ascendRange(t.root, greaterOrEqual, lessThan, iterator)
}

func (t *LLRB) ascendRange(h *Node, inf, sup uint64, iterator ItemIterator) bool {
	if h == nil {
		return true
	}
	if !less(h.Value, sup) {
		return t.ascendRange(h.Left, inf, sup, iterator)
	}
	if less(h.Value, inf) {
		return t.ascendRange(h.Right, inf, sup, iterator)
	}

	if !t.ascendRange(h.Left, inf, sup, iterator) {
		return false
	}

	if !iterator(h.Value) {
		return false
	}

	return t.ascendRange(h.Right, inf, sup, iterator)
}

// AscendGreaterOrEqual will call iterator once for each element greater or equal to
// pivot in ascending order. It will stop whenever the iterator returns false.
func (t *LLRB) AscendGreaterOrEqual(pivot uint64, iterator ItemIterator) {
	t.ascendGreaterOrEqual(t.root, pivot, iterator)
}

func (t *LLRB) ascendGreaterOrEqual(h *Node, pivot uint64, iterator ItemIterator) bool {
	if h == nil {
		return true
	}
	if !less(h.Value, pivot) {
		if !t.ascendGreaterOrEqual(h.Left, pivot, iterator) {
			return false
		}
		if !iterator(h.Value) {
			return false
		}
	}
	return t.ascendGreaterOrEqual(h.Right, pivot, iterator)
}

func (t *LLRB) AscendLessThan(pivot uint64, iterator ItemIterator) {
	t.ascendLessThan(t.root, pivot, iterator)
}

func (t *LLRB) ascendLessThan(h *Node, pivot uint64, iterator ItemIterator) bool {
	if h == nil {
		return true
	}
	if !t.ascendLessThan(h.Left, pivot, iterator) {
		return false
	}
	if !iterator(h.Value) {
		return false
	}
	if less(h.Value, pivot) {
		return t.ascendLessThan(h.Left, pivot, iterator)
	}
	return true
}

// DescendLessOrEqual will call iterator once for each element less than the
// pivot in descending order. It will stop whenever the iterator returns false.
func (t *LLRB) DescendLessOrEqual(pivot uint64, iterator ItemIterator) {
	t.descendLessOrEqual(t.root, pivot, iterator)
}

func (t *LLRB) descendLessOrEqual(h *Node, pivot uint64, iterator ItemIterator) bool {
	if h == nil {
		return true
	}
	if less(h.Value, pivot) || !less(pivot, h.Value) {
		if !t.descendLessOrEqual(h.Right, pivot, iterator) {
			return false
		}
		if !iterator(h.Value) {
			return false
		}
	}
	return t.descendLessOrEqual(h.Left, pivot, iterator)
}

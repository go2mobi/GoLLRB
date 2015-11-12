// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

import (
	"math"
	"math/rand"
	"testing"
)

func TestCases(t *testing.T) {
	tree := New()
	tree.ReplaceOrInsert(1)
	tree.ReplaceOrInsert(1)
	if tree.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !tree.Has(1) {
		t.Errorf("expecting to find key=1")
	}

	tree.Delete(1)
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(1) {
		t.Errorf("not expecting to find key=1")
	}

	tree.Delete(1)
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(1) {
		t.Errorf("not expecting to find key=1")
	}
}

func TestReverseInsertOrder(t *testing.T) {
	tree := New()
	n := uint64(100)
	for i := uint64(0); i < n; i++ {
		tree.ReplaceOrInsert(n - i)
	}
	i := uint64(0)
	tree.AscendGreaterOrEqual(0, func(item uint64) bool {
		i++
		if item != i {
			t.Errorf("bad order: got %d, expect %d", item, i)
		}
		return true
	})
}
func TestRange(t *testing.T) {
	tree := New()
	order := []uint64{
		1, 2, 3, 4, 5, 6, 7, 8, 9,
	}
	for _, i := range order {
		v := i
		tree.ReplaceOrInsert(v)
	}
	k := 2

	if tree.Len() != len(order) {
		t.Fatalf("Incorrect number of elements added")
	}

	tree.AscendRange(3, 6, func(item uint64) bool {
		if k > 5 {
			t.Fatalf("returned more items than expected")
		}
		i1 := order[k]
		i2 := item
		if i1 != i2 {
			t.Errorf("expecting %v, got %v", i1, i2)
		}
		k++
		return true
	})
}

func TestRandomInsertOrder(t *testing.T) {
	tree := New()
	n := 1000
	perm := rand.Perm(n)

	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(uint64(perm[i]))
	}

	j := uint64(0)
	tree.AscendGreaterOrEqual(0, func(item uint64) bool {
		if item != j {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

func TestRandomReplace(t *testing.T) {
	tree := New()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(uint64(perm[i]))
	}
	perm = rand.Perm(n)
	for i := 0; i < n; i++ {
		v := uint64(perm[i])
		if replaced := tree.ReplaceOrInsert(v); !replaced {
			t.Errorf("error replacing")
		}
	}
}

func TestRandomInsertSequentialDelete(t *testing.T) {
	tree := New()
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(uint64(perm[i]))
	}
	for i := 0; i < n; i++ {
		tree.Delete(uint64(i))
	}
}

func TestRandomInsertDeleteNonExistent(t *testing.T) {
	tree := New()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(uint64(perm[i]))
	}
	if tree.Delete(uint64(200)) {
		t.Errorf("deleted non-existent item")
	}
	if tree.Delete(uint64(300)) {
		t.Errorf("deleted non-existent item")
	}
	for i := 0; i < n; i++ {
		if u := tree.Delete(uint64(i)); !u {
			t.Errorf("delete failed")
		}
	}
	if tree.Delete(uint64(200)) {
		t.Errorf("deleted non-existent item")
	}
	if tree.Delete(uint64(300)) {
		t.Errorf("deleted non-existent item")
	}
}

func TestRandomInsertPartialDeleteOrder(t *testing.T) {
	tree := New()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(uint64(perm[i]))
	}
	for i := 1; i < n-1; i++ {
		tree.Delete(uint64(i))
	}
	j := 0
	tree.AscendGreaterOrEqual(uint64(0), func(item uint64) bool {
		switch j {
		case 0:
			if item != uint64(0) {
				t.Errorf("expecting 0")
			}
		case 1:
			if item != uint64(n-1) {
				t.Errorf("expecting %d", n-1)
			}
		}
		j++
		return true
	})
}

func TestRandomInsertStats(t *testing.T) {
	tree := New()
	n := 100000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(uint64(perm[i]))
	}
	avg, _ := tree.HeightStats()
	expAvg := math.Log2(float64(n)) - 1.5
	if math.Abs(avg-expAvg) >= 2.0 {
		t.Errorf("too much deviation from expected average height")
	}
}

func BenchmarkInsert(b *testing.B) {
	tree := New()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(uint64(b.N - i))
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	tree := New()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(uint64(b.N - i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Delete(uint64(i))
	}
}

func BenchmarkDeleteMin(b *testing.B) {
	b.StopTimer()
	tree := New()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(uint64(b.N - i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.DeleteMin()
	}
}

func TestInsertNoReplace(t *testing.T) {
	tree := New()
	n := 1000
	for q := 0; q < 2; q++ {
		perm := rand.Perm(n)
		for i := 0; i < n; i++ {
			tree.InsertNoReplace(uint64(perm[i]))
		}
	}
	j := 0
	tree.AscendGreaterOrEqual(uint64(0), func(item uint64) bool {
		if item != uint64(j/2) {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

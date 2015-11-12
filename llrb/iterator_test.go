package llrb

import (
	"reflect"
	"testing"
)

func TestAscendGreaterOrEqual(t *testing.T) {
	tree := New()
	tree.InsertNoReplace(4)
	tree.InsertNoReplace(6)
	tree.InsertNoReplace(1)
	tree.InsertNoReplace(3)

	var ary []uint64
	tree.AscendGreaterOrEqual(0, func(i uint64) bool {
		ary = append(ary, i)
		return true
	})
	expected := []uint64{1, 3, 4, 6}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.AscendGreaterOrEqual(3, func(i uint64) bool {
		ary = append(ary, i)
		return true
	})
	expected = []uint64{3, 4, 6}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.AscendGreaterOrEqual(2, func(i uint64) bool {
		ary = append(ary, i)
		return true
	})
	expected = []uint64{3, 4, 6}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
}

func TestDescendLessOrEqual(t *testing.T) {
	tree := New()
	tree.InsertNoReplace(4)
	tree.InsertNoReplace(6)
	tree.InsertNoReplace(1)
	tree.InsertNoReplace(3)

	var ary []uint64
	tree.DescendLessOrEqual(10, func(i uint64) bool {
		ary = append(ary, i)
		return true
	})
	expected := []uint64{6, 4, 3, 1}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.DescendLessOrEqual(4, func(i uint64) bool {
		ary = append(ary, i)
		return true
	})
	expected = []uint64{4, 3, 1}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
	ary = nil
	tree.DescendLessOrEqual(5, func(i uint64) bool {
		ary = append(ary, i)
		return true
	})
	expected = []uint64{4, 3, 1}
	if !reflect.DeepEqual(ary, expected) {
		t.Errorf("expected %v but got %v", expected, ary)
	}
}

package test

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
)

var cmpSortSlices = Cmp(cmpopts.SortSlices(func(i, j int) bool {
	return i < j
}))

func TestCmp_Eq(t *testing.T) {
	a := []int{3, 5, 1, 6, 7}
	b := []int{1, 7, 6, 3, 5}
	Eq(t, a, b, cmpSortSlices)
}

func TestCmp_NotEq(t *testing.T) {
	a := []int{3, 5, 1, 6, 0}
	b := []int{1, 7, 6, 3, 5}
	NotEq(t, a, b, cmpSortSlices)
}

func TestCmp_SliceContains(t *testing.T) {
	a := [][]int{{1}, {1, 2}}
	SliceContains(t, a, []int{2, 1}, cmpSortSlices)
}

func TestCmp_SliceNotContains(t *testing.T) {
	a := [][]int{{1}, {1, 2}}
	SliceNotContains(t, a, []int{3, 1}, cmpSortSlices)
}

func TestCmp_MapContainsValues(t *testing.T) {
	m1 := map[string][]int{
		"one": {1, 3, 5, 7},
	}
	MapContainsValues(t, m1, [][]int{{7, 5, 1, 3}}, cmpSortSlices)
}

func TestCmp_MapNotContainsValues(t *testing.T) {
	m1 := map[string][]int{
		"one": {1, 3, 5, 7},
	}
	MapNotContainsValues(t, m1, [][]int{{0, 5, 1, 3}}, cmpSortSlices)
}

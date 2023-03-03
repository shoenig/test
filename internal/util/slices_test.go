package util

import (
	"strconv"
	"testing"
)

func TestCloneSliceFunc(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		result := CloneSliceFunc([]int{}, func(i int) string {
			return strconv.Itoa(i)
		})
		if len(result) > 0 {
			t.Fatal("expected empty slice")
		}
	})

	t.Run("non empty", func(t *testing.T) {
		original := []int{1, 4, 5}
		result := CloneSliceFunc(original, func(i int) string {
			return strconv.Itoa(i)
		})
		if len(result) != 3 {
			t.Fatal("expected length of 3")
		}
		if result[0] != "1" {
			t.Fatal("expected result[0] == 1")
		}
		if result[1] != "4" {
			t.Fatal("expected result[1] == 4")
		}
		if result[2] != "5" {
			t.Fatal("expected result[2] == 5")
		}
	})
}

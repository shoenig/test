// Code generated via scripts/generate.sh. DO NOT EDIT.

package must

import (
	"fmt"
	"strings"
)

var t = new(myT)

// myT is a substitute for testing.T for use in examples
type myT struct{}

func (t *myT) Errorf(s string, args ...any) {
	s = fmt.Sprintf(s, args...)
	fmt.Println(s)
}

func (t *myT) Fatalf(s string, args ...any) {
	s = fmt.Sprintf(s, args...)
	fmt.Println(s)
}

func (t *myT) Helper() {
	// nothing
}

type myContainer[T comparable] struct {
	items map[T]struct{}
}

func newContainer[T comparable](items ...T) *myContainer[T] {
	c := &myContainer[T]{items: make(map[T]struct{})}
	for _, item := range items {
		c.items[item] = struct{}{}
	}
	return c
}

func (c *myContainer[T]) Contains(item T) bool {
	_, exists := c.items[item]
	return exists
}

func (c *myContainer[T]) Empty() bool {
	return len(c.items) == 0
}

func ExampleAscending() {
	nums := []int{1, 3, 4, 4, 9}
	Ascending(t, nums)
	// Output:
}

// AscendingCmp

// AscendingFunc

// AscendingLess

func ExampleBetween() {
	lower, upper := 3, 9
	value := 5
	Between(t, lower, value, upper)
	// Output:
}

func ExampleBetweenExclusive() {
	lower, upper := 2, 8
	value := 4
	BetweenExclusive(t, lower, value, upper)
	// Output:
}

func ExampleContains() {
	// container implements .Contains method
	container := newContainer(2, 4, 6, 8)
	Contains[int](t, 4, container)
	// Output:
}

func ExampleContainsSubset() {
	// container implements .Contains method
	container := newContainer(1, 2, 3, 4, 5, 6)
	ContainsSubset[int](t, []int{2, 4, 6}, container)
	// Output:
}

func ExampleDescending() {
	nums := []int{9, 6, 5, 4, 4, 2, 1}
	Descending(t, nums)
	// Output:
}

func ExampleDescendingCmp() {
	nums := []int{9, 5, 3, 3, 1, -2}
	DescendingCmp(t, nums, func(a, b int) int {
		return a - b
	})
	// Output:
}

func ExampleDescendingFunc() {
	words := []string{"Foo", "baz", "Bar", "AND"}
	DescendingFunc(t, words, func(a, b string) bool {
		lowerA := strings.ToLower(a)
		lowerB := strings.ToLower(b)
		return lowerA < lowerB
	})
	// Output:
}

// DescendingLess

// DirExists

// DirExistsFS

// DirNotExists

// DirNotExistsFS

func ExampleEmpty() {
	// container implements .Empty method
	container := newContainer[string]()
	Empty(t, container)
	// Output:
}

func ExampleEq() {
	actual := "hello"
	Eq(t, "hello", actual)
	// Output:
}

// EqError

// EqFunc

// EqJSON

// EqOp

// Equal

// Error

// ErrorContains

// ErrorIs

// False

// FileContains

// FileContainsFS

// FileExists

// FileExistsFS

// FileMode

// FileModeFS

// FileNotExists

// FileNotExistsFS

// FilePathValid

// Greater

// GreaterEq

// InDelta

// InDeltaSlice

// Len

// Length

// Less

// LessEq

// Lesser

// MapContainsKey
func ExampleMapContainsKey() {
	numbers := map[string]int{"one": 1, "two": 2, "three": 3}
	MapContainsKey(t, numbers, "one")
	// Output:
}

// MapContainsKeys
func ExampleMapContainsKeys() {
	numbers := map[string]int{"one": 1, "two": 2, "three": 3}
	keys := []string{"one", "two"}
	MapContainsKeys(t, numbers, keys)
	// Output:
}

// MapContainsValues
func ExampleMapContainsValues() {
	numbers := map[string]int{"one": 1, "two": 2, "three": 3}
	values := []int{1, 2}
	MapContainsValues(t, numbers, values)
	// Output:
}

// MapContainsValuesEqual

// MapContainsValuesFunc

// MapEmpty
func ExampleMapEmpty() {
	m := make(map[int]int)
	MapEmpty(t, m)
	// Output:
}

// MapEq
func ExampleMapEq() {
	m1 := map[string]int{"one": 1, "two": 2, "three": 3}
	m2 := map[string]int{"one": 1, "two": 2, "three": 3}
	MapEq(t, m1, m2)
	// Output:
}

// MapEqFunc

// MapEqual

// MapLen

// MapNotContainsKey

// MapNotContainsKeys

// MapNotContainsValues

// MapNotContainsValuesEqual

// MapNotContainsValuesFunc

// MapNotEmpty

// Max

// Min

// Negative

// Nil

// NoError

// NonNegative

// NonPositive

// NonZero

// NotContains

// NotEmpty

// NotEq

// NotEqFunc

// NotEqOp

// NotEqual

// NotNil

// One

// Positive

// RegexCompiles

// RegexCompilesPOSIX

// RegexMatch

// Size

// SliceContains

// SliceContainsAll

// SliceContainsEqual

// SliceContainsFunc

// SliceContainsOp

// SliceContainsSubset

// SliceEmpty

// SliceEqFunc

// SliceEqual

// SliceLen

// SliceNotContains

// SliceNotContainsFunc

// SliceNotEmpty

// StrContains

// StrContainsAny

// StrContainsFields

// StrContainsFold

// StrCount

// StrEqFold

// StrHasPrefix

// StrHasSuffix

// StrNotContains
func ExampleSliceNotContains() {
	StrNotContains(t, "public static void main", "def")
	// Output:
}

// StrNotContainsAny

// StrNotContainsFold

// StrNotEqFold

// StrNotHasPrefix

// StructEqual

// True
func ExampleTrue() {
	True(t, true)
	// Output:
}

// UUIDv4

// Unreachable

// ValidJSON

// ValidJSONBytes

// Wait

// Zero

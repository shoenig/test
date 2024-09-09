// Code generated via scripts/generate.sh. DO NOT EDIT.

// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

//go:build unix

package must

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shoenig/test/wait"
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
	return c.Size() == 0
}

func (c *myContainer[T]) Size() int {
	return len(c.items)
}

type score int

func (s score) Less(other score) bool {
	return s < other
}

func (s score) Equal(other score) bool {
	return s == other
}

type scores []score

func (s scores) Min() score {
	min := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] < min {
			min = s[i]
		}
	}
	return min
}

func (s scores) Max() score {
	max := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] > max {
			max = s[i]
		}
	}
	return max
}

func (s scores) Len() int {
	return len(s)
}

type employee struct {
	first string
	last  string
	id    int
}

func (e *employee) Equal(o *employee) bool {
	if e == nil || o == nil {
		return e == o
	}
	switch {
	case e.first != o.first:
		return false
	case e.last != o.last:
		return false
	case e.id != o.id:
		return false
	}
	return true
}

func (e *employee) Copy() *employee {
	return &employee{
		first: e.first,
		last:  e.last,
		id:    e.id,
	}
}

func ExampleAscending() {
	nums := []int{1, 3, 4, 4, 9}
	Ascending(t, nums)
	// Output:
}

func ExampleAscendingCmp() {
	labels := []string{"Fun", "great", "Happy", "joyous"}
	AscendingCmp(t, labels, func(a, b string) int {
		A := strings.ToLower(a)
		B := strings.ToLower(b)
		switch {
		case A == B:
			return 0
		case A < B:
			return -1
		default:
			return 1
		}
	})
	// Output:
}

func ExampleAscendingFunc() {
	labels := []string{"Fun", "great", "Happy", "joyous"}
	AscendingFunc(t, labels, func(a, b string) bool {
		A := strings.ToLower(a)
		B := strings.ToLower(b)
		return A < B
	})
	// Output:
}

func ExampleAscendingLess() {
	nums := []score{4, 6, 7, 9}
	AscendingLess(t, nums)
	// Output:
}

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

func ExampleDescendingLess() {
	nums := []score{9, 6, 3, 1, 0}
	DescendingLess(t, nums)
	// Output:
}

func ExampleDirExists() {
	DirExists(t, "/tmp")
	// Output:
}

func ExampleDirExistsFS() {
	DirExistsFS(t, os.DirFS("/"), "tmp")
	// Output:
}

func ExampleDirNotExists() {
	DirNotExists(t, "/does/not/exist")
	// Output:
}

func ExampleDirNotExistsFS() {
	DirNotExistsFS(t, os.DirFS("/"), "does/not/exist")
	// Output:
}

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

func ExampleEqError() {
	err := errors.New("undefined error")
	EqError(t, err, "undefined error")
	// Output:
}

func ExampleEqFunc() {
	EqFunc(t, "abcd", "dcba", func(a, b string) bool {
		if len(a) != len(b) {
			return false
		}
		l := len(a)
		for i := 0; i < l; i++ {
			if a[i] != b[l-1-i] {
				return false
			}
		}
		return true
	})
	// Output:
}

func ExampleEqJSON() {
	a := `{"foo":"bar","numbers":[1,2,3]}`
	b := `{"numbers":[1,2,3],"foo":"bar"}`
	EqJSON(t, a, b)
	// Output:
}

func ExampleEqOp() {
	EqOp(t, 123, 123)
	// Output:
}

func ExampleEqual() {
	// score implements .Equal method
	Equal(t, score(1000), score(1000))
	// Output:
}

func ExampleError() {
	Error(t, errors.New("error"))
	// Output:
}

func ExampleErrorContains() {
	err := errors.New("error beer not found")
	ErrorContains(t, err, "beer")
	// Output:
}

func ExampleErrorIs() {
	e1 := errors.New("e1")
	e2 := errors.New("e2")
	e3 := errors.New("e3")
	errorChain := errors.Join(e1, e2, e3)
	ErrorIs(t, errorChain, e2)
	// Output:
}

func ExampleErrorAs() {
	e1 := errors.New("e1")
	e2 := FakeError("foo")
	e3 := errors.New("e3")
	errorChain := errors.Join(e1, e2, e3)
	var target FakeError
	ErrorAs(t, errorChain, &target)
	fmt.Println(target.Error())
	// Output: foo
}

func ExampleFalse() {
	False(t, 1 == int('a'))
	// Output:
}

func ExampleFileContains() {
	_ = os.WriteFile("/tmp/example", []byte("foo bar baz"), fs.FileMode(0600))
	FileContains(t, "/tmp/example", "bar")
	// Output:
}

func ExampleFileContainsFS() {
	_ = os.WriteFile("/tmp/example", []byte("foo bar baz"), fs.FileMode(0600))
	FileContainsFS(t, os.DirFS("/tmp"), "example", "bar")
	// Output:
}

func ExampleFileExists() {
	_ = os.WriteFile("/tmp/example", []byte{}, fs.FileMode(0600))
	FileExists(t, "/tmp/example")
	// Output:
}

func ExampleFileExistsFS() {
	_ = os.WriteFile("/tmp/example", []byte{}, fs.FileMode(0600))
	FileExistsFS(t, os.DirFS("/tmp"), "example")
	// Output:
}

func ExampleFileMode() {
	_ = os.WriteFile("/tmp/example_fm", []byte{}, fs.FileMode(0600))
	FileMode(t, "/tmp/example_fm", fs.FileMode(0600))
	// Output:
}

func ExampleFileModeFS() {
	_ = os.WriteFile("/tmp/example_fm", []byte{}, fs.FileMode(0600))
	FileModeFS(t, os.DirFS("/tmp"), "example_fm", fs.FileMode(0600))
	// Output:
}

func ExampleFileNotExists() {
	FileNotExists(t, "/tmp/not_existing_file")
	// Output:
}

func ExampleFileNotExistsFS() {
	FileNotExistsFS(t, os.DirFS("/tmp"), "not_existing_file")
	// Output:
}

func ExampleFilePathValid() {
	FilePathValid(t, "foo/bar/baz")
	// Output:
}

func ExampleGreater() {
	Greater(t, 30, 42)
	// Output:
}

func ExampleGreaterEq() {
	GreaterEq(t, 30.1, 30.3)
	// Output:
}

func ExampleInDelta() {
	InDelta(t, 30.5, 30.54, .1)
	// Output:
}

func ExampleInDeltaSlice() {
	nums := []int{51, 48, 55, 49, 52}
	base := []int{52, 44, 51, 51, 47}
	InDeltaSlice(t, nums, base, 5)
	// Output:
}

func ExampleLen() {
	nums := []int{1, 3, 5, 9}
	Len(t, 4, nums)
	// Output:
}

func ExampleLength() {
	s := scores{89, 93, 91, 99, 88}
	Length(t, 5, s)
	// Output:
}

func ExampleLess() {
	// compare using < operator
	s := score(50)
	Less(t, 66, s)
	// Output:
}

func ExampleLessEq() {
	s := score(50)
	LessEq(t, 50, s)
	// Output:
}

func ExampleLesser() {
	// compare using .Less method
	s := score(50)
	Lesser(t, 66, s)
	// Output:
}

func ExampleMapContainsKey() {
	numbers := map[string]int{"one": 1, "two": 2, "three": 3}
	MapContainsKey(t, numbers, "one")
	// Output:
}

func ExampleMapContainsKeys() {
	numbers := map[string]int{"one": 1, "two": 2, "three": 3}
	keys := []string{"one", "two"}
	MapContainsKeys(t, numbers, keys)
	// Output:
}

func ExampleMapContainsValues() {
	numbers := map[string]int{"one": 1, "two": 2, "three": 3}
	values := []int{1, 2}
	MapContainsValues(t, numbers, values)
	// Output:
}

func ExampleMapContainsValuesEqual() {
	// employee implements .Equal
	m := map[int]*employee{
		0: {first: "armon", id: 101},
		1: {first: "mitchell", id: 100},
		2: {first: "dave", id: 102},
	}
	expect := []*employee{
		{first: "armon", id: 101},
		{first: "dave", id: 102},
	}
	MapContainsValuesEqual(t, m, expect)
	// Output:
}

func ExampleMapContainsValuesFunc() {
	m := map[int]string{
		0: "Zero",
		1: "ONE",
		2: "two",
	}
	f := func(a, b string) bool {
		return strings.EqualFold(a, b)
	}
	MapContainsValuesFunc(t, m, []string{"one", "two"}, f)
	// Output:
}

func ExampleMapEmpty() {
	m := make(map[int]int)
	MapEmpty(t, m)
	// Output:
}

func ExampleMapEq() {
	m1 := map[string]int{"one": 1, "two": 2, "three": 3}
	m2 := map[string]int{"one": 1, "two": 2, "three": 3}
	MapEq(t, m1, m2)
	// Output:
}

func ExampleMapEqFunc() {
	m1 := map[int]string{
		0: "Zero",
		1: "one",
		2: "TWO",
	}
	m2 := map[int]string{
		0: "ZERO",
		1: "ONE",
		2: "TWO",
	}
	MapEqFunc(t, m1, m2, func(a, b string) bool {
		return strings.EqualFold(a, b)
	})
	// Output:
}

func ExampleMapEqual() {
	armon := &employee{first: "armon", id: 101}
	mitchell := &employee{first: "mitchell", id: 100}
	m1 := map[int]*employee{
		0: mitchell,
		1: armon,
	}
	m2 := map[int]*employee{
		0: mitchell,
		1: armon,
	}
	MapEqual(t, m1, m2)
	// Output:
}

func ExampleMapEqOp() {
	m1 := map[int]string{
		1: "one",
		2: "two",
	}
	m2 := map[int]string{
		1: "one",
		2: "two",
	}
	MapEqOp(t, m1, m2)
	// Output:
}

func ExampleMapLen() {
	m := map[int]string{
		1: "one",
		2: "two",
	}
	MapLen(t, 2, m)
	// Output:
}

func ExampleMapNotContainsKey() {
	m := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	MapNotContainsKey(t, m, "four")
	// Output:
}

func ExampleMapNotContainsKeys() {
	m := map[string]int{
		"one": 1,
		"two": 2,
	}
	MapNotContainsKeys(t, m, []string{"three", "four"})
	// Output:
}

func ExampleMapNotContainsValues() {
	m := map[int]string{
		1: "one",
		2: "two",
	}
	MapNotContainsValues(t, m, []string{"three", "four"})
	// Output:
}

func ExampleMapNotContainsValuesEqual() {
	m := map[int]*employee{
		0: {first: "mitchell", id: 100},
		1: {first: "armon", id: 101},
	}
	MapNotContainsValuesEqual(t, m, []*employee{
		{first: "dave", id: 103},
	})
	// Output:
}

func ExampleMapNotContainsValuesFunc() {
	m := map[int]string{
		1: "One",
		2: "TWO",
		3: "three",
	}
	f := func(a, b string) bool {
		return strings.EqualFold(a, b)
	}
	MapNotContainsValuesFunc(t, m, []string{"four", "five"}, f)
	// Output:
}

func ExampleMapNotEmpty() {
	m := map[string]int{
		"one": 1,
	}
	MapNotEmpty(t, m)
	// Output:
}

func ExampleMax() {
	s := scores{89, 88, 91, 90, 87}
	Max[score](t, 91, s)
	// Output:
}

func ExampleMin() {
	s := scores{89, 88, 90, 91}
	Min[score](t, 88, s)
	// Output:
}

func ExampleNegative() {
	Negative(t, -9)
	// Output:
}

func ExampleNil() {
	var e *employee
	Nil(t, e)
	// Output:
}

func ExampleNoError() {
	var err error
	NoError(t, err)
	// Output:
}

func ExampleNonNegative() {
	NonNegative(t, 4)
	// Output:
}

func ExampleNonPositive() {
	NonPositive(t, -3)
	// Output:
}

func ExampleNonZero() {
	NonZero(t, .001)
	// Output:
}

func ExampleNotContains() {
	c := newContainer("mage", "warrior", "priest", "paladin", "hunter")
	NotContains[string](t, "rogue", c)
	// Output:
}

func ExampleNotEmpty() {
	c := newContainer("one", "two", "three")
	NotEmpty(t, c)
	// Output:
}

func ExampleNotEq() {
	NotEq(t, "one", "two")
	// Output:
}

func ExampleNotEqFunc() {
	NotEqFunc(t, 4.1, 5.2, func(a, b float64) bool {
		return math.Round(a) == math.Round(b)
	})
	// Output:
}

func ExampleNotEqOp() {
	NotEqOp(t, 1, 2)
	// Output:
}

func ExampleNotEqual() {
	e1 := &employee{first: "alice"}
	e2 := &employee{first: "bob"}
	NotEqual(t, e1, e2)
	// Output:
}

func ExampleNotNil() {
	e := &employee{first: "bob"}
	NotNil(t, e)
	// Output:
}

func ExampleOne() {
	One(t, 1)
	// Output:
}

func ExamplePositive() {
	Positive(t, 42)
	// Output:
}

func ExampleRegexCompiles() {
	RegexCompiles(t, `[a-z]{7}`)
	// Output:
}

func ExampleRegexCompilesPOSIX() {
	RegexCompilesPOSIX(t, `[a-z]{3}`)
	// Output:
}

func ExampleRegexMatch() {
	re := regexp.MustCompile(`[a-z]{6}`)
	RegexMatch(t, re, "cookie")
	// Output:
}

func ExampleSize() {
	c := newContainer("pie", "brownie", "cake", "cookie")
	Size(t, 4, c)
	// Output:
}

func ExampleSliceContains() {
	drinks := []string{"ale", "lager", "cider", "wine"}
	SliceContains(t, drinks, "cider")
	// Output:
}

func ExampleSliceContainsAll() {
	nums := []int{2, 4, 6, 7, 8}
	SliceContainsAll(t, nums, []int{7, 8, 2, 6, 4})
	// Output:
}

func ExampleSliceContainsEqual() {
	dave := &employee{first: "dave", id: 8}
	armon := &employee{first: "armon", id: 2}
	mitchell := &employee{first: "mitchell", id: 1}
	employees := []*employee{dave, armon, mitchell}
	SliceContainsEqual(t, employees, &employee{first: "dave", id: 8})
	// Output:
}

func ExampleSliceContainsFunc() {
	// comparing slice to element of same type
	words := []string{"UP", "DoWn", "LefT", "RiGHT"}
	SliceContainsFunc(t, words, "left", func(a, b string) bool {
		return strings.EqualFold(a, b)
	})

	// comparing slice to element of different type
	nums := []string{"2", "4", "6", "8"}
	SliceContainsFunc(t, nums, 4, func(a string, b int) bool {
		return a == strconv.Itoa(b)
	})
	// Output:
}

func ExampleSliceContainsOp() {
	nums := []int{1, 2, 3, 4, 5}
	SliceContainsOp(t, nums, 3)
	// Output:
}

func ExampleSliceContainsSubset() {
	nums := []int{10, 20, 30, 40, 50}
	SliceContainsSubset(t, nums, []int{40, 10, 30})
	// Output:
}

func ExampleSliceEmpty() {
	var ints []int
	SliceEmpty(t, ints)
	// Output:
}

func ExampleSliceEqFunc() {
	ints := []int{2, 4, 6}
	strings := []string{"2", "4", "6"}
	SliceEqFunc(t, ints, strings, func(exp string, value int) bool {
		return strconv.Itoa(value) == exp
	})
	// Output:
}

func ExampleSliceEqual() {
	// type employee implements .Equal
	dave := &employee{first: "dave"}
	armon := &employee{first: "armon"}
	mitchell := &employee{first: "mitchell"}
	s1 := []*employee{dave, armon, mitchell}
	s2 := []*employee{dave, armon, mitchell}
	SliceEqual(t, s1, s2)
	// Output:
}

func ExampleSliceEqOp() {
	s1 := []int{1, 3, 3, 7}
	s2 := []int{1, 3, 3, 7}
	SliceEqOp(t, s1, s2)
	// Output:
}

func ExampleSliceLen() {
	SliceLen(t, 4, []float64{32, 1.2, 0.01, 9e4})
	// Output:
}

func ExampleSliceNotContains() {
	SliceNotContains(t, []int{1, 2, 4, 5}, 3)
	// Output:
}

func ExampleSliceNotContainsFunc() {
	// comparing slice to element of same type
	f := func(a, b int) bool {
		return a == b
	}
	SliceNotContainsFunc(t, []int{10, 20, 30}, 50, f)

	// comparing slice to element of different type
	g := func(s string, b int) bool {
		return strconv.Itoa(b) == s
	}
	SliceNotContainsFunc(t, []string{"1", "2", "3"}, 5, g)
	//Output:
}

func ExampleSliceNotEmpty() {
	SliceNotEmpty(t, []int{2, 4, 6, 8})
	// Output:
}

func ExampleStrContains() {
	StrContains(t, "Visit https://github.com today!", "https://")
	// Output:
}

func ExampleStrContainsAny() {
	StrContainsAny(t, "glyph", "aeiouy")
	// Output:
}

func ExampleStrContainsFields() {
	StrContainsFields(t, "apple banana cherry grape strawberry", []string{"banana", "grape"})
	// Output:
}

func ExampleStrContainsFold() {
	StrContainsFold(t, "one two three", "TWO")
	// Output:
}

func ExampleStrCount() {
	StrCount(t, "see sally sell sea shells by the sea shore", "se", 4)
	// Output:
}

func ExampleStrEqFold() {
	StrEqFold(t, "So MANY test Cases!", "so many test cases!")
	// Output:
}

func ExampleStrHasPrefix() {
	StrHasPrefix(t, "hello", "hello world!")
	// Output:
}

func ExampleStrHasSuffix() {
	StrHasSuffix(t, "world!", "hello world!")
	// Output:
}

func ExampleStrNotContains() {
	StrNotContains(t, "public static void main", "def")
	// Output:
}

func ExampleStrNotContainsAny() {
	StrNotContainsAny(t, "The quick brown fox", "alyz")
	// Output:
}

func ExampleStrNotContainsFold() {
	StrNotContainsFold(t, "This is some text.", "Absent")
	// Output:
}

func ExampleStrNotEqFold() {
	StrNotEqFold(t, "This Is SOME text.", "THIS is some TEXT!")
	// Output:
}

func ExampleStrNotHasPrefix() {
	StrNotHasPrefix(t, "public static void main", "private")
	// Output:
}

func ExampleStructEqual() {
	original := &employee{
		first: "mitchell",
		last:  "hashimoto",
		id:    1,
	}
	StructEqual(t, original, Tweaks[*employee]{{
		Field: "first",
		Apply: func(e *employee) { e.first = "modified" },
	}, {
		Field: "last",
		Apply: func(e *employee) { e.last = "modified" },
	}, {
		Field: "id",
		Apply: func(e *employee) { e.id = 999 },
	}})
	// Output:
}

func ExampleTrue() {
	True(t, true)
	// Output:
}

func ExampleUUIDv4() {
	UUIDv4(t, "60bf6bb2-dceb-c986-2d47-07ac5d14f247")
	// Output:
}

func ExampleUnreachable() {
	if "foo" < "bar" {
		Unreachable(t)
	}
	// Output:
}

func ExampleValidJSON() {
	js := `{"key": ["v1", "v2"]}`
	ValidJSON(t, js)
	// Output:
}

func ExampleValidJSONBytes() {
	js := []byte(`{"key": ["v1", "v2"]}`)
	ValidJSONBytes(t, js)
	// Output:
}

func ExampleWait_initial_success() {
	Wait(t, wait.InitialSuccess(
		wait.BoolFunc(func() bool {
			// will be retried until returns true
			// or timeout is exceeded
			return true
		}),
		wait.Timeout(1*time.Second),
		wait.Gap(100*time.Millisecond),
	))
	// Output:
}

func ExampleWait_continual_success() {
	Wait(t, wait.ContinualSuccess(
		wait.BoolFunc(func() bool {
			// will be retried until timeout expires
			// and will fail test if false is ever returned
			return true
		}),
		wait.Timeout(1*time.Second),
		wait.Gap(100*time.Millisecond),
	))
	// Output:
}

func ExampleZero() {
	Zero(t, 0)
	Zero(t, 0.0)
	// Output:
}

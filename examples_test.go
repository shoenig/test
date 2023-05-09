//go:build unix

package test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
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
	return len(c.items) == 0
}

type score int

func (s score) Less(other score) bool {
	return s < other
}

func (s score) Equal(other score) bool {
	return s == other
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

// DirNotExists
func ExampleDirNotExists() {
	DirNotExists(t, "/does/not/exist")
	// Output:
}

// DirNotExistsFS
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

// EqError
func ExampleEqError() {
	err := errors.New("undefined error")
	EqError(t, err, "undefined error")
	// Output:
}

// EqFunc
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

// EqJSON
func ExampleEqJSON() {
	a := `{"foo":"bar","numbers":[1,2,3]}`
	b := `{"numbers":[1,2,3],"foo":"bar"}`
	EqJSON(t, a, b)
	// Output:
}

// EqOp
func ExampleEqOp() {
	EqOp(t, 123, 123)
	// Output:
}

// Equal
func ExampleEqual() {
	// score implements .Equal method
	Equal(t, score(1000), score(1000))
	// Output:
}

// Error
func ExampleError() {
	Error(t, errors.New("error"))
	// Output:
}

// ErrorContains
func ExampleErrorContains() {
	err := errors.New("error beer not found")
	ErrorContains(t, err, "beer")
	// Output:
}

// ErrorIs
func ExampleErrorIs() {
	e1 := errors.New("e1")
	e2 := errors.New("e2")
	e3 := errors.New("e3")
	errorChain := errors.Join(e1, e2, e3)
	ErrorIs(t, errorChain, e2)
	// Output:
}

// False
func ExampleFalse() {
	False(t, 1 == int('a'))
	// Output:
}

// FileContains
func ExampleFileContains() {
	_ = os.WriteFile("/tmp/example", []byte("foo bar baz"), fs.FileMode(0600))
	FileContains(t, "/tmp/example", "bar")
	// Output:
}

// FileContainsFS
func ExampleFileContainsFS() {
	_ = os.WriteFile("/tmp/example", []byte("foo bar baz"), fs.FileMode(0600))
	FileContainsFS(t, os.DirFS("/tmp"), "example", "bar")
	// Output:
}

// FileExists
func ExampleFileExists() {
	_ = os.WriteFile("/tmp/example", []byte{}, fs.FileMode(0600))
	FileExists(t, "/tmp/example")
	// Output:
}

// FileExistsFS
func ExampleFileExistsFS() {
	_ = os.WriteFile("/tmp/example", []byte{}, fs.FileMode(0600))
	FileExistsFS(t, os.DirFS("/tmp"), "example")
	// Output:
}

// FileMode
func ExampleFileMode() {
	_ = os.WriteFile("/tmp/example_fm", []byte{}, fs.FileMode(0600))
	FileMode(t, "/tmp/example_fm", fs.FileMode(0600))
	// Output:
}

// FileModeFS
func ExampleFileModeFS() {
	_ = os.WriteFile("/tmp/example_fm", []byte{}, fs.FileMode(0600))
	FileModeFS(t, os.DirFS("/tmp"), "example_fm", fs.FileMode(0600))
	// Output:
}

// FileNotExists
func ExampleFileNotExists() {
	FileNotExists(t, "/tmp/not_existing_file")
	// Output:
}

// FileNotExistsFS
func ExampleFileNotExistsFS() {
	FileNotExistsFS(t, os.DirFS("/tmp"), "not_existing_file")
	// Output:
}

// FilePathValid
func ExampleFilePathValid() {
	FilePathValid(t, "foo/bar/baz")
	// Output:
}

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

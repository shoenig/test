//go:build unix

package test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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
	f, _ := os.CreateTemp("", "example")
	defer f.Close()
	f.Write([]byte("foo bar baz"))

	FileContains(t, f.Name(), "bar")
	// Output:
}

// FileContainsFS
func ExampleFileContainsFS() {
	f, _ := os.CreateTemp("", "example")
	fName := filepath.Base(f.Name())
	defer f.Close()
	f.Write([]byte("foo bar baz"))

	FileContainsFS(t, os.DirFS(os.TempDir()), fName, "bar")
	// Output:
}

// FileExists
func ExampleFileExists() {
	f, _ := os.CreateTemp("", "example")
	defer f.Close()

	FileExists(t, f.Name())
	// Output:
}

// FileExistsFS
func ExampleFileExistsFS() {
	f, _ := os.CreateTemp("", "example")
	fName := filepath.Base(f.Name())
	defer f.Close()

	FileExistsFS(t, os.DirFS(os.TempDir()), fName)
	// Output:
}

// FileMode
func ExampleFileMode() {
	f, _ := os.CreateTemp("", "example")
	defer f.Close()

	FileMode(t, f.Name(), fs.FileMode(0600))
	// Output:
}

// FileModeFS
func ExampleFileModeFS() {
	f, _ := os.CreateTemp("", "example")
	fName := filepath.Base(f.Name())
	defer f.Close()

	FileModeFS(t, os.DirFS(os.TempDir()), fName, fs.FileMode(0600))
	// Output:
}

// FileNotExists
func ExampleFileNotExists() {
	FileNotExists(t, "/tmp/not_existing_file")
	// Output:
}

// FileNotExistsFS
func ExampleFileNotExistsFS() {
	FileNotExistsFS(t, os.DirFS(os.TempDir()), "not_existing_file")
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

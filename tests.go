// Package test provides a modern generic testing assertions library.
package test

import (
	"fmt"
	"strings"

	"github.com/shoenig/test/interfaces"
	"github.com/shoenig/test/internal/assertions"
	"github.com/shoenig/test/internal/constraints"
)

// T is the minimal set of functions to be implemented by any testing framework
// compatible with the test package.
type T interface {
	Helper()
	Errorf(string, ...any)
}

const pass = ""

func fail(t T, msg string, args ...any) {
	c := assertions.Caller()
	s := c + fmt.Sprintf(msg, args...)
	t.Errorf("\n" + strings.TrimSpace(s) + "\n")
}

func invoke(t T, result string) {
	if result != pass {
		fail(t, result)
	}
}

// Nil asserts a is nil.
func Nil(t T, a any) {
	t.Helper()
	invoke(t, assertions.Nil(a))
}

// NotNil asserts a is not nil.
func NotNil(t T, a any) {
	t.Helper()
	invoke(t, assertions.NotNil(a))
}

// True asserts that condition is true.
func True(t T, condition bool) {
	t.Helper()
	invoke(t, assertions.True(condition))
}

// False asserts condition is false.
func False(t T, condition bool) {
	t.Helper()
	invoke(t, assertions.False(condition))
}

// Error asserts err is a non-nil error.
func Error(t T, err error) {
	t.Helper()
	invoke(t, assertions.Error(err))
}

// EqError asserts err contains message msg.
func EqError(t T, err error, msg string) {
	t.Helper()
	invoke(t, assertions.EqError(err, msg))
}

// ErrorIs asserts err
func ErrorIs(t T, err error, target error) {
	t.Helper()
	invoke(t, assertions.ErrorIs(err, target))
}

// NoError asserts err is a nil error.
func NoError(t T, err error) {
	t.Helper()
	invoke(t, assertions.NoError(err))
}

// Eq asserts a and b are equal using cmp.Equal.
func Eq[A any](t T, a, b A) {
	t.Helper()
	invoke(t, assertions.Eq(a, b))
}

// EqOp asserts a == b.
func EqOp[C comparable](t T, a, b C) {
	t.Helper()
	invoke(t, assertions.EqOp(a, b))
}

// EqFunc asserts a and b are equal using eq.
func EqFunc[A any](t T, a, b A, eq func(a, b A) bool) {
	t.Helper()
	invoke(t, assertions.EqFunc(a, b, eq))
}

// NotEq asserts a != b.
func NotEq[C comparable](t T, a, b C) {
	t.Helper()
	invoke(t, assertions.NotEq(a, b))
}

// NotEqFunc asserts a and b are not equal using eq.
func NotEqFunc[A any](t T, a, b A, eq func(a, b A) bool) {
	t.Helper()
	invoke(t, assertions.NotEqFunc(a, b, eq))
}

// EqJSON asserts a and b are equivalent JSON.
func EqJSON(t T, a, b string) {
	t.Helper()
	invoke(t, assertions.EqJSON(a, b))
}

// EqSliceFunc asserts elements of a and b are the same using eq.
func EqSliceFunc[A any](t T, a, b []A, eq func(a, b A) bool) {
	t.Helper()
	invoke(t, assertions.EqSliceFunc(a, b, eq))
}

// Equals asserts a.Equals(b).
func Equals[E interfaces.EqualsFunc[E]](t T, a, b E) {
	t.Helper()
	invoke(t, assertions.Equals(a, b))
}

// NotEquals asserts !a.Equals(b).
func NotEquals[E interfaces.EqualsFunc[E]](t T, a, b E) {
	t.Helper()
	invoke(t, assertions.NotEquals(a, b))
}

// EqualsSlice asserts a[n].Equals(b[n]) for each element n in slices a and b.
func EqualsSlice[E interfaces.EqualsFunc[E]](t T, a, b []E) {
	t.Helper()
	invoke(t, assertions.EqualsSlice(a, b))
}

// Lesser asserts a.Less(b).
func Lesser[L interfaces.LessFunc[L]](t T, a, b L) {
	t.Helper()
	invoke(t, assertions.Lesser(a, b))
}

// EmptySlice asserts slice is empty.
func EmptySlice[A any](t T, slice []A) {
	t.Helper()
	invoke(t, assertions.EmptySlice(slice))
}

// LenSlice asserts slice is of length n.
func LenSlice[A any](t T, n int, slice []A) {
	t.Helper()
	invoke(t, assertions.LenSlice(n, slice))
}

// Contains asserts item exists in slice using cmp.Equal function.
func Contains[A any](t T, slice []A, item A) {
	t.Helper()
	invoke(t, assertions.Contains(slice, item))
}

// ContainsOp asserts item exists in slice using == operator.
func ContainsOp[C comparable](t T, slice []C, item C) {
	t.Helper()
	invoke(t, assertions.ContainsOp(slice, item))
}

// ContainsFunc asserts item exists in slice, using eq to compare elements.
func ContainsFunc[A any](t T, slice []A, item A, eq func(a, b A) bool) {
	t.Helper()
	invoke(t, assertions.ContainsFunc(slice, item, eq))
}

// ContainsEquals asserts item exists in slice, using Equals to compare elements.
func ContainsEquals[E interfaces.EqualsFunc[E]](t T, slice []E, item E) {
	t.Helper()
	invoke(t, assertions.ContainsEquals(slice, item))
}

// Less asserts a < b.
func Less[O constraints.Ordered](t T, a, b O) {
	t.Helper()
	invoke(t, assertions.Less(a, b))
}

// LessEq asserts a <= b.
func LessEq[O constraints.Ordered](t T, a, b O) {
	t.Helper()
	invoke(t, assertions.LessEq(a, b))
}

// Greater asserts a > b.
func Greater[O constraints.Ordered](t T, a, b O) {
	t.Helper()
	invoke(t, assertions.Greater(a, b))
}

// GreaterEq asserts a >= b.
func GreaterEq[O constraints.Ordered](t T, a, b O) {
	t.Helper()
	invoke(t, assertions.GreaterEq(a, b))
}

// InDelta asserts a and b are within delta of each other.
func InDelta[N interfaces.Number](t T, a, b, delta N) {
	t.Helper()
	invoke(t, assertions.InDelta(a, b, delta))
}

// InDeltaSlice asserts each element a[n] is within delta of b[n].
func InDeltaSlice[N interfaces.Number](t T, a, b []N, delta N) {
	t.Helper()
	invoke(t, assertions.InDeltaSlice(a, b, delta))
}

// MapEq asserts maps a and b contain the same key/value pairs, using
// cmp.Equal function to compare values.
func MapEq[M1, M2 interfaces.Map[K, V], K comparable, V any](t T, a M1, b M2) {
	t.Helper()
	invoke(t, assertions.MapEq(a, b))
}

// MapEqFunc asserts maps a and b contain the same key/value pairs, using eq to
// compare values.
func MapEqFunc[M1, M2 interfaces.Map[K, V], K comparable, V any](t T, a M1, b M2, eq func(V, V) bool) {
	t.Helper()
	invoke(t, assertions.MapEqFunc(a, b, eq))
}

// MapEquals asserts maps a and b contain the same key/value pairs, using Equals
// method to compare values
func MapEquals[M interfaces.MapEqualsFunc[K, V], K comparable, V interfaces.EqualsFunc[V]](t T, a, b M) {
	t.Helper()
	invoke(t, assertions.MapEquals(a, b))
}

// MapLen asserts map is of size n.
func MapLen[M map[K]V, K comparable, V any](t T, n int, m M) {
	t.Helper()
	invoke(t, assertions.MapLen(n, m))
}

// MapEmpty asserts map is empty.
func MapEmpty[M map[K]V, K comparable, V any](t T, m M) {
	t.Helper()
	invoke(t, assertions.MapEmpty(m))
}

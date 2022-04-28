package must

import (
	"github.com/shoenig/test"
	"github.com/shoenig/test/internal/constraints"
)

// T is like test.T except implements FailNow instead of Fail.
type T interface {
	Helper()
	FailNow()
	Logf(string, ...any)
}

// shim is an implementation of test.T, which converts the
// call to Fail() to a call to FailNow(), stopping the test
// case immediately
type shim struct {
	t T
}

func (s *shim) Helper() {
	s.t.Helper()
}

func (s *shim) Fail() {
	s.t.FailNow()
}

func (s *shim) Logf(msg string, args ...any) {
	s.t.Logf(msg, args...)
}

// convert will translate T into test.T so we can re-use the test package implementations.
func convert(t T) test.T {
	return &shim{t}
}

// Nil asserts a is nil.
func Nil(t T, a any) {
	t.Helper()
	test.Nil(convert(t), a)
}

// NotNil asserts a is not nil.
func NotNil(t T, a any) {
	t.Helper()
	test.NotNil(convert(t), a)
}

// True asserts that condition is true.
func True(t T, condition bool) {
	t.Helper()
	test.True(convert(t), condition)
}

// False asserts condition is false.
func False(t T, condition bool) {
	t.Helper()
	test.False(convert(t), condition)
}

// Error asserts err is a non-nil error.
func Error(t T, err error) {
	t.Helper()
	test.Error(convert(t), err)
}

// EqError asserts error contains value msg.
func EqError(t T, err error, msg string) {
	t.Helper()
	test.EqError(convert(t), err, msg)
}

// ErrorIs asserts err
func ErrorIs(t T, err error, target error) {
	t.Helper()
	test.ErrorIs(convert(t), err, target)
}

// NoError asserts err is a nil error.
func NoError(t T, err error) {
	t.Helper()
	test.NoError(convert(t), err)
}

// Eq asserts a and b are equal using cmp.Equal.
func Eq[A any](t T, a, b A) {
	t.Helper()
	test.Eq(convert(t), a, b)
}

// EqCmp asserts a == b.
func EqCmp[C comparable](t T, a, b C) {
	t.Helper()
	test.EqCmp(convert(t), a, b)
}

// EqFunc asserts a and b are equal using eq.
func EqFunc[A any](t T, a, b A, eq func(a, b A) bool) {
	t.Helper()
	test.EqFunc(convert(t), a, b, eq)
}

// NotEq asserts a != b.
func NotEq[C comparable](t T, a, b C) {
	t.Helper()
	test.NotEq(convert(t), a, b)
}

// NotEqFunc asserts a and b are not equal using eq.
func NotEqFunc[A any](t T, a, b A, eq func(a, b A) bool) {
	t.Helper()
	test.NotEqFunc(convert(t), a, b, eq)
}

// EqJSON asserts a and b are equivalent JSON.
func EqJSON(t T, a, b string) {
	t.Helper()
	test.EqJSON(convert(t), a, b)
}

// EqSliceFunc asserts elements of a and b are the same using eq.
func EqSliceFunc[A any](t T, a, b []A, eq func(a, b A) bool) {
	t.Helper()
	test.EqSliceFunc(convert(t), a, b, eq)
}

// Equals asserts a.Equals(b).
func Equals[E test.EqualsFunc[E]](t T, a, b E) {
	t.Helper()
	test.Equals(convert(t), a, b)
}

// NotEquals asserts !a.Equals(b).
func NotEquals[E test.EqualsFunc[E]](t T, a, b E) {
	t.Helper()
	test.NotEquals(convert(t), a, b)
}

// EqualsSlice asserts a[n].Equals(b[n]) for each element n in slices a and b.
func EqualsSlice[E test.EqualsFunc[E]](t T, a, b []E) {
	t.Helper()
	test.EqualsSlice(convert(t), a, b)
}

// Lesser asserts a.Less(b).
func Lesser[L test.LessFunc[L]](t T, a, b L) {
	t.Helper()
	test.Lesser(convert(t), a, b)
}

// EmptySlice asserts slice is empty.
func EmptySlice[A any](t T, slice []A) {
	t.Helper()
	test.EmptySlice(convert(t), slice)
}

// LenSlice asserts slice is of length n.
func LenSlice[A any](t T, n int, slice []A) {
	t.Helper()
	test.LenSlice(convert(t), n, slice)
}

// Contains asserts item exists in slice using cmp.Equal function.
func Contains[A any](t T, slice []A, item A) {
	t.Helper()
	test.Contains(convert(t), slice, item)
}

// ContainsCmp asserts item exists in slice using == operator.
func ContainsCmp[C comparable](t T, slice []C, item C) {
	t.Helper()
	test.ContainsCmp(convert(t), slice, item)
}

// ContainsFunc asserts item exists in slice, using eq to compare elements.
func ContainsFunc[A any](t T, slice []A, item A, eq func(a, b A) bool) {
	t.Helper()
	test.ContainsFunc(convert(t), slice, item, eq)
}

// ContainsEquals asserts item exists in slice, using Equals to compare elements.
func ContainsEquals[E test.EqualsFunc[E]](t T, slice []E, item E) {
	t.Helper()
	test.ContainsEquals(convert(t), slice, item)
}

// Less asserts a < b.
func Less[O constraints.Ordered](t T, a, b O) {
	t.Helper()
	test.Less(convert(t), a, b)
}

// LessEq asserts a <= b.
func LessEq[O constraints.Ordered](t T, a, b O) {
	t.Helper()
	test.LessEq(convert(t), a, b)
}

// Greater asserts a > b.
func Greater[O constraints.Ordered](t T, a, b O) {
	t.Helper()
	test.Greater(convert(t), a, b)
}

// GreaterEq asserts a >= b.
func GreaterEq[O constraints.Ordered](t T, a, b O) {
	t.Helper()
	test.GreaterEq(convert(t), a, b)
}

// InDelta asserts a and b are within delta of each other.
func InDelta[N test.Number](t T, a, b, delta N) {
	t.Helper()
	test.InDelta(convert(t), a, b, delta)
}

// InDeltaSlice asserts each element a[n] is within delta of b[n].
func InDeltaSlice[N test.Number](t T, a, b []N, delta N) {
	t.Helper()
	test.InDeltaSlice(convert(t), a, b, delta)
}

// MapEq asserts maps a and b contain the same key/value pairs, using
// cmp.Equal function to compare values.
func MapEq[M1, M2 ~map[K]V, K comparable, V any](t T, a M1, b M2) {
	t.Helper()
	test.MapEq[M1, M2, K, V](convert(t), a, b)
}

// MapEqFunc asserts maps a and b contain the same key/value pairs, using eq to
// compare values.
func MapEqFunc[M ~map[K]V, K comparable, V any](t T, a, b M, eq func(V, V) bool) {
	t.Helper()
	test.MapEqFunc[M, K, V](convert(t), a, b, eq)
}

// MapEquals asserts maps a and b contain the same key/value pairs, using Equals
// method to compare values
func MapEquals[M ~map[K]V, K comparable, V test.EqualsFunc[V]](t T, a, b M) {
	t.Helper()
	test.MapEquals[M, K, V](convert(t), a, b)
}

// MapLen asserts map is of size n.
func MapLen[M ~map[K]V, K comparable, V any](t T, n int, m M) {
	t.Helper()
	test.MapLen[M, K, V](convert(t), n, m)
}

// MapEmpty asserts map is empty.
func MapEmpty[M ~map[K]V, K comparable, V any](t T, m M) {
	t.Helper()
	test.MapEmpty[M, K, V](convert(t), m)
}

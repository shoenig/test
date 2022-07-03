package test

import (
	"io/fs"
	"regexp"

	"github.com/shoenig/test/interfaces"
	"github.com/shoenig/test/internal/assertions"
	"github.com/shoenig/test/internal/constraints"
)

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

// NotEq asserts a and b are not equal using cmp.Equal.
func NotEq[A any](t T, a, b A) {
	t.Helper()
	invoke(t, assertions.NotEq(a, b))
}

// NotEqOp asserts a != b.
func NotEqOp[C comparable](t T, a, b C) {
	t.Helper()
	invoke(t, assertions.NotEqOp(a, b))
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

// Empty asserts slice is empty.
//
// Convenience function for EmptySlice.
func Empty[A any](t T, slice []A) {
	t.Helper()
	EmptySlice(t, slice)
}

// LenSlice asserts slice is of length n.
func LenSlice[A any](t T, n int, slice []A) {
	t.Helper()
	invoke(t, assertions.LenSlice(n, slice))
}

// Len asserts slice is of length n.
//
// Convenience function for LenSlice.
func Len[A any](t T, n int, slice []A) {
	t.Helper()
	LenSlice(t, n, slice)
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

// ContainsString asserts s contains sub.
func ContainsString(t T, s, sub string) {
	t.Helper()
	invoke(t, assertions.ContainsString(s, sub))
}

// Positive asserts n > 0.
func Positive[N interfaces.Number](t T, n N) {
	t.Helper()
	invoke(t, assertions.Positive(n))
}

// Negative asserts n < 0.
func Negative[N interfaces.Number](t T, n N) {
	t.Helper()
	invoke(t, assertions.Negative(n))
}

// Zero asserts n == 0.
func Zero[N interfaces.Number](t T, n N) {
	t.Helper()
	invoke(t, assertions.Zero(n))
}

// NonZero asserts n != 0.
func NonZero[N interfaces.Number](t T, n N) {
	t.Helper()
	invoke(t, assertions.NonZero(n))
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

// Ascending asserts slice[n] <= slice[n+1] for each element n.
func Ascending[O constraints.Ordered](t T, slice []O) {
	t.Helper()
	invoke(t, assertions.Ascending(slice))
}

// AscendingFunc asserts slice[n] is less than slice[n+1] for each element n using the less comparator.
func AscendingFunc[A any](t T, slice []A, less func(A, A) bool) {
	t.Helper()
	invoke(t, assertions.AscendingFunc(slice, less))
}

// AscendingLess asserts slice[n].Less(slice[n+1]) for each element n.
func AscendingLess[L interfaces.LessFunc[L]](t T, slice []L) {
	t.Helper()
	invoke(t, assertions.AscendingLess(slice))
}

// Descending asserts slice[n] >= slice[n+1] for each element n.
func Descending[O constraints.Ordered](t T, slice []O) {
	t.Helper()
	invoke(t, assertions.Descending(slice))
}

// DescendingFunc asserts slice[n+1] is less than slice[n] for each element n using the less comparator.
func DescendingFunc[A any](t T, slice []A, less func(A, A) bool) {
	t.Helper()
	invoke(t, assertions.DescendingFunc(slice, less))
}

// DescendingLess asserts slice[n+1].Less(slice[n]) for each element n.
func DescendingLess[L interfaces.LessFunc[L]](t T, slice []L) {
	t.Helper()
	invoke(t, assertions.DescendingLess(slice))
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
	invoke(t, assertions.MapEq[M1, M2, K, V](a, b))
}

// MapEqFunc asserts maps a and b contain the same key/value pairs, using eq to
// compare values.
func MapEqFunc[M1, M2 interfaces.Map[K, V], K comparable, V any](t T, a M1, b M2, eq func(V, V) bool) {
	t.Helper()
	invoke(t, assertions.MapEqFunc[M1, M2, K, V](a, b, eq))
}

// MapEquals asserts maps a and b contain the same key/value pairs, using Equals
// method to compare values
func MapEquals[M interfaces.MapEqualsFunc[K, V], K comparable, V interfaces.EqualsFunc[V]](t T, a, b M) {
	t.Helper()
	invoke(t, assertions.MapEquals[M, K, V](a, b))
}

// MapLen asserts map is of size n.
func MapLen[M ~map[K]V, K comparable, V any](t T, n int, m M) {
	t.Helper()
	invoke(t, assertions.MapLen[M, K, V](n, m))
}

// MapEmpty asserts map is empty.
func MapEmpty[M ~map[K]V, K comparable, V any](t T, m M) {
	t.Helper()
	invoke(t, assertions.MapEmpty[M, K, V](m))
}

// MapContainsKeys asserts m contains each key in keys.
func MapContainsKeys[M ~map[K]V, K comparable, V any](t T, m M, keys []K) {
	t.Helper()
	invoke(t, assertions.MapContainsKeys[M, K, V](m, keys))
}

// MapContainsValues asserts m contains each value in values.
func MapContainsValues[M ~map[K]V, K comparable, V any](t T, m M, values []V) {
	t.Helper()
	invoke(t, assertions.MapContainsValues[M, K, V](m, values))
}

// MapContainsValuesFunc asserts m contains each value in values using the eq function.
func MapContainsValuesFunc[M ~map[K]V, K comparable, V any](t T, m M, values []V, eq func(V, V) bool) {
	t.Helper()
	invoke(t, assertions.MapContainsValuesFunc[M, K, V](m, values, eq))
}

func MapContainsValuesEquals[M ~map[K]V, K comparable, V interfaces.EqualsFunc[V]](t T, m M, values []V) {
	t.Helper()
	invoke(t, assertions.MapContainsValuesEquals[M, K, V](m, values))
}

// FileExists asserts file exists on system.
//
// Often os.DirFS is used to interact with the the host filesystem.
// Example,
// FileExists(t, os.DirFS("/etc"), "hosts")
func FileExists(t T, system fs.FS, file string) {
	t.Helper()
	invoke(t, assertions.FileExists(system, file))
}

// FileNotExists asserts file does not exist on system.
//
// Often os.DirFS is used to interact with the host filesystem.
// Example,
// FileNotExist(t, os.DirFS("/bin"), "exploit.exe")
func FileNotExists(t T, system fs.FS, file string) {
	t.Helper()
	invoke(t, assertions.FileNotExists(system, file))
}

// DirExists asserts directory exists on system.
//
// Often os.DirFS is used to interact with the host filesystem.
// Example,
// DirExists(t, os.DirFS("/usr/local"), "bin")
func DirExists(t T, system fs.FS, directory string) {
	t.Helper()
	invoke(t, assertions.DirExists(system, directory))
}

// DirNotExists asserts directory does not exist on system.
//
// Often os.DirFS is used to interact with the host filesystem.
// Example,
// DirNotExists(t, os.DirFS("/tmp"), "scratch")
func DirNotExists(t T, system fs.FS, directory string) {
	t.Helper()
	invoke(t, assertions.DirNotExists(system, directory))
}

// FileMode asserts the file or directory at path has exactly
// the given permission bits.
//
// Often os.DirFS is used to interact with the host filesystem.
// Example,
// FileMode(t, os.DirFS("/bin"), "find", 0655)
func FileMode(t T, system fs.FS, path string, permissions fs.FileMode) {
	t.Helper()
	invoke(t, assertions.FileMode(system, path, permissions))
}

// FileContains asserts the file contains content as a substring.
//
// Often os.DirFS is used to interact with the host filesystem.
// Example,
// FileContains(t, os.DirFS("/etc"), "hosts", "localhost")
func FileContains(t T, system fs.FS, file, content string) {
	t.Helper()
	invoke(t, assertions.FileContains(system, file, content))
}

// FilePathValid asserts path is a valid file path.
func FilePathValid(t T, path string) {
	t.Helper()
	invoke(t, assertions.FilePathValid(path))
}

// RegexMatch asserts regular expression re matches string s.
func RegexMatch(t T, re *regexp.Regexp, s string) {
	t.Helper()
	invoke(t, assertions.RegexMatch(re, s))
}

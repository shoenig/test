// Code generated via scripts/generate.sh. DO NOT EDIT.

package must

import (
	"io/fs"
	"os"
	"regexp"
	"strings"

	"github.com/shoenig/test/interfaces"
	"github.com/shoenig/test/internal/assertions"
	"github.com/shoenig/test/internal/constraints"
)

// Nil asserts a is nil.
func Nil(t T, a any, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Nil(a), scripts...)
}

// NotNil asserts a is not nil.
func NotNil(t T, a any, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.NotNil(a), scripts...)
}

// True asserts that condition is true.
func True(t T, condition bool, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.True(condition), scripts...)
}

// False asserts condition is false.
func False(t T, condition bool, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.False(condition), scripts...)
}

// Unreachable asserts a code path is not executed.
func Unreachable(t T, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Unreachable(), scripts...)
}

// Error asserts err is a non-nil error.
func Error(t T, err error, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Error(err), scripts...)
}

// EqError asserts err contains message msg.
func EqError(t T, err error, msg string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.EqError(err, msg), scripts...)
}

// ErrorIs asserts err
func ErrorIs(t T, err error, target error, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.ErrorIs(err, target), scripts...)
}

// NoError asserts err is a nil error.
func NoError(t T, err error, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.NoError(err), scripts...)
}

// ErrorContains asserts err contains sub.
func ErrorContains(t T, err error, sub string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.ErrorContains(err, sub), scripts...)
}

// Eq asserts exp and val are equal using cmp.Equal.
func Eq[A any](t T, exp, val A, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Eq(exp, val), scripts...)
}

// EqOp asserts exp == val.
func EqOp[C comparable](t T, exp, val C, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.EqOp(exp, val), scripts...)
}

// EqFunc asserts exp and val are equal using eq.
func EqFunc[A any](t T, exp, val A, eq func(a, b A) bool, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.EqFunc(exp, val, eq), scripts...)
}

// NotEq asserts exp and val are not equal using cmp.Equal.
func NotEq[A any](t T, exp, val A, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.NotEq(exp, val), scripts...)
}

// NotEqOp asserts exp != val.
func NotEqOp[C comparable](t T, exp, val C, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.NotEqOp(exp, val), scripts...)
}

// NotEqFunc asserts exp and val are not equal using eq.
func NotEqFunc[A any](t T, exp, val A, eq func(a, b A) bool, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.NotEqFunc(exp, val, eq), scripts...)
}

// EqJSON asserts exp and val are equivalent JSON.
func EqJSON(t T, exp, val string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.EqJSON(exp, val), scripts...)
}

// Equal asserts val.Equal(exp).
func Equal[E interfaces.EqualFunc[E]](t T, exp, val E, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Equal(exp, val), scripts...)
}

// NotEqual asserts !val.Equal(exp).
func NotEqual[E interfaces.EqualFunc[E]](t T, exp, val E, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.NotEqual(exp, val), scripts...)
}

// Lesser asserts val.Less(exp).
func Lesser[L interfaces.LessFunc[L]](t T, exp, val L, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Lesser(exp, val), scripts...)
}

// SliceEqFunc asserts elements of exp and val are the same using eq.
func SliceEqFunc[A any](t T, exp, val []A, eq func(a, b A) bool, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.EqSliceFunc(exp, val, eq), scripts...)
}

// SliceEqual asserts val[n].Equal(exp[n]) for each element n.
func SliceEqual[E interfaces.EqualFunc[E]](t T, exp, val []E, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceEqual(exp, val), scripts...)
}

// SliceEmpty asserts slice is empty.
func SliceEmpty[A any](t T, slice []A, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceEmpty(slice), scripts...)
}

// SliceNotEmpty asserts slice is not empty.
func SliceNotEmpty[A any](t T, slice []A, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceNotEmpty(slice), scripts...)
}

// SliceLen asserts slice is of length n.
func SliceLen[A any](t T, n int, slice []A, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceLen(n, slice), scripts...)
}

// Len asserts slice is of length n.
//
// Shorthand function for SliceLen. For checking Len() of a struct,
// use the Length() assertion.
func Len[A any](t T, n int, slice []A, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceLen(n, slice), scripts...)
}

// SliceContainsOp asserts item exists in slice using == operator.
func SliceContainsOp[C comparable](t T, slice []C, item C, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceContainsOp(slice, item), scripts...)
}

// SliceContainsFunc asserts item exists in slice, using eq to compare elements.
func SliceContainsFunc[A any](t T, slice []A, item A, eq func(a, b A) bool, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceContainsFunc(slice, item, eq), scripts...)
}

// SliceContainsEqual asserts item exists in slice, using Equal to compare elements.
func SliceContainsEqual[E interfaces.EqualFunc[E]](t T, slice []E, item E, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceContainsEqual(slice, item), scripts...)
}

// SliceContains asserts item exists in slice, using cmp.Equal to compare elements.
func SliceContains[A any](t T, slice []A, item A, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceContains(slice, item), scripts...)
}

// SliceNotContains asserts item does not exist in slice, using cmp.Equal to
// compare elements.
func SliceNotContains[A any](t T, slice []A, item A, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceNotContains(slice, item), scripts...)
}

// SliceContainsAll asserts slice and items contain the same elements, but in
// no particular order. The number of elements in slice and items must be the
// same.
func SliceContainsAll[A any](t T, slice, items []A, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceContainsAll(slice, items), scripts...)
}

// SliceContainsSubset asserts slice contains each item in items, in no particular
// order. There could be additional elements in slice not in items.
func SliceContainsSubset[A any](t T, slice, items []A, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.SliceContainsSubset(slice, items), scripts...)
}

// Positive asserts n > 0.
func Positive[N interfaces.Number](t T, n N, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Positive(n), scripts...)
}

// NonPositive asserts n ≤ 0.
func NonPositive[N interfaces.Number](t T, n N, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.NonPositive(n), scripts...)
}

// Negative asserts n < 0.
func Negative[N interfaces.Number](t T, n N, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Negative(n), scripts...)
}

// NonNegative asserts n >= 0.
func NonNegative[N interfaces.Number](t T, n N, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.NonNegative(n), scripts...)
}

// Zero asserts n == 0.
func Zero[N interfaces.Number](t T, n N, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Zero(n), scripts...)
}

// NonZero asserts n != 0.
func NonZero[N interfaces.Number](t T, n N, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.NonZero(n), scripts...)
}

// One asserts n == 1.
func One[N interfaces.Number](t T, n N, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.One(n), scripts...)
}

// Less asserts val < exp.
func Less[O constraints.Ordered](t T, exp, val O, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Less(exp, val), scripts...)
}

// LessEq asserts val ≤ exp.
func LessEq[O constraints.Ordered](t T, exp, val O, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.LessEq(exp, val), scripts...)
}

// Greater asserts val > exp.
func Greater[O constraints.Ordered](t T, exp, val O, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Greater(exp, val), scripts...)
}

// GreaterEq asserts val ≥ exp.
func GreaterEq[O constraints.Ordered](t T, exp, val O, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.GreaterEq(exp, val), scripts...)
}

// Between asserts lower ≤ val ≤ upper.
func Between[O constraints.Ordered](t T, lower, val, upper O, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Between(lower, val, upper), scripts...)
}

// BetweenExclusive asserts lower < val < upper.
func BetweenExclusive[O constraints.Ordered](t T, lower, val, upper O, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.BetweenExclusive(lower, val, upper), scripts...)
}

// Ascending asserts slice[n] ≤ slice[n+1] for each element n.
func Ascending[O constraints.Ordered](t T, slice []O, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Ascending(slice), scripts...)
}

// AscendingFunc asserts slice[n] is less than slice[n+1] for each element n using the less comparator.
func AscendingFunc[A any](t T, slice []A, less func(A, A) bool, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.AscendingFunc(slice, less), scripts...)
}

// AscendingLess asserts slice[n].Less(slice[n+1]) for each element n.
func AscendingLess[L interfaces.LessFunc[L]](t T, slice []L, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.AscendingLess(slice), scripts...)
}

// Descending asserts slice[n] ≥ slice[n+1] for each element n.
func Descending[O constraints.Ordered](t T, slice []O, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Descending(slice), scripts...)
}

// DescendingFunc asserts slice[n+1] is less than slice[n] for each element n using the less comparator.
func DescendingFunc[A any](t T, slice []A, less func(A, A) bool, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.DescendingFunc(slice, less), scripts...)
}

// DescendingLess asserts slice[n+1].Less(slice[n]) for each element n.
func DescendingLess[L interfaces.LessFunc[L]](t T, slice []L, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.DescendingLess(slice), scripts...)
}

// InDelta asserts a and b are within delta of each other.
func InDelta[N interfaces.Number](t T, a, b, delta N, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.InDelta(a, b, delta), scripts...)
}

// InDeltaSlice asserts each element a[n] is within delta of b[n].
func InDeltaSlice[N interfaces.Number](t T, a, b []N, delta N, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.InDeltaSlice(a, b, delta), scripts...)
}

// MapEq asserts maps exp and val contain the same key/val pairs, using
// cmp.Equal function to compare vals.
func MapEq[M1, M2 interfaces.Map[K, V], K comparable, V any](t T, exp M1, val M2, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapEq(exp, val), scripts...)
}

// MapEqFunc asserts maps exp and val contain the same key/val pairs, using eq to
// compare vals.
func MapEqFunc[M1, M2 interfaces.Map[K, V], K comparable, V any](t T, exp M1, val M2, eq func(V, V) bool, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapEqFunc(exp, val, eq), scripts...)
}

// MapEqual asserts maps exp and val contain the same key/val pairs, using Equal
// method to compare vals
func MapEqual[M interfaces.MapEqualFunc[K, V], K comparable, V interfaces.EqualFunc[V]](t T, exp, val M, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapEqual(exp, val), scripts...)
}

// MapLen asserts map is of size n.
func MapLen[M ~map[K]V, K comparable, V any](t T, n int, m M, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapLen(n, m), scripts...)
}

// MapEmpty asserts map is empty.
func MapEmpty[M ~map[K]V, K comparable, V any](t T, m M, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapEmpty(m), scripts...)
}

// MapNotEmpty asserts map is not empty.
func MapNotEmpty[M ~map[K]V, K comparable, V any](t T, m M, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapNotEmpty(m), scripts...)
}

// MapContainsKeys asserts m contains each key in keys.
func MapContainsKeys[M ~map[K]V, K comparable, V any](t T, m M, keys []K, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapContainsKeys(m, keys), scripts...)
}

// MapNotContainsKeys asserts m does not contain any key in keys.
func MapNotContainsKeys[M ~map[K]V, K comparable, V any](t T, m M, keys []K, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapNotContainsKeys(m, keys), scripts...)
}

// MapContainsValues asserts m contains each val in vals.
func MapContainsValues[M ~map[K]V, K comparable, V any](t T, m M, vals []V, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapContainsValues(m, vals), scripts...)
}

// MapNotContainsValues asserts m does not contain any value in vals.
func MapNotContainsValues[M ~map[K]V, K comparable, V any](t T, m M, vals []V, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapNotContainsValues(m, vals), scripts...)
}

// MapContainsValuesFunc asserts m contains each val in vals using the eq function.
func MapContainsValuesFunc[M ~map[K]V, K comparable, V any](t T, m M, vals []V, eq func(V, V) bool, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapContainsValuesFunc(m, vals, eq), scripts...)
}

// MapNotContainsValuesFunc asserts m does not contain any value in vals using the eq function.
func MapNotContainsValuesFunc[M ~map[K]V, K comparable, V any](t T, m M, vals []V, eq func(V, V) bool, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapNotContainsValuesFunc(m, vals, eq), scripts...)
}

// MapContainsValuesEqual asserts m contains each val in vals using the V.Equal method.
func MapContainsValuesEqual[M ~map[K]V, K comparable, V interfaces.EqualFunc[V]](t T, m M, vals []V, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapContainsValuesEqual(m, vals), scripts...)
}

// MapNotContainsValuesEqual asserts m does not contain any value in vals using the V.Equal method.
func MapNotContainsValuesEqual[M ~map[K]V, K comparable, V interfaces.EqualFunc[V]](t T, m M, vals []V, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.MapNotContainsValuesEqual(m, vals), scripts...)
}

// FileExistsFS asserts file exists on the fs.FS filesystem.
//
// Example,
// FileExistsFS(t, os.DirFS("/etc"), "hosts")
func FileExistsFS(t T, system fs.FS, file string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.FileExistsFS(system, file), scripts...)
}

// FileExists asserts file exists on the OS filesystem.
func FileExists(t T, file string, scripts ...PostScript) {
	t.Helper()
	file = strings.TrimPrefix(file, "/")
	invoke(t, assertions.FileExistsFS(os.DirFS(fsRoot), file), scripts...)
}

// FileNotExistsFS asserts file does not exist on the fs.FS filesystem.
//
// Example,
// FileNotExist(t, os.DirFS("/bin"), "exploit.exe")
func FileNotExistsFS(t T, system fs.FS, file string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.FileNotExistsFS(system, file), scripts...)
}

// FileNotExists asserts file does not exist on the OS filesystem.
func FileNotExists(t T, file string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.FileNotExistsFS(os.DirFS(fsRoot), file), scripts...)
}

// DirExistsFS asserts directory exists on the fs.FS filesystem.
//
// Example,
// DirExistsFS(t, os.DirFS("/usr/local"), "bin")
func DirExistsFS(t T, system fs.FS, directory string, scripts ...PostScript) {
	t.Helper()
	directory = strings.TrimPrefix(directory, "/")
	invoke(t, assertions.DirExistsFS(system, directory), scripts...)
}

// DirExists asserts directory exists on the OS filesystem.
func DirExists(t T, directory string, scripts ...PostScript) {
	t.Helper()
	directory = strings.TrimPrefix(directory, "/")
	invoke(t, assertions.DirExistsFS(os.DirFS(fsRoot), directory), scripts...)
}

// DirNotExistsFS asserts directory does not exist on the fs.FS filesystem.
//
// Example,
// DirNotExistsFS(t, os.DirFS("/tmp"), "scratch")
func DirNotExistsFS(t T, system fs.FS, directory string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.DirNotExistsFS(system, directory), scripts...)
}

// DirNotExists asserts directory does not exist on the OS filesystem.
func DirNotExists(t T, directory string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.DirNotExistsFS(os.DirFS(fsRoot), directory), scripts...)
}

// FileModeFS asserts the file or directory at path on fs.FS has exactly the given permission bits.
//
// Example,
// FileModeFS(t, os.DirFS("/bin"), "find", 0655)
func FileModeFS(t T, system fs.FS, path string, permissions fs.FileMode, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.FileModeFS(system, path, permissions), scripts...)
}

// FileMode asserts the file or directory at path on the OS filesystem has exactly the given permission bits.
func FileMode(t T, path string, permissions fs.FileMode, scripts ...PostScript) {
	t.Helper()
	path = strings.TrimPrefix(path, "/")
	invoke(t, assertions.FileModeFS(os.DirFS(fsRoot), path, permissions), scripts...)
}

// FileContainsFS asserts the file on fs.FS contains content as a substring.
//
// Often os.DirFS is used to interact with the host filesystem.
// Example,
// FileContainsFS(t, os.DirFS("/etc"), "hosts", "localhost")
func FileContainsFS(t T, system fs.FS, file, content string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.FileContainsFS(system, file, content), scripts...)
}

// FileContains asserts the file on the OS filesystem contains content as a substring.
func FileContains(t T, file, content string, scripts ...PostScript) {
	t.Helper()
	file = strings.TrimPrefix(file, "/")
	invoke(t, assertions.FileContainsFS(os.DirFS(fsRoot), file, content), scripts...)
}

// FilePathValid asserts path is a valid file path.
func FilePathValid(t T, path string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.FilePathValid(path), scripts...)
}

// StrEqFold asserts exp and val are equivalent, ignoring case.
func StrEqFold(t T, exp, val string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrEqFold(exp, val), scripts...)
}

// StrNotEqFold asserts exp and val are not equivalent, ignoring case.
func StrNotEqFold(t T, exp, val string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrNotEqFold(exp, val), scripts...)
}

// StrContains asserts s contains substring sub.
func StrContains(t T, s, sub string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrContains(s, sub), scripts...)
}

// StrContainsFold asserts s contains substring sub, ignoring case.
func StrContainsFold(t T, s, sub string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrContainsFold(s, sub), scripts...)
}

// StrNotContains asserts s does not contain substring sub.
func StrNotContains(t T, s, sub string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrNotContains(s, sub), scripts...)
}

// StrNotContainsFold asserts s does not contain substring sub, ignoring case.
func StrNotContainsFold(t T, s, sub string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrNotContainsFold(s, sub), scripts...)
}

// StrContainsAny asserts s contains at least one character in chars.
func StrContainsAny(t T, s, chars string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrContainsAny(s, chars), scripts...)
}

// StrNotContainsAny asserts s does not contain any character in chars.
func StrNotContainsAny(t T, s, chars string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrNotContainsAny(s, chars), scripts...)
}

// StrCount asserts s contains exactly count instances of substring sub.
func StrCount(t T, s, sub string, count int, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrCount(s, sub, count), scripts...)
}

// StrContainsFields asserts that fields is a subset of the result of strings.Fields(s).
func StrContainsFields(t T, s string, fields []string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrContainsFields(s, fields), scripts...)
}

// StrHasPrefix asserts that s starts with prefix.
func StrHasPrefix(t T, prefix, s string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrHasPrefix(prefix, s), scripts...)
}

// StrNotHasPrefix asserts that s does not start with prefix.
func StrNotHasPrefix(t T, prefix, s string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrNotHasPrefix(prefix, s), scripts...)
}

// StrHasSuffix asserts that s ends with suffix.
func StrHasSuffix(t T, suffix, s string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrHasSuffix(suffix, s), scripts...)
}

// StrNotHasSuffix asserts that s does not end with suffix.
func StrNotHasSuffix(t T, suffix, s string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.StrNotHasSuffix(suffix, s), scripts...)
}

// RegexMatch asserts regular expression re matches string s.
func RegexMatch(t T, re *regexp.Regexp, s string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.RegexMatch(re, s), scripts...)
}

// RegexCompiles asserts expr compiles as a valid regular expression.
func RegexCompiles(t T, expr string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.RegexpCompiles(expr), scripts...)
}

// RegexCompilesPOSIX asserts expr compiles as a valid POSIX regular expression.
func RegexCompilesPOSIX(t T, expr string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.RegexpCompilesPOSIX(expr), scripts...)
}

// UUIDv4 asserts id meets the criteria of a v4 UUID.
func UUIDv4(t T, id string, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.UUIDv4(id), scripts...)
}

// Size asserts s.Size() is equal to exp.
func Size(t T, exp int, s interfaces.SizeFunc, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Size(exp, s), scripts...)
}

// Length asserts l.Len() is equal to exp.
func Length(t T, exp int, l interfaces.LengthFunc, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Length(exp, l), scripts...)
}

// Empty asserts e.Empty() is true.
func Empty(t T, e interfaces.EmptyFunc, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Empty(e), scripts...)
}

// NotEmpty asserts e.Empty() is false.
func NotEmpty(t T, e interfaces.EmptyFunc, scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.NotEmpty(e), scripts...)
}

// Contains asserts container.ContainsFunc(element) is true.
func Contains[C any](t T, element C, container interfaces.ContainsFunc[C], scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.Contains(element, container), scripts...)
}

// NotContains asserts container.ContainsFunc(element) is false.
func NotContains[C any](t T, element C, container interfaces.ContainsFunc[C], scripts ...PostScript) {
	t.Helper()
	invoke(t, assertions.NotContains(element, container), scripts...)
}

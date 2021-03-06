package assertions

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/shoenig/test/interfaces"
	"github.com/shoenig/test/internal/constraints"
)

const depth = 4

func Caller() string {
	_, file, line, ok := runtime.Caller(depth)
	if ok {
		file = filepath.Base(file)
		return fmt.Sprintf("%s:%d: ", file, line)
	}
	return "[???]"
}

// diff creates a diff of a and b using cmp.Diff if possible, falling back to printing
// the Go string values of both types (e.g. contains unexported fields).
func diff[A, B any](a A, b B) (s string) {
	defer func() {
		if r := recover(); r != nil {
			s = fmt.Sprintf("↪ comparison ↷\na: %#v\nb: %#v\n", a, b)
		}
	}()
	s = "↪ differential ↷\n" + cmp.Diff(a, b)
	return
}

// equal compares a and b using cmp.Equal if possible, falling back to reflect.DeepEqual
// (e.g. contains unexported fields).
func equal[A, B any](a A, b B) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			result = reflect.DeepEqual(a, b)
		}
	}()
	result = cmp.Equal(a, b)
	return
}

func contains[C comparable](slice []C, item C) bool {
	found := false
	for i := 0; i < len(slice); i++ {
		if slice[i] == item {
			found = true
			break
		}
	}
	return found
}

func containsFunc[A any](slice []A, item A, eq func(a, b A) bool) bool {
	found := false
	for i := 0; i < len(slice); i++ {
		if eq(slice[i], item) {
			found = true
			break
		}
	}
	return found
}

func isNil(a any) bool {
	// comparable check only works for simple types
	if a == nil {
		return true
	}

	// check for non-nil nil types
	value := reflect.ValueOf(a)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func Nil(a any) (s string) {
	if !isNil(a) {
		s = "expected to be nil; is not nil"
	}
	return
}

func NotNil(a any) (s string) {
	if isNil(a) {
		s = "expected to not be nil; is nil"
	}
	return
}

func True(condition bool) (s string) {
	if !condition {
		s = "expected condition to be true; is false"
	}
	return
}

func False(condition bool) (s string) {
	if condition {
		s = "expected condition to be false; is true"
	}
	return
}

func Unreachable() (s string) {
	s = "expected not to execute this code path"
	return
}

func Error(err error) (s string) {
	if err == nil {
		s = "expected non-nil error; is nil"
	}
	return
}

func EqError(err error, msg string) (s string) {
	if err == nil {
		s = "expected error; got nil"
		return
	}
	e := err.Error()
	if e != msg {
		s = "expected matching error strings\n"
		s += fmt.Sprintf("↪ msg: %q\n", msg)
		s += fmt.Sprintf("↪ err: %q\n", e)
	}
	return
}

func ErrorIs(err error, target error) (s string) {
	if err == nil {
		s = "expected error; got nil"
		return
	}
	if !errors.Is(err, target) {
		s = "expected errors.Is match\n"
		s += fmt.Sprintf("↪ error: %v\n", err)
		s += fmt.Sprintf("↪ target: %v\n", err)
	}
	return
}

func NoError(err error) (s string) {
	if err != nil {
		s = "expected nil error\n"
		s += fmt.Sprintf("↪ error: %v", err)
	}
	return
}

func Eq[A any](a, b A) (s string) {
	if !equal(a, b) {
		s = "expected equality via cmp.Equal function\n"
		s += diff(a, b)
	}
	return
}

func NotEq[A any](a, b A) (s string) {
	if equal(a, b) {
		s = "expected inequality via cmp.Equal function\n"
	}
	return
}

func EqOp[C comparable](a, b C) (s string) {
	if a != b {
		s = "expected equality via ==\n"
		s += diff(a, b)
	}
	return
}

func EqFunc[A any](a, b A, eq func(a, b A) bool) (s string) {
	if !eq(a, b) {
		s = "expected equality via 'eq' function\n"
		s += diff(a, b)
	}
	return
}

func NotEqOp[C comparable](a, b C) (s string) {
	if a == b {
		s = "expected inequality via !="
	}
	return
}

func NotEqFunc[A any](a, b A, eq func(a, b A) bool) (s string) {
	if eq(a, b) {
		s = "expected inequality via 'eq' function"
	}
	return
}

func EqJSON(a, b string) (s string) {
	var expA, expB any

	if err := json.Unmarshal([]byte(a), &expA); err != nil {
		s = fmt.Sprintf("failed to unmarshal first argument as json: %v", err)
		return
	}

	if err := json.Unmarshal([]byte(b), &expB); err != nil {
		s = fmt.Sprintf("failed to unmarshal second argument as json: %v", err)
		return
	}

	if !reflect.DeepEqual(expA, expB) {
		jsonA, _ := json.Marshal(expA)
		jsonB, _ := json.Marshal(expB)
		s = "expected equality via json marshalling\n"
		s += diff(string(jsonA), string(jsonB))
		return
	}

	return
}

func EqSliceFunc[A any](a, b []A, eq func(a, b A) bool) (s string) {
	lenA, lenB := len(a), len(b)

	if lenA != lenB {
		s = "expected slices of same length\n"
		s += fmt.Sprintf("↪ len(slice a): %d\n", lenA)
		s += fmt.Sprintf("↪ len(slice b): %d\n", lenB)
		s += diff(a, b)
		return
	}

	miss := false
	for i := 0; i < lenA; i++ {
		if !eq(a[i], b[i]) {
			miss = true
			break
		}
	}

	if miss {
		s = "expected slice equality via 'eq' function\n"
		s += diff(a, b)
		return
	}

	return
}

func Equals[E interfaces.EqualsFunc[E]](a, b E) (s string) {
	if !a.Equals(b) {
		s = "expected equality via .Equals method\n"
		s += diff(a, b)
	}
	return
}

func NotEquals[E interfaces.EqualsFunc[E]](a, b E) (s string) {
	if a.Equals(b) {
		s = "expected inequality via .Equals method\n"
		s += diff(a, b)
	}
	return
}

func EqualsSlice[E interfaces.EqualsFunc[E]](a, b []E) (s string) {
	lenA, lenB := len(a), len(b)

	if lenA != lenB {
		s = "expected slices of same length\n"
		s += fmt.Sprintf("↪ len(slice a): %d\n", lenA)
		s += fmt.Sprintf("↪ len(slice b): %d\n", lenB)
		s += diff(a, b)
		return
	}

	for i := 0; i < lenA; i++ {
		if !a[i].Equals(b[i]) {
			s += "expected slice equality via .Equals method\n"
			s += diff(a[i], b[i])
			return
		}
	}
	return
}

func Lesser[L interfaces.LessFunc[L]](a, b L) (s string) {
	if !a.Less(b) {
		s = "expected to be less via .Less method\n"
		s += diff(a, b)
	}
	return
}

func EmptySlice[A any](slice []A) (s string) {
	if len(slice) != 0 {
		s = "expected slice to be empty\n"
		s += fmt.Sprintf("↪ len(slice): %d\n", len(slice))
	}
	return
}

func LenSlice[A any](n int, slice []A) (s string) {
	if l := len(slice); l != n {
		s = "expected slice to be different length\n"
		s += fmt.Sprintf("↪ len(slice): %d, expected: %d\n", l, n)
	}
	return
}

func Contains[A any](slice []A, item A) (s string) {
	if !containsFunc(slice, item, func(a, b A) bool {
		return equal(a, b)
	}) {
		s = "expected slice to contain missing item via cmp.Equal function\n"
		s += fmt.Sprintf("↪ slice is missing %#v\n", item)
	}
	return
}

func ContainsOp[C comparable](slice []C, item C) (s string) {
	if !contains(slice, item) {
		s = "expected slice to contain missing item via == operator\n"
		s += fmt.Sprintf("↪ slice is missing %#v\n", item)
	}
	return
}

func ContainsFunc[A any](slice []A, item A, eq func(a, b A) bool) (s string) {
	if !containsFunc(slice, item, eq) {
		s = "expected slice to contain missing item via 'eq' function\n"
		s += fmt.Sprintf("↪ slice is missing %#v\n", item)
	}
	return
}

func ContainsEquals[E interfaces.EqualsFunc[E]](slice []E, item E) (s string) {
	if !containsFunc(slice, item, E.Equals) {
		s = "expected slice to contain missing item via .Equals method\n"
		s += fmt.Sprintf("↪ slice is missing %#v\n", item)
	}
	return
}

func ContainsString(original, sub string) (s string) {
	if !strings.Contains(original, sub) {
		s = "expected to contain substring\n"
		s += fmt.Sprintf("↪ str: %s\n", original)
		s += fmt.Sprintf("↪ sub: %s\n", sub)
	}
	return
}

func Positive[N interfaces.Number](n N) (s string) {
	if !(n > 0) {
		s = "expected positive value\n"
		s += fmt.Sprintf("↪ n: %v\n", n)
	}
	return
}

func Negative[N interfaces.Number](n N) (s string) {
	if n > 0 {
		s = "expected negative value\n"
		s += fmt.Sprintf("↪ n: %v\n", n)
	}
	return
}

func Zero[N interfaces.Number](n N) (s string) {
	if n != 0 {
		s = "expected zero\n"
		s += fmt.Sprintf("↪ n: %v\n", n)
	}
	return
}

func NonZero[N interfaces.Number](n N) (s string) {
	if n == 0 {
		s = "expected non-zero\n"
		s += fmt.Sprintf("↪ n: %v\n", n)
	}
	return
}

func Less[O constraints.Ordered](a, b O) (s string) {
	if !(a < b) {
		s = fmt.Sprintf("expected %v < %v", a, b)
	}
	return
}

func LessEq[O constraints.Ordered](a, b O) (s string) {
	if !(a <= b) {
		s = fmt.Sprintf("expected %v <= %v", a, b)
	}
	return
}

func Greater[O constraints.Ordered](a, b O) (s string) {
	if !(a > b) {
		s = fmt.Sprintf("expected %v > %v", a, b)
	}
	return
}

func GreaterEq[O constraints.Ordered](a, b O) (s string) {
	if !(a >= b) {
		s = fmt.Sprintf("expected %v >= %v", a, b)
	}
	return
}

func Ascending[O constraints.Ordered](slice []O) (s string) {
	for i := 0; i < len(slice)-1; i++ {
		if slice[i] > slice[i+1] {
			s = fmt.Sprintf("expected slice[%d] <= slice[%d]\n", i, i+1)
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i, slice[i])
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i+1, slice[i+1])
			return
		}
	}
	return
}

func AscendingFunc[A any](slice []A, less func(a, b A) bool) (s string) {
	for i := 0; i < len(slice)-1; i++ {
		if !less(slice[i], slice[i+1]) {
			s = fmt.Sprintf("expected less(slice[%d], slice[%d])\n", i, i+1)
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i, slice[i])
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i+1, slice[i+1])
			return
		}
	}
	return
}

func AscendingLess[L interfaces.LessFunc[L]](slice []L) (s string) {
	for i := 0; i < len(slice)-1; i++ {
		if !slice[i].Less(slice[i+1]) {
			s = fmt.Sprintf("expected slice[%d].Less(slice[%d])\n", i, i+1)
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i, slice[i])
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i+1, slice[i+1])
			return
		}
	}
	return
}

func Descending[O constraints.Ordered](slice []O) (s string) {
	for i := 0; i < len(slice)-1; i++ {
		if slice[i] < slice[i+1] {
			s = fmt.Sprintf("expected slice[%d] >= slice[%d]\n", i, i+1)
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i, slice[i])
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i+1, slice[i+1])
			return
		}
	}
	return
}

func DescendingFunc[A any](slice []A, less func(a, b A) bool) (s string) {
	for i := 0; i < len(slice)-1; i++ {
		if !less(slice[i+1], slice[i]) {
			s = fmt.Sprintf("expected less(slice[%d], slice[%d])\n", i+1, i)
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i, slice[i])
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i+1, slice[i+1])
			return
		}
	}
	return
}

func DescendingLess[L interfaces.LessFunc[L]](slice []L) (s string) {
	for i := 0; i < len(slice)-1; i++ {
		if !(slice[i+1].Less(slice[i])) {
			s = fmt.Sprintf("expected slice[%d].Less(slice[%d])\n", i+1, i)
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i, slice[i])
			s += fmt.Sprintf("↪ slice[%d]: %v\n", i+1, slice[i+1])
			return
		}
	}
	return
}

func InDelta[N interfaces.Number](a, b, delta N) (s string) {
	var zero N

	if !interfaces.Numeric(delta) {
		s = fmt.Sprintf("delta must be numeric; got %v", delta)
		return
	}

	if delta <= zero {
		s = fmt.Sprintf("delta must be positive; got %v", delta)
		return
	}

	if !interfaces.Numeric(a) {
		s = fmt.Sprintf("first argument must be numeric; got %v", a)
		return
	}

	if !interfaces.Numeric(b) {
		s = fmt.Sprintf("second argument must be numeric; got %v", b)
		return
	}

	difference := a - b
	if difference < -delta || difference > delta {
		s = fmt.Sprintf("%v and %v not within %v", a, b, delta)
		return
	}

	return
}

func InDeltaSlice[N interfaces.Number](a, b []N, delta N) (s string) {
	if len(a) != len(b) {
		s = "expected slices of same length\n"
		s += fmt.Sprintf("↪ len(slice a): %d\n", len(a))
		s += fmt.Sprintf("↪ len(slice b): %d\n", len(b))
		return
	}

	for i := 0; i < len(a); i++ {
		if s = InDelta(a[i], b[i], delta); s != "" {
			return
		}
	}
	return
}

func MapEq[M1, M2 interfaces.Map[K, V], K comparable, V any](a M1, b M2) (s string) {
	lenA, lenB := len(a), len(b)

	if lenA != lenB {
		s = "expected maps of same length\n"
		s += fmt.Sprintf("↪ len(map a): %d\n", lenA)
		s += fmt.Sprintf("↪ len(map b): %d\n", lenB)
		return
	}

	for key, valueA := range a {
		valueB, exists := b[key]
		if !exists {
			s = "expected maps of same keys\n"
			s += diff(a, b)
			return
		}

		if !cmp.Equal(valueA, valueB) {
			s = "expected maps of same values via cmp.Diff function\n"
			s += diff(a, b)
			return
		}
	}
	return
}

func MapEqFunc[M1, M2 interfaces.Map[K, V], K comparable, V any](a M1, b M2, eq func(V, V) bool) (s string) {
	lenA, lenB := len(a), len(b)

	if lenA != lenB {
		s = "expected maps of same length\n"
		s += fmt.Sprintf("↪ len(map a): %d\n", lenA)
		s += fmt.Sprintf("↪ len(map b): %d\n", lenB)
		return
	}

	for key, valueA := range a {
		valueB, exists := b[key]
		if !exists {
			s = "expected maps of same keys\n"
			s += diff(a, b)
			return
		}

		if !eq(valueA, valueB) {
			s = "expected maps of same values via 'eq' function\n"
			s += diff(a, b)
			return
		}
	}
	return
}

func MapEquals[M interfaces.MapEqualsFunc[K, V], K comparable, V interfaces.EqualsFunc[V]](a, b M) (s string) {
	lenA, lenB := len(a), len(b)

	if lenA != lenB {
		s = "expected maps of same length\n"
		s += fmt.Sprintf("↪ len(map a): %d\n", lenA)
		s += fmt.Sprintf("↪ len(map b): %d\n", lenB)
		return
	}

	for key, valueA := range a {
		valueB, exists := b[key]
		if !exists {
			s = "expected maps of same keys\n"
			s += diff(a, b)
			return
		}

		if !(valueB).Equals(valueA) {
			s = "expected maps of same values via .Equals method\n"
			s += diff(a, b)
			return
		}
	}

	return
}

func MapLen[M ~map[K]V, K comparable, V any](n int, m M) (s string) {
	if l := len(m); l != n {
		s = "expected map to be different length\n"
		s += fmt.Sprintf("↪ len(map): %d, expected: %d\n", l, n)
	}
	return
}

func MapEmpty[M ~map[K]V, K comparable, V any](m M) (s string) {
	if l := len(m); l > 0 {
		s = "expected map to be empty\n"
		s += fmt.Sprintf("↪ len(map): %d\n", l)
	}
	return
}

func MapContainsKeys[M ~map[K]V, K comparable, V any](m M, keys []K) (s string) {
	var missing []K
	for _, key := range keys {
		if _, exists := m[key]; !exists {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		s = "expected map to contain keys\n"
		for _, key := range missing {
			s += fmt.Sprintf("↪ key: %v\n", key)
		}
	}
	return
}

func mapContains[M ~map[K]V, K comparable, V any](m M, values []V, eq func(V, V) bool) (s string) {
	var missing []V
	for _, wanted := range values {
		found := false
		for _, v := range m {
			if equal[V](wanted, v) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, wanted)
		}
	}

	if len(missing) > 0 {
		s = "expected map to contain values\n"
		for _, value := range missing {
			s += fmt.Sprintf("↪ value: %v\n", value)
		}
	}
	return
}

func MapContainsValues[M ~map[K]V, K comparable, V any](m M, values []V) (s string) {
	return mapContains[M, K, V](m, values, func(a, b V) bool {
		return equal(a, b)
	})
}

func MapContainsValuesFunc[M ~map[K]V, K comparable, V any](m M, values []V, eq func(V, V) bool) (s string) {
	return mapContains[M, K, V](m, values, eq)
}

func MapContainsValuesEquals[M ~map[K]V, K comparable, V interfaces.EqualsFunc[V]](m M, values []V) (s string) {
	return mapContains[M, K, V](m, values, func(a, b V) bool {
		return a.Equals(b)
	})
}

func FileExists(system fs.FS, file string) (s string) {
	info, err := fs.Stat(system, file)
	if errors.Is(err, fs.ErrNotExist) {
		s = "expected file to exist\n"
		s += fmt.Sprintf("↪ name: %s\n", file)
		s += fmt.Sprintf("↪ error: %s\n", err)
		return
	}
	// other errors - file probably exists but cannot be read
	if info.IsDir() {
		s = "expected file but is a directory\n"
		s += fmt.Sprintf("↪ name: %s\n", file)
		return
	}
	return
}

func FileNotExists(system fs.FS, file string) (s string) {
	_, err := fs.Stat(system, file)
	if !errors.Is(err, fs.ErrNotExist) {
		s = "expected file to not exist\n"
		s += fmt.Sprintf("↪ name: %s\n", file)
		return
	}
	return
}

func DirExists(system fs.FS, directory string) (s string) {
	info, err := fs.Stat(system, directory)
	if os.IsNotExist(err) {
		s = "expected directory to exist\n"
		s += fmt.Sprintf("↪ name: %s\n", directory)
		s += fmt.Sprintf("↪ error: %s\n", err)
		return
	}
	// other errors - directory probably exists but cannot be read
	if !info.IsDir() {
		s = "expected directory but is a file\n"
		s += fmt.Sprintf("↪ name: %s\n", directory)
		return
	}
	return
}

func DirNotExists(system fs.FS, directory string) (s string) {
	_, err := fs.Stat(system, directory)
	if !errors.Is(err, fs.ErrNotExist) {
		s = "expected directory to not exist\n"
		s += fmt.Sprintf("↪ name: %s\n", directory)
		return
	}
	return
}

func FileMode(system fs.FS, path string, permissions fs.FileMode) (s string) {
	info, err := fs.Stat(system, path)
	if err != nil {
		s = "expected to stat path\n"
		s += fmt.Sprintf("↪ name: %s\n", path)
		s += fmt.Sprintf("↪ error: %s\n", err)
		return
	}

	mode := info.Mode()
	if permissions != mode {
		s = "expected different file permissions\n"
		s += fmt.Sprintf("↪ name: %s\n", path)
		s += fmt.Sprintf("↪ exp: %s\n", permissions)
		s += fmt.Sprintf("↪ got: %s\n", mode)
	}
	return
}

func FileContains(system fs.FS, file, content string) (s string) {
	b, err := fs.ReadFile(system, file)
	if err != nil {
		s = "expected to read file\n"
		s += fmt.Sprintf("↪ name: %s\n", file)
		s += fmt.Sprintf("↪ error: %s\n", err)
		return
	}
	actual := string(b)
	if !strings.Contains(string(b), content) {
		s = "expected file contents\n"
		s += fmt.Sprintf("↪ name: %s\n", file)
		s += fmt.Sprintf("↪ wanted: %s\n", content)
		s += fmt.Sprintf("↪ actual: %s\n", actual)
		return
	}
	return
}

func FilePathValid(path string) (s string) {
	if !fs.ValidPath(path) {
		s = "expected valid file path\n"
	}
	return
}

func RegexMatch(re *regexp.Regexp, target string) (s string) {
	if !re.MatchString(target) {
		s = "expected regexp match\n"
		s += fmt.Sprintf("↪ s: %s\n", target)
		s += fmt.Sprintf("↪ re: %s\n", re)
	}
	return
}

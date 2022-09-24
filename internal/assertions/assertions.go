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
			s = fmt.Sprintf("↪ Assertion | comparison ↷\na: %#v\nb: %#v\n", a, b)
		}
	}()
	s = "↪ Assertion | differential ↷\n" + cmp.Diff(a, b)
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

func Eq[A any](expectation, value A) (s string) {
	if !equal(expectation, value) {
		s = "expected equality via cmp.Equal function\n"
		s += diff(expectation, value)
	}
	return
}

func NotEq[A any](expectation, value A) (s string) {
	if equal(expectation, value) {
		s = "expected inequality via cmp.Equal function\n"
	}
	return
}

func EqOp[C comparable](expectation, value C) (s string) {
	if expectation != value {
		s = "expected equality via ==\n"
		s += diff(expectation, value)
	}
	return
}

func EqFunc[A any](expectation, value A, eq func(a, b A) bool) (s string) {
	if !eq(expectation, value) {
		s = "expected equality via 'eq' function\n"
		s += diff(expectation, value)
	}
	return
}

func NotEqOp[C comparable](expectation, value C) (s string) {
	if expectation == value {
		s = "expected inequality via !="
	}
	return
}

func NotEqFunc[A any](expectation, value A, eq func(a, b A) bool) (s string) {
	if eq(expectation, value) {
		s = "expected inequality via 'eq' function"
	}
	return
}

func EqJSON(expectation, value string) (s string) {
	var expA, expB any

	if err := json.Unmarshal([]byte(expectation), &expA); err != nil {
		s = fmt.Sprintf("failed to unmarshal first argument as json: %v", err)
		return
	}

	if err := json.Unmarshal([]byte(value), &expB); err != nil {
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

func EqSliceFunc[A any](expectation, value []A, eq func(a, b A) bool) (s string) {
	lenA, lenB := len(expectation), len(value)

	if lenA != lenB {
		s = "expected slices of same length\n"
		s += fmt.Sprintf("↪ len(exp): %d\n", lenA)
		s += fmt.Sprintf("↪ len(val): %d\n", lenB)
		s += diff(expectation, value)
		return
	}

	miss := false
	for i := 0; i < lenA; i++ {
		if !eq(expectation[i], value[i]) {
			miss = true
			break
		}
	}

	if miss {
		s = "expected slice equality via 'eq' function\n"
		s += diff(expectation, value)
		return
	}

	return
}

func Equal[E interfaces.EqualFunc[E]](expectation, value E) (s string) {
	if !value.Equal(expectation) {
		s = "expected equality via .Equal method\n"
		s += diff(expectation, value)
	}
	return
}

func NotEqual[E interfaces.EqualFunc[E]](expectation, value E) (s string) {
	if value.Equal(expectation) {
		s = "expected inequality via .Equal method\n"
		s += diff(expectation, value)
	}
	return
}

func SliceEqual[E interfaces.EqualFunc[E]](expectation, value []E) (s string) {
	lenA, lenB := len(expectation), len(value)

	if lenA != lenB {
		s = "expected slices of same length\n"
		s += fmt.Sprintf("↪ len(exp): %d\n", lenA)
		s += fmt.Sprintf("↪ len(val): %d\n", lenB)
		s += diff(expectation, value)
		return
	}

	for i := 0; i < lenA; i++ {
		if !expectation[i].Equal(value[i]) {
			s += "expected slice equality via .Equal method\n"
			s += diff(expectation[i], value[i])
			return
		}
	}
	return
}

func Lesser[L interfaces.LessFunc[L]](expectation, value L) (s string) {
	if !value.Less(expectation) {
		s = "expected value to be less via .Less method\n"
		s += diff(expectation, value)
	}
	return
}

func SliceEmpty[A any](slice []A) (s string) {
	if len(slice) != 0 {
		s = "expected slice to be empty\n"
		s += fmt.Sprintf("↪ len(slice): %d\n", len(slice))
	}
	return
}

func SliceNotEmpty[A any](slice []A) (s string) {
	if len(slice) == 0 {
		s = "expected slice to not be empty\n"
		s += fmt.Sprintf("↪ len(slice): %d\n", len(slice))
	}
	return
}

func SliceLen[A any](n int, slice []A) (s string) {
	if l := len(slice); l != n {
		s = "expected slice to be different length\n"
		s += fmt.Sprintf("↪ len(slice): %d, expected: %d\n", l, n)
	}
	return
}

func SliceContainsOp[C comparable](slice []C, item C) (s string) {
	if !contains(slice, item) {
		s = "expected slice to contain missing item via == operator\n"
		s += fmt.Sprintf("↪ slice is missing %#v\n", item)
	}
	return
}

func SliceContainsFunc[A any](slice []A, item A, eq func(a, b A) bool) (s string) {
	if !containsFunc(slice, item, eq) {
		s = "expected slice to contain missing item via 'eq' function\n"
		s += fmt.Sprintf("↪ slice is missing %#v\n", item)
	}
	return
}

func SliceContainsEqual[E interfaces.EqualFunc[E]](slice []E, item E) (s string) {
	if !containsFunc(slice, item, E.Equal) {
		s = "expected slice to contain missing item via .Equal method\n"
		s += fmt.Sprintf("↪ slice is missing %#v\n", item)
	}
	return
}

func SliceContains[A any](slice []A, item A) (s string) {
	for _, i := range slice {
		if cmp.Equal(i, item) {
			return
		}
	}
	s = "expected slice to contain missing item via cmp.Equal method\n"
	s += fmt.Sprintf("↪ slice is missing %#v\n", item)
	return
}

func SliceNotContains[A any](slice []A, item A) (s string) {
	for _, i := range slice {
		if cmp.Equal(i, item) {
			s = "expected slice to not contain item but it does\n"
			s += fmt.Sprintf("↪ unwanted item %#v\n", item)
			return
		}
	}
	return
}

func SliceContainsAll[A any](slice, items []A) (s string) {
OUTER:
	for _, target := range items {
		var item A
		for _, item = range slice {
			if cmp.Equal(target, item) {
				continue OUTER
			}
		}
		s = "expected slice to contain missing item\n"
		s += fmt.Sprintf("↪ slice is missing %#v\n", item)
		return
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

func One[N interfaces.Number](n N) (s string) {
	if n != 1 {
		s = "expected one\n"
		s += fmt.Sprintf("↪ n: %v\n", n)
	}
	return
}

func Less[O constraints.Ordered](expectation, value O) (s string) {
	if !(value < expectation) {
		s = fmt.Sprintf("expected %v < %v", value, expectation)
	}
	return
}

func LessEq[O constraints.Ordered](expectation, value O) (s string) {
	if !(value <= expectation) {
		s = fmt.Sprintf("expected %v ≤ %v", value, expectation)
	}
	return
}

func Greater[O constraints.Ordered](expectation, value O) (s string) {
	if !(value > expectation) {
		s = fmt.Sprintf("expected %v > %v", value, expectation)
	}
	return
}

func GreaterEq[O constraints.Ordered](expectation, value O) (s string) {
	if !(value >= expectation) {
		s = fmt.Sprintf("expected %v ≥ %v", value, expectation)
	}
	return
}

func Between[O constraints.Ordered](lower, value, upper O) (s string) {
	if value < lower || value > upper {
		s = fmt.Sprintf("expected value in range (%v ≤ value ≤ %v)\n", lower, upper)
		s += fmt.Sprintf("↪ value: %v\n", value)
		return
	}
	return
}

func BetweenExclusive[O constraints.Ordered](lower, value, upper O) (s string) {
	if value <= lower || value >= upper {
		s = fmt.Sprintf("expected value in range (%v < value < %v)\n", lower, upper)
		s += fmt.Sprintf("↪ value: %v\n", value)
		return
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

func MapEq[M1, M2 interfaces.Map[K, V], K comparable, V any](expectation M1, value M2) (s string) {
	lenA, lenB := len(expectation), len(value)

	if lenA != lenB {
		s = "expected maps of same length\n"
		s += fmt.Sprintf("↪ len(exp): %d\n", lenA)
		s += fmt.Sprintf("↪ len(val): %d\n", lenB)
		return
	}

	for key, valueA := range expectation {
		valueB, exists := value[key]
		if !exists {
			s = "expected maps of same keys\n"
			s += diff(expectation, value)
			return
		}

		if !cmp.Equal(valueA, valueB) {
			s = "expected maps of same values via cmp.Diff function\n"
			s += diff(expectation, value)
			return
		}
	}
	return
}

func MapEqFunc[M1, M2 interfaces.Map[K, V], K comparable, V any](expectation M1, value M2, eq func(V, V) bool) (s string) {
	lenA, lenB := len(expectation), len(value)

	if lenA != lenB {
		s = "expected maps of same length\n"
		s += fmt.Sprintf("↪ len(exp): %d\n", lenA)
		s += fmt.Sprintf("↪ len(val): %d\n", lenB)
		return
	}

	for key, valueA := range expectation {
		valueB, exists := value[key]
		if !exists {
			s = "expected maps of same keys\n"
			s += diff(expectation, value)
			return
		}

		if !eq(valueA, valueB) {
			s = "expected maps of same values via 'eq' function\n"
			s += diff(expectation, value)
			return
		}
	}
	return
}

func MapEqual[M interfaces.MapEqualFunc[K, V], K comparable, V interfaces.EqualFunc[V]](expectation, value M) (s string) {
	lenA, lenB := len(expectation), len(value)

	if lenA != lenB {
		s = "expected maps of same length\n"
		s += fmt.Sprintf("↪ len(exp): %d\n", lenA)
		s += fmt.Sprintf("↪ len(val): %d\n", lenB)
		return
	}

	for key, valueA := range expectation {
		valueB, exists := value[key]
		if !exists {
			s = "expected maps of same keys\n"
			s += diff(expectation, value)
			return
		}

		if !(valueB).Equal(valueA) {
			s = "expected maps of same values via .Equal method\n"
			s += diff(expectation, value)
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

func MapNotEmpty[M ~map[K]V, K comparable, V any](m M) (s string) {
	if l := len(m); l == 0 {
		s = "expected map to not be empty\n"
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

func MapContainsValuesEqual[M ~map[K]V, K comparable, V interfaces.EqualFunc[V]](m M, values []V) (s string) {
	return mapContains[M, K, V](m, values, func(a, b V) bool {
		return a.Equal(b)
	})
}

func FileExistsFS(system fs.FS, file string) (s string) {
	info, err := fs.Stat(system, file)
	if errors.Is(err, fs.ErrNotExist) {
		s = "expected file to exist\n"
		s += fmt.Sprintf("↪  name: %s\n", file)
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

func FileNotExistsFS(system fs.FS, file string) (s string) {
	_, err := fs.Stat(system, file)
	if !errors.Is(err, fs.ErrNotExist) {
		s = "expected file to not exist\n"
		s += fmt.Sprintf("↪ name: %s\n", file)
		return
	}
	return
}

func DirExistsFS(system fs.FS, directory string) (s string) {
	info, err := fs.Stat(system, directory)
	if os.IsNotExist(err) {
		s = "expected directory to exist\n"
		s += fmt.Sprintf("↪  name: %s\n", directory)
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

func DirNotExistsFS(system fs.FS, directory string) (s string) {
	_, err := fs.Stat(system, directory)
	if !errors.Is(err, fs.ErrNotExist) {
		s = "expected directory to not exist\n"
		s += fmt.Sprintf("↪ name: %s\n", directory)
		return
	}
	return
}

func FileModeFS(system fs.FS, path string, permissions fs.FileMode) (s string) {
	info, err := fs.Stat(system, path)
	if err != nil {
		s = "expected to stat path\n"
		s += fmt.Sprintf("↪  name: %s\n", path)
		s += fmt.Sprintf("↪ error: %s\n", err)
		return
	}

	mode := info.Mode()
	if permissions != mode {
		s = "expected different file permissions\n"
		s += fmt.Sprintf("↪ name: %s\n", path)
		s += fmt.Sprintf("↪  exp: %s\n", permissions)
		s += fmt.Sprintf("↪  got: %s\n", mode)
	}
	return
}

func FileContainsFS(system fs.FS, file, content string) (s string) {
	b, err := fs.ReadFile(system, file)
	if err != nil {
		s = "expected to read file\n"
		s += fmt.Sprintf("↪  name: %s\n", file)
		s += fmt.Sprintf("↪ error: %s\n", err)
		return
	}
	actual := string(b)
	if !strings.Contains(string(b), content) {
		s = "expected file contents\n"
		s += fmt.Sprintf("↪   name: %s\n", file)
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

func StrEqFold(expectation, value string) (s string) {
	if !strings.EqualFold(expectation, value) {
		s = "expected strings to be equal ignoring case\n"
		s += fmt.Sprintf("↪ exp: %s\n", expectation)
		s += fmt.Sprintf("↪ val: %s\n", value)
	}
	return
}

func StrNotEqFold(expectation, value string) (s string) {
	if strings.EqualFold(expectation, value) {
		s = "expected strings to not be equal ignoring case; but they are\n"
		s += fmt.Sprintf("↪ exp: %s\n", expectation)
		s += fmt.Sprintf("↪ val: %s\n", value)
	}
	return
}

func StrContains(str, sub string) (s string) {
	if !strings.Contains(str, sub) {
		s = "expected string to contain substring; it does not\n"
		s += fmt.Sprintf("↪ substring: %s\n", sub)
		s += fmt.Sprintf("↪    string: %s\n", str)
	}
	return
}

func StrContainsFold(str, sub string) (s string) {
	upperS := strings.ToUpper(str)
	upperSub := strings.ToUpper(sub)
	return StrContains(upperS, upperSub)
}

func StrNotContains(str, sub string) (s string) {
	if strings.Contains(str, sub) {
		s = "expected string to not contain substring; but it does\n"
		s += fmt.Sprintf("↪ substring: %s\n", sub)
		s += fmt.Sprintf("↪    string: %s\n", str)
	}
	return
}

func StrNotContainsFold(str, sub string) (s string) {
	upperS := strings.ToUpper(str)
	upperSub := strings.ToUpper(sub)
	return StrNotContains(upperS, upperSub)
}

func StrContainsAny(str, chars string) (s string) {
	if !strings.ContainsAny(str, chars) {
		s = "expected string to contain one or more code points\n"
		s += fmt.Sprintf("↪ code-points: %s\n", chars)
		s += fmt.Sprintf("↪      string: %s\n", str)
	}
	return
}

func StrNotContainsAny(str, chars string) (s string) {
	if strings.ContainsAny(str, chars) {
		s = "expected string to not contain code points; but it does\n"
		s += fmt.Sprintf("↪ code-points: %s\n", chars)
		s += fmt.Sprintf("↪      string: %s\n", str)
	}
	return
}

func StrCount(str, sub string, exp int) (s string) {
	count := strings.Count(str, sub)
	if count != exp {
		s = fmt.Sprintf("expected string to contain %d non-overlapping cases of substring\n", exp)
		s += fmt.Sprintf("↪ count: %d\n", count)
	}
	return
}

func StrContainsFields(str string, fields []string) (s string) {
	set := make(map[string]struct{}, len(fields))
	for _, field := range strings.Fields(str) {
		set[field] = struct{}{}
	}
	var missing []string
	for _, field := range fields {
		if _, exists := set[field]; !exists {
			missing = append(missing, field)
		}
	}
	if len(missing) > 0 {
		s = fmt.Sprintf("expected fields of string to contain subset of values\n")
		s += fmt.Sprintf("↪ missing: %s\n", strings.Join(missing, ", "))
	}
	return
}

func StrHasPrefix(str, prefix string) (s string) {
	if !strings.HasPrefix(str, prefix) {
		s = "expected string to have prefix\n"
		s += fmt.Sprintf("↪ prefix: %s\n", prefix)
		s += fmt.Sprintf("↪ string: %s\n", str)
	}
	return
}

func StrNotHasPrefix(str, prefix string) (s string) {
	if strings.HasPrefix(str, prefix) {
		s = "expected string to not have prefix; but it does\n"
		s += fmt.Sprintf("↪ prefix: %s\n", prefix)
		s += fmt.Sprintf("↪ string: %s\n", str)
	}
	return
}

func StrHasSuffix(str, suffix string) (s string) {
	if !strings.HasSuffix(str, suffix) {
		s = "expected string to have suffix\n"
		s += fmt.Sprintf("↪ suffix: %s\n", suffix)
		s += fmt.Sprintf("↪ string: %s\n", str)
	}
	return
}

func StrNotHasSuffix(str, suffix string) (s string) {
	if strings.HasSuffix(str, suffix) {
		s = "expected string to not have suffix; but it does\n"
		s += fmt.Sprintf("↪ suffix: %s\n", suffix)
		s += fmt.Sprintf("↪ string: %s\n", str)
	}
	return
}

func RegexMatch(re *regexp.Regexp, target string) (s string) {
	if !re.MatchString(target) {
		s = "expected regexp match\n"
		s += fmt.Sprintf("↪  regex: %s\n", re)
		s += fmt.Sprintf("↪ string: %s\n", target)
	}
	return
}

func RegexpCompiles(expr string) (s string) {
	if _, err := regexp.Compile(expr); err != nil {
		s = "expected regular expression to compile\n"
		s += fmt.Sprintf("↪ regex: %s\n", expr)
		s += fmt.Sprintf("↪ error: %v\n", err)
	}
	return
}

func RegexpCompilesPOSIX(expr string) (s string) {
	if _, err := regexp.CompilePOSIX(expr); err != nil {
		s = "expected regular expression to compile (posix)\n"
		s += fmt.Sprintf("↪ regex: %s\n", expr)
		s += fmt.Sprintf("↪ error: %v\n", err)
	}
	return
}

// a10b173d-1427-432d-8a27-b12eada42feb
var uuid4Re = regexp.MustCompile(`^[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{12}$`)

func UUIDv4(id string) (s string) {
	if !uuid4Re.MatchString(id) {
		s = "expected well-formed v4 UUID\n"
		s += "↪ format: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX\n"
		s += "↪ actual: " + id + "\n"
	}
	return
}

func Length(n int, length interfaces.LengthFunc) (s string) {
	if l := length.Len(); l != n {
		s = "expected different length\n"
		s += fmt.Sprintf("↪ length:   %d\n↪ expected: %d\n", l, n)
	}
	return
}

func Size(n int, size interfaces.SizeFunc) (s string) {
	if l := size.Size(); l != n {
		s = "expected different size\n"
		s += fmt.Sprintf("↪ size:     %d\n↪ expected: %d\n", l, n)
	}
	return
}

func Empty(e interfaces.EmptyFunc) (s string) {
	if !e.Empty() {
		s = "expected to be empty, but was not\n"
	}
	return
}

func NotEmpty(e interfaces.EmptyFunc) (s string) {
	if e.Empty() {
		s = "expected to not be empty, but is\n"
	}
	return
}

func Contains[C any](i C, c interfaces.Contains[C]) (s string) {
	if !c.Contains(i) {
		s = "expected to contain element, but does not\n"
	}
	return
}

func NotContains[C any](i C, c interfaces.Contains[C]) (s string) {
	if c.Contains(i) {
		s = "expected not to contain element, but it does\n"
	}
	return
}

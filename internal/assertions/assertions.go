package assertions

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"

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
			s = fmt.Sprintf("↪ difference!\na: %#v\nb: %#v\n", a, b)
		}
	}()
	s = "↪ difference:\n" + cmp.Diff(a, b)
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

func Nil(a any) (s string) {
	if a != nil {
		s = "expected to be nil; is not nil"
	}
	return
}

func NotNil(a any) (s string) {
	if a == nil {
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

func Error(err error) (s string) {
	if err == nil {
		s = "expected non-nil error; is nil"
	}
	return
}

func EqError(err error, msg string) (s string) {
	e := err.Error()
	if e != msg {
		s = "expected matching error strings\n"
		s += fmt.Sprintf("↪ msg: %q\n", msg)
		s += fmt.Sprintf("↪ err: %q\n", e)
	}
	return
}

func ErrorIs(err error, target error) (s string) {
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

func NotEq[C comparable](a, b C) (s string) {
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

func MapLen[M map[K]V, K comparable, V any](n int, m M) (s string) {
	if l := len(m); l != n {
		s = "expected map to be different length\n"
		s += fmt.Sprintf("↪ len(map): %d, expected: %d\n", l, n)
	}
	return
}

func MapEmpty[M map[K]V, K comparable, V any](m M) (s string) {
	if l := len(m); l > 0 {
		s = "expected map to be empty\n"
		s += fmt.Sprintf("↪ len(map): %d\n", l)
	}
	return
}

package test

import (
	"encoding/json"
	"errors"
	"math"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/shoenig/test/internal/constraints"
)

// Nil asserts a is nil.
func Nil(t T, a any) {
	t.Helper()

	if a != nil {
		t.Fatalf(";; expected to be nil; is not nil")
	}
}

// NotNil asserts a is not nil.
func NotNil(t T, a any) {
	t.Helper()

	if a == nil {
		t.Fatalf(";; expected to not be nil; is nil")
	}
}

// True asserts that condition is true.
func True(t T, condition bool) {
	t.Helper()

	if !condition {
		t.Fatalf(";; expected condition to be true; is false")
	}
}

// False asserts condition is false.
func False(t T, condition bool) {
	t.Helper()

	if condition {
		t.Fatalf(";; expected condition to be false; is true")
	}
}

// Error asserts err is a non-nil error.
func Error(t T, err error) {
	t.Helper()

	if err == nil {
		t.Fatalf(";; expected non-nil error; is nil")
	}
}

// ErrorIs asserts err
func ErrorIs(t T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Logf("error: %v\n", err)
		t.Logf("target: %v\n", target)
		t.Fatalf(";; expected errors.Is match")
	}
}

// NoError asserts err is a nil error.
func NoError(t T, err error) {
	t.Helper()

	if err != nil {
		t.Logf("error: %v", err)
		t.Fatalf(";; expected nil error")
	}
}

// Eq asserts a and b are equal using cmp.Equal.
func Eq[A any](t T, a, b A) {
	t.Helper()

	if !cmp.Equal(a, b) {
		t.Logf(cmp.Diff(a, b))
		t.Fatalf(";; expected equality via cmp.Equal function")
	}
}

// EqCmp asserts a == b.
func EqCmp[C comparable](t T, a, b C) {
	t.Helper()

	if a != b {
		t.Logf(cmp.Diff(a, b))
		t.Fatalf(";; expected equality via ==")
	}
}

// EqFunc asserts a and b are equal using eq.
func EqFunc[A any](t T, a, b A, eq func(a, b A) bool) {
	t.Helper()

	if !eq(a, b) {
		t.Logf(cmp.Diff(a, b))
		t.Fatalf(";; expected equality via 'eq' function")
	}
}

// NotEq asserts a != b.
func NotEq[C comparable](t T, a, b C) {
	t.Helper()

	if a == b {
		t.Fatalf(";; expected inequality via !=")
	}
}

// NotEqFunc asserts a and b are not equal using eq.
func NotEqFunc[A any](t T, a, b A, eq func(a, b A) bool) {
	t.Helper()

	if eq(a, b) {
		t.Fatalf(";; expected inequality via 'eq' function")
	}
}

// EqJSON asserts a and b are equivalent JSON.
func EqJSON(t T, a, b string) {
	t.Helper()

	var expA, expB any

	if err := json.Unmarshal([]byte(a), &expA); err != nil {
		t.Fatalf("failed to unmarshal first argument as json: %v", err)
	}

	if err := json.Unmarshal([]byte(b), &expB); err != nil {
		t.Fatalf("failed to unmarshal second argument as json: %v", err)
	}

	if !reflect.DeepEqual(expA, expB) {
		jsonA, _ := json.Marshal(expA)
		jsonB, _ := json.Marshal(expB)
		t.Logf(cmp.Diff(string(jsonA), string(jsonB)))
		t.Fatalf(";; expected equality via json marshalling")
	}
}

// EqSliceFunc asserts elements of a and b are the same using eq.
func EqSliceFunc[A any](t T, a, b []A, eq func(a, b A) bool) {
	t.Helper()

	lenA, lenB := len(a), len(b)

	if lenA != lenB {
		t.Logf("len(slice a): %d\n", lenA)
		t.Logf("len(slice b): %d\n", lenB)
		t.Logf(cmp.Diff(a, b))
		t.Fatalf(";; expected slices of same length")
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
		t.Logf(cmp.Diff(a, b))
		t.Fatalf(";; expected slice equality via 'eq' function")
		return
	}
}

// Equals asserts a.Equals(b).
func Equals[E EqualsFunc[E]](t T, a, b E) {
	t.Helper()

	if !a.Equals(b) {
		t.Logf(cmp.Diff(a, b))
		t.Fatalf(";; expected equality via .Equals method")
	}
}

// NotEquals asserts !a.Equals(b).
func NotEquals[E EqualsFunc[E]](t T, a, b E) {
	t.Helper()

	if a.Equals(b) {
		t.Fatalf(";; expected inequality via .Equals method")
	}
}

// EqualsSlice asserts a[n].Equals(b[n]) for each element n in slices a and b.
func EqualsSlice[E EqualsFunc[E]](t T, a, b []E) {
	t.Helper()

	lenA, lenB := len(a), len(b)

	if lenA != lenB {
		t.Logf("len(slice a): %d\n", lenA)
		t.Logf("len(slice b): %d\n", lenB)
		t.Logf(cmp.Diff(a, b))
		t.Fatalf(";; expected slices of same length")
		return
	}

	for i := 0; i < lenA; i++ {
		if !a[i].Equals(b[i]) {
			t.Logf(cmp.Diff(a, b))
			t.Fatalf(";; expected slice equality via .Equals method")
		}
	}
}

// Lesser asserts a.Less(b).
func Lesser[L LessFunc[L]](t T, a, b L) {
	t.Helper()

	if !a.Less(b) {
		t.Logf(cmp.Diff(a, b))
		t.Fatalf(";; expected to be less via .Less method")
	}
}

// EmptySlice asserts slice is empty.
func EmptySlice[A any](t T, slice []A) {
	t.Helper()

	if len(slice) != 0 {
		t.Logf("len(slice): %d\n", len(slice))
		t.Fatalf(";; expected slice to be empty")
	}
}

// LenSlice asserts slice is of length n.
func LenSlice[A any](t T, n int, slice []A) {
	t.Helper()

	l := len(slice)
	if l != n {
		t.Logf("len(slice): %d, expected: %d\n", l, n)
		t.Fatalf(";; expected slice to be different length")
	}
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

// Contains asserts item exists in slice using cmp.Equal function.
func Contains[A any](t T, slice []A, item A) {
	t.Helper()

	if !containsFunc(slice, item, func(a, b A) bool {
		return cmp.Equal(a, b)
	}) {
		t.Logf("slice is missing %#v\n", item)
		t.Fatalf(";; expected slice to contain missing item via cmp.Equal function")
	}
}

// ContainsCmp asserts item exists in slice using == operator.
func ContainsCmp[C comparable](t T, slice []C, item C) {
	t.Helper()

	if !contains(slice, item) {
		t.Logf("slice is missing %#v\n", item)
		t.Fatalf(";; expected slice to contain missing item via == operator")
	}
}

// ContainsFunc asserts item exists in slice, using eq to compare elements.
func ContainsFunc[A any](t T, slice []A, item A, eq func(a, b A) bool) {
	t.Helper()

	if !containsFunc(slice, item, eq) {
		t.Logf("slice is missing %#v\n", item)
		t.Fatalf(";; expected slice to contain missing item via 'eq' function")
	}
}

// ContainsEquals asserts item exists in slice, using Equals to compare elements.
func ContainsEquals[E EqualsFunc[E]](t T, slice []E, item E) {
	t.Helper()

	if !containsFunc(slice, item, E.Equals) {
		t.Logf("slice is missing %#v\n", item)
		t.Fatalf(";; expected slice to contain missing item via .Equals method")
	}
}

// Less asserts a < b.
func Less[O constraints.Ordered](t T, a, b O) {
	t.Helper()

	if !(a < b) {
		t.Fatalf(";; expected %v < %v", a, b)
	}
}

// LessEq asserts a <= b.
func LessEq[O constraints.Ordered](t T, a, b O) {
	t.Helper()

	if !(a <= b) {
		t.Fatalf(";; expected %v <= %v", a, b)
	}
}

// Greater asserts a > b.
func Greater[O constraints.Ordered](t T, a, b O) {
	t.Helper()

	if !(a > b) {
		t.Fatalf(";; expected %v > %v", a, b)
	}
}

// GreaterEq asserts a >= b.
func GreaterEq[O constraints.Ordered](t T, a, b O) {
	t.Helper()

	if !(a >= b) {
		t.Fatalf(";; expected %v >= %v", a, b)
	}
}

// Number is float, integer, or complex.
type Number interface {
	constraints.Ordered
	constraints.Float | constraints.Integer | constraints.Complex
}

// Numeric returns false if n is Inf/NaN.
//
// Always returns true for integral values.
func Numeric[N Number](n N) bool {
	check := func(f float64) bool {
		if math.IsNaN(f) {
			return false
		} else if math.IsInf(f, 0) {
			return false
		}
		return true
	}
	return check(float64(n))
}

// InDelta asserts a and b are within delta of each other.
func InDelta[N Number](t T, a, b, delta N) {
	t.Helper()

	var zero N

	if !Numeric(delta) {
		t.Fatalf(";; delta must be numeric; got %v", delta)
	}

	if delta <= zero {
		t.Fatalf(";; delta must be positive; got %v", delta)
	}

	if !Numeric(a) {
		t.Fatalf(";; first argument must be numeric; got %v", a)
	}

	if !Numeric(b) {
		t.Fatalf(";; second argument must be numeric; got %v", b)
	}

	difference := a - b
	if difference < -delta || difference > delta {
		t.Fatalf(";; %v and %v not within %v", a, b, delta)
	}
}

// InDeltaSlice asserts each element a[n] is within delta of b[n].
func InDeltaSlice[N Number](t T, a, b []N, delta N) {
	t.Helper()

	if len(a) != len(b) {
		t.Logf("len(slice a): %d\n", len(a))
		t.Logf("len(slice b): %d\n", len(b))
		t.Fatalf(";; expected slices of same length")
	}

	for i := 0; i < len(a); i++ {
		InDelta(t, a[i], b[i], delta)
	}
}

// MapEq asserts maps a and b contain the same key/value pairs, using
// cmp.Equal function to compare values.
func MapEq[M1, M2 ~map[K]V, K comparable, V any](t T, a M1, b M2) {
	t.Helper()

	lenA, lenB := len(a), len(b)

	if lenA != lenB {
		t.Logf("len(map a): %d\n", lenA)
		t.Logf("len(map b): %d\n", lenB)
		t.Fatalf(";; expected maps of same length")
		return
	}

	for key, valueA := range a {
		valueB, exists := b[key]
		if !exists {
			t.Logf(cmp.Diff(a, b))
			t.Fatalf(";; expected maps of same keys")
			return
		}

		if !cmp.Equal(valueA, valueB) {
			t.Logf(cmp.Diff(a, b))
			t.Fatalf(";; expected maps of same values via cmp.Diff function")
			return
		}
	}
}

// MapEqFunc asserts maps a and b contain the same key/value pairs, using eq to
// compare values.
func MapEqFunc[M1 ~map[K]V1, M2 ~map[K]V2, K comparable, V1, V2 any](t T, a M1, b M2, eq func(V1, V2) bool) {
	t.Helper()

	lenA, lenB := len(a), len(b)

	if lenA != lenB {
		t.Logf("len(map a): %d\n", lenA)
		t.Logf("len(map b): %d\n", lenB)
		t.Fatalf(";; expected maps of same length")
		return
	}

	for key, valueA := range a {
		valueB, exists := b[key]
		if !exists {
			t.Logf(cmp.Diff(a, b))
			t.Fatalf(";; expected maps of same keys")
			return
		}

		if !eq(valueA, valueB) {
			t.Logf(cmp.Diff(a, b))
			t.Fatalf(";; expected maps of same values via 'eq' function")
			return
		}
	}
}

// MapEquals asserts maps a and b contain the same key/value pairs, using Equals
// method to compare values
func MapEquals[M ~map[K]V, K comparable, V EqualsFunc[V]](t T, a, b M) {
	t.Helper()

	lenA, lenB := len(a), len(b)

	if lenA != lenB {
		t.Logf("len(map a): %d\n", lenA)
		t.Logf("len(map b): %d\n", lenB)
		t.Fatalf(";; expected maps of same length")
		return
	}

	for key, valueA := range a {
		valueB, exists := b[key]
		if !exists {
			t.Logf(cmp.Diff(a, b))
			t.Fatalf(";; expected maps of same keys")
			return
		}

		if !valueB.Equals(valueA) {
			t.Logf(cmp.Diff(a, b))
			t.Fatalf(";; expected maps of same values via .Equals method")
			return
		}
	}
}

// MapLen asserts map is of size n.
func MapLen[M ~map[K]V, K comparable, V any](t T, n int, m M) {
	t.Helper()

	l := len(m)
	if l != n {
		t.Logf("len(map): %d, expected: %d\n", l, n)
		t.Fatalf(";; expected map to be different length")
	}
}

// MapEmpty asserts map is empty.
func MapEmpty[M ~map[K]V, K comparable, V any](t T, m M) {
	t.Helper()

	if l := len(m); l > 0 {
		t.Logf("len(map): %d\n")
		t.Fatalf(";; expected map to be empty")
	}
}

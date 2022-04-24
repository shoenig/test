package test

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"
)

type internalTest struct {
	t        *testing.T
	trigger  bool
	helper   bool
	exp, msg string
}

func (it *internalTest) Helper() {
	it.helper = true
}

func (it *internalTest) Fatalf(msg string, args ...any) {
	if !it.trigger {
		it.trigger = true
		it.msg = fmt.Sprintf(msg, args...)
	}
}

func (it *internalTest) assert() {
	if !it.helper {
		it.t.Fatal("should be marked as helper")
	}
	if !it.trigger {
		it.t.Fatalf("condition expected to trigger; did not")
	}
	if it.exp != it.msg {
		it.t.Fatalf("expected message %q ... got %q", it.exp, it.msg)
	}
}

func newCase(t *testing.T, msg string) *internalTest {
	return &internalTest{
		t:       t,
		trigger: false,
		exp:     msg,
	}
}

func TestNil(t *testing.T) {
	tc := newCase(t, `expected to be nil; is not nil`)
	t.Cleanup(tc.assert)

	Nil(tc, 42)
	Nil(tc, "hello")
	Nil(tc, time.UTC)
}

func TestNilf(t *testing.T) {
	tc := newCase(t, `a message: 42`)
	t.Cleanup(tc.assert)

	Nilf(tc, 0, "a message: %d", 42)
}

func TestNotNil(t *testing.T) {
	tc := newCase(t, `expected to not be nil; is nil`)
	t.Cleanup(tc.assert)

	NotNil(tc, nil)
}

func TestNotNilf(t *testing.T) {
	tc := newCase(t, `a message: 42`)
	t.Cleanup(tc.assert)

	NotNilf(tc, nil, "a message: %d", 42)
}

func TestTrue(t *testing.T) {
	tc := newCase(t, `expected condition to be true; is false`)
	t.Cleanup(tc.assert)

	True(tc, false)
}

func TestTruef(t *testing.T) {
	tc := newCase(t, `a message: 42`)
	t.Cleanup(tc.assert)

	Truef(tc, false, "a message: %d", 42)
}

func TestFalse(t *testing.T) {
	tc := newCase(t, `expected condition to be false; is true`)
	t.Cleanup(tc.assert)

	False(tc, true)
}

func TestFalsef(t *testing.T) {
	tc := newCase(t, `a message: 42`)
	t.Cleanup(tc.assert)

	Falsef(tc, true, "a message: %d", 42)
}

func TestError(t *testing.T) {
	tc := newCase(t, `expected non-nil error; is nil`)
	t.Cleanup(tc.assert)

	Error(tc, nil)
}

func TestErrorIs(t *testing.T) {
	tc := newCase(t, `expected foo errors.Is bar`)
	t.Cleanup(tc.assert)

	e1 := errors.New("foo")
	e2 := errors.New("bar")
	ErrorIs(tc, e1, e2)
}

func TestNoError(t *testing.T) {
	tc := newCase(t, `expected nil error, got "hello"`)
	t.Cleanup(tc.assert)

	NoError(tc, errors.New("hello"))
}

func TestEq(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		tc := newCase(t, `expected 42 == 43`)
		t.Cleanup(tc.assert)

		Eq(tc, 42, 43)
	})

	t.Run("string", func(t *testing.T) {
		tc := newCase(t, `expected foo == bar`)
		t.Cleanup(tc.assert)

		Eq(tc, "foo", "bar")
	})

	t.Run("duration", func(t *testing.T) {
		tc := newCase(t, `expected 2s == 3m0s`)
		t.Cleanup(tc.assert)

		a := 2 * time.Second
		b := 3 * time.Minute
		Eq(tc, a, b)
	})
}

func TestEqf(t *testing.T) {
	tc := newCase(t, `a number: 42`)
	t.Cleanup(tc.assert)

	Eqf(tc, 1, 2, "a number: %d", 42)
}

func TestNotEq(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		tc := newCase(t, `expected 42 != 42`)
		t.Cleanup(tc.assert)
		NotEq(tc, 42, 42)
	})

	t.Run("string", func(t *testing.T) {
		tc := newCase(t, `expected foo != foo`)
		t.Cleanup(tc.assert)
		NotEq(tc, "foo", "foo")
	})

	t.Run("duration", func(t *testing.T) {
		tc := newCase(t, `expected 3s != 3s`)
		t.Cleanup(tc.assert)
		NotEq(tc, 3*time.Second, 3*time.Second)
	})
}

func TestNotEqf(t *testing.T) {
	tc := newCase(t, `a number: 42`)
	t.Cleanup(tc.assert)

	NotEqf(tc, 42, 42, "a number: %d", 42)
}

func TestEqJSON(t *testing.T) {
	tc := newCase(t, `json strings are not the same; {"a":1,"b":2} vs. {"a":9,"b":2}`)
	t.Cleanup(tc.assert)

	EqJSON(tc, `{"a":1, "b":2}`, `{"b":2, "a":9}`)
}

// Person implements the Equals and Less functions.
type Person struct {
	ID   int
	Name string
}

func (p *Person) Equals(o *Person) bool {
	return p.ID == o.ID
}

func (p *Person) Less(o *Person) bool {
	return p.ID < o.ID
}

func TestEquals(t *testing.T) {
	tc := newCase(t, `expected to be equal: &{100 Alice}, &{150 Alice}`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 150, Name: "Alice"}

	Equals(tc, a, b)
}

func TestEqualsf(t *testing.T) {
	tc := newCase(t, `a number: 42`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 150, Name: "Alice"}

	Equalsf(tc, a, b, "a number: %d", 42)
}

func TestNotEquals(t *testing.T) {
	tc := newCase(t, `expected to be not equal: &{100 Alice}, &{100 Alice}`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 100, Name: "Alice"}

	NotEquals(tc, a, b)
}

func TestNotEqualsf(t *testing.T) {
	tc := newCase(t, `a number: 42`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 100, Name: "Alice"}

	NotEqualsf(tc, a, b, "a number: %d", 42)
}

func TestLesser(t *testing.T) {
	tc := newCase(t, `expected to be less; &{200 Alice}, &{100 Bob}`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 200, Name: "Alice"}
	b := &Person{ID: 100, Name: "Bob"}

	Lesser(tc, a, b)
}

func TestEmpty(t *testing.T) {
	tc := newCase(t, `expected slice to be empty; is len 2`)
	t.Cleanup(tc.assert)

	Empty(tc, []int{1, 2})
}

func TestLen(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected slice to be length 2; is 3`)
		t.Cleanup(tc.assert)
		Len(tc, 2, []string{"a", "b", "c"})
	})

	t.Run("numbers", func(t *testing.T) {
		tc := newCase(t, `expected slice to be length 3; is 2`)
		t.Cleanup(tc.assert)
		Len(tc, 3, []int{8, 9})
	})
}

func TestLenf(t *testing.T) {
	tc := newCase(t, `a number: 42`)
	t.Cleanup(tc.assert)

	Lenf(tc, 3, []int{1, 2}, "a number: %d", 42)
}

func TestContains(t *testing.T) {
	t.Run("numbers", func(t *testing.T) {
		tc := newCase(t, `expected slice to contain 7 but does not: [3 4 5]`)
		t.Cleanup(tc.assert)
		Contains(tc, []int{3, 4, 5}, 7)
	})

	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected slice to contain bob but does not: [alice carl]`)
		t.Cleanup(tc.assert)
		Contains(tc, []string{"alice", "carl"}, "bob")
	})
}

func TestContainsf(t *testing.T) {
	tc := newCase(t, `a number: 42`)
	t.Cleanup(tc.assert)
	Containsf(tc, []int{1, 2, 3}, 4, "a number: %d", 42)
}

func TestLess(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		tc := newCase(t, `expected 7 < 5`)
		t.Cleanup(tc.assert)
		Less(tc, 7, 5)
	})

	t.Run("floats", func(t *testing.T) {
		tc := newCase(t, `expected 7.7 < 5.5`)
		t.Cleanup(tc.assert)
		Less(tc, 7.7, 5.5)
	})

	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected foo < bar`)
		t.Cleanup(tc.assert)
		Less(tc, "foo", "bar")
	})

	t.Run("equal", func(t *testing.T) {
		tc := newCase(t, `expected 7 < 7`)
		t.Cleanup(tc.assert)
		Less(tc, 7, 7)
	})
}

func TestLessf(t *testing.T) {
	tc := newCase(t, `a number: 42`)
	t.Cleanup(tc.assert)
	Lessf(tc, 7, 5, "a number: %d", 42)
}

func TestLessEq(t *testing.T) {
	tc := newCase(t, `expected 7 <= 5`)
	t.Cleanup(tc.assert)
	LessEq(tc, 7, 5)
}

func TestLessEqf(t *testing.T) {
	tc := newCase(t, `a number: 42`)
	t.Cleanup(tc.assert)
	LessEqf(tc, 7, 5, "a number: %d", 42)
}

func TestGreater(t *testing.T) {
	t.Run("integer", func(t *testing.T) {
		tc := newCase(t, `expected 5 > 7`)
		t.Cleanup(tc.assert)
		Greater(tc, 5, 7)
	})

	t.Run("floats", func(t *testing.T) {
		tc := newCase(t, `expected 5.5 > 7.7`)
		t.Cleanup(tc.assert)
		Greater(tc, 5.5, 7.7)
	})

	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected bar > foo`)
		t.Cleanup(tc.assert)
		Greater(tc, "bar", "foo")
	})

	t.Run("equal", func(t *testing.T) {
		tc := newCase(t, `expected bar > bar`)
		t.Cleanup(tc.assert)
		Greater(tc, "bar", "bar")
	})
}

func TestGreaterf(t *testing.T) {
	tc := newCase(t, `a number: 42`)
	t.Cleanup(tc.assert)
	Greaterf(tc, 5, 7, "a number: %d", 42)
}

func TestGreaterEq(t *testing.T) {
	tc := newCase(t, `expected 5 >= 7`)
	t.Cleanup(tc.assert)
	GreaterEq(tc, 5, 7)
}

func TestGreaterEqf(t *testing.T) {
	tc := newCase(t, `a number: 42`)
	t.Cleanup(tc.assert)
	GreaterEqf(tc, 5, 7, "a number: %d", 42)
}

func TestInDelta(t *testing.T) {
	t.Run("inf delta", func(t *testing.T) {
		tc := newCase(t, `delta must be numeric; got +Inf`)
		t.Cleanup(tc.assert)
		InDelta(tc, 100.0, 101.0, math.Inf(1))
	})

	t.Run("nan delta", func(t *testing.T) {
		tc := newCase(t, `delta must be numeric; got NaN`)
		t.Cleanup(tc.assert)
		InDelta(tc, 100.0, 101.0, math.NaN())
	})

	t.Run("negative delta", func(t *testing.T) {
		tc := newCase(t, `delta must be positive; got -3.5`)
		t.Cleanup(tc.assert)
		InDelta(tc, 100.0, 101.0, -3.5)
	})

	t.Run("inf arg1", func(t *testing.T) {
		tc := newCase(t, `first argument must be numeric; got +Inf`)
		t.Cleanup(tc.assert)
		InDelta(tc, math.Inf(1), 101.0, 1.0)
	})

	t.Run("inf arg2", func(t *testing.T) {
		tc := newCase(t, `second argument must be numeric; got +Inf`)
		t.Cleanup(tc.assert)
		InDelta(tc, 100.0, math.Inf(1), 1.0)
	})

	t.Run("float", func(t *testing.T) {
		tc := newCase(t, `100.1 and 101.5 not within 0.7`)
		t.Cleanup(tc.assert)
		InDelta(tc, 100.1, 101.5, 0.7)
	})

	t.Run("int", func(t *testing.T) {
		tc := newCase(t, `50 and 70 not within 10`)
		t.Cleanup(tc.assert)
		InDelta(tc, 50, 70, 10)
	})
}

func TestInDeltaSlice(t *testing.T) {
	t.Run("different length", func(t *testing.T) {
		tc := newCase(t, `slices not of same length; 2 != 3`)
		t.Cleanup(tc.assert)
		InDeltaSlice(tc, []int{2, 3}, []int{2, 3, 4}, 2)
	})

	t.Run("float", func(t *testing.T) {
		tc := newCase(t, `25 and 42 not within 5`)
		t.Cleanup(tc.assert)
		InDeltaSlice(tc, []int{10, 25, 300}, []int{11, 42, 299}, 5)
	})
}

func TestMapEq(t *testing.T) {
	t.Run("different length", func(t *testing.T) {
		tc := newCase(t, `maps are different size; 1 vs 2; map[a:1] != map[a:1 b:2]`)
		t.Cleanup(tc.assert)
		a := map[string]int{"a": 1}
		b := map[string]int{"a": 1, "b": 2}
		MapEq(tc, a, b)
	})

	t.Run("different keys", func(t *testing.T) {
		tc := newCase(t, `map keys are different; map[1:a 2:b] != map[1:a 3:c]`)
		t.Cleanup(tc.assert)
		a := map[int]string{1: "a", 2: "b"}
		b := map[int]string{1: "a", 3: "c"}
		MapEq(tc, a, b)
	})

	t.Run("different values", func(t *testing.T) {
		tc := newCase(t, `value for key b different; map[a:amp b:bar] != map[a:amp b:foo]`)
		t.Cleanup(tc.assert)
		a := map[string]string{"a": "amp", "b": "bar"}
		b := map[string]string{"a": "amp", "b": "foo"}
		MapEq(tc, a, b)
	})
}

func TestMapEmpty(t *testing.T) {
	tc := newCase(t, `expected map to be empty; is size 2`)
	t.Cleanup(tc.assert)
	m := map[string]int{"a": 1, "b": 2}
	MapEmpty(tc, m)
}

func TestMapLen(t *testing.T) {
	tc := newCase(t, `expected map to be size 2; is 3`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	MapLen(tc, 2, m)
}

func TestMapLenf(t *testing.T) {
	tc := newCase(t, `a number: 42`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	MapLenf(tc, 2, m, "a number: %d", 42)
}

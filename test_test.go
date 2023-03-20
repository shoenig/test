package test

import (
	"errors"
	"math"
	"os"
	"regexp"
	"runtime"
	"testing"
	"time"

	"github.com/shoenig/test/wait"
)

func needsOS(t *testing.T, os string) {
	if os != runtime.GOOS {
		t.Skip("not supported on this OS")
	}
}

func TestNil(t *testing.T) {
	tc := newCase(t, `expected to be nil; is not nil`)
	t.Cleanup(tc.assert)

	Nil(tc, 42)
	Nil(tc, "hello")
	Nil(tc, time.UTC)
	Nil(tc, []string{"foo"})
	Nil(tc, map[string]int{"foo": 1})
}

func TestNil_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	Nil(tc, 42, tc.TestPostScript("nil"))
}

func TestNotNil(t *testing.T) {
	tc := newCase(t, `expected to not be nil; is nil`)
	t.Cleanup(tc.assert)

	var s []string
	var m map[string]int

	NotNil(tc, nil)
	NotNil(tc, s)
	NotNil(tc, m)
}

func TestNotNil_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	NotNil(tc, nil, tc.TestPostScript("not nil"))
}

func TestTrue(t *testing.T) {
	tc := newCase(t, `expected condition to be true; is false`)
	t.Cleanup(tc.assert)

	True(tc, false)
}

func TestTrue_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	True(tc, false, tc.TestPostScript("true"))
}

func TestFalse(t *testing.T) {
	tc := newCase(t, `expected condition to be false; is true`)
	t.Cleanup(tc.assert)

	False(tc, true)
}

func TestFalse_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	False(tc, true, tc.TestPostScript("false"))
}

func TestUnreachable(t *testing.T) {
	tc := newCase(t, `expected not to execute this code path`)
	t.Cleanup(tc.assert)

	Unreachable(tc)
}

func TestUnreachable_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	Unreachable(tc, tc.TestPostScript("unreachable"))
}

func TestError(t *testing.T) {
	tc := newCase(t, `expected non-nil error; is nil`)
	t.Cleanup(tc.assert)

	Error(tc, nil)
}

func TestError_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	Error(tc, nil, tc.TestPostScript("error"))
}

func TestEqError(t *testing.T) {
	tc := newCase(t, `expected matching error strings`)
	t.Cleanup(tc.assert)

	EqError(tc, errors.New("oops"), "blah")
}

func TestEqError_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	EqError(tc, errors.New("oops"), "blah", tc.TestPostScript("eq error"))
}

func TestEqError_nil(t *testing.T) {
	tc := newCase(t, `expected error; got nil`)
	t.Cleanup(tc.assert)

	EqError(tc, nil, "blah")
}

func TestErrorIs(t *testing.T) {
	tc := newCase(t, `expected errors.Is match`)
	t.Cleanup(tc.assert)

	e1 := errors.New("foo")
	e2 := errors.New("bar")
	ErrorIs(tc, e1, e2)
}

func TestErrorIs_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	e1 := errors.New("foo")
	e2 := errors.New("bar")
	ErrorIs(tc, e1, e2, tc.TestPostScript("error is"))
}

func TestErrorIs_nil(t *testing.T) {
	tc := newCase(t, `expected error; got nil`)
	t.Cleanup(tc.assert)

	err := errors.New("oops")
	ErrorIs(tc, nil, err)
}

func TestNoError(t *testing.T) {
	tc := newCase(t, `expected nil error`)
	t.Cleanup(tc.assert)

	NoError(tc, errors.New("hello"))
}

func TestNoError_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	NoError(tc, errors.New("hello"), tc.TestPostScript("no error"))
}

func TestErrorContains(t *testing.T) {
	tc := newCase(t, `expected error to contain substring`)
	t.Cleanup(tc.assert)

	ErrorContains(tc, errors.New("something bad"), "oops")
}

func TestEq(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		tc := newCase(t, `expected equality via cmp.Equal function`)
		t.Cleanup(tc.assert)

		Eq(tc, 42, 43)
	})

	t.Run("string", func(t *testing.T) {
		tc := newCase(t, `expected equality via cmp.Equal function`)
		t.Cleanup(tc.assert)

		Eq(tc, "foo", "bar")
	})

	t.Run("duration", func(t *testing.T) {
		tc := newCase(t, `expected equality via cmp.Equal function`)
		t.Cleanup(tc.assert)

		a := 2 * time.Second
		b := 3 * time.Minute
		Eq(tc, a, b)
	})

	t.Run("person", func(t *testing.T) {
		tc := newCase(t, `expected equality via cmp.Equal function`)
		t.Cleanup(tc.assert)

		p1 := Person{ID: 100, Name: "Alice"}
		p2 := Person{ID: 101, Name: "Bob"}
		Eq(tc, p1, p2)
	})

	t.Run("slice", func(t *testing.T) {
		tc := newCase(t, `expected equality via cmp.Equal function`)
		t.Cleanup(tc.assert)

		a := []int{1, 2, 3, 4}
		b := []int{1, 2, 9, 4}
		Eq(tc, a, b)
	})
}

func TestEq_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	Eq(tc, 1, 2, tc.TestPostScript("eq"))
}

func TestEqOp(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		tc := newCase(t, `expected equality via ==`)
		t.Cleanup(tc.assert)
		EqOp(tc, "foo", "bar")
	})
}

func TestEqOp_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	EqOp(tc, "foo", "bar", tc.TestPostScript("eq op"))
}

func TestEqFunc(t *testing.T) {
	tc := newCase(t, `expected equality via 'eq' function`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 101, Name: "Bob"}

	EqFunc(tc, a, b, func(a, b *Person) bool {
		return a.ID == b.ID && a.Name == b.Name
	})
}

func TestEqFunc_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	EqFunc(tc, "hello", "world", func(a, b string) bool {
		return a == b
	}, tc.TestPostScript("eq func"))
}

func TestNotEq(t *testing.T) {
	tc := newCase(t, `expected inequality via cmp.Equal function`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 100, Name: "Alice"}

	NotEq(tc, a, b)
}

func TestNotEq_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	NotEq(tc, 1, 1, tc.TestPostScript("not eq"))
}

func TestNotEqOp(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		tc := newCase(t, `expected inequality via !=`)
		t.Cleanup(tc.assert)
		NotEqOp(tc, 42, 42)
	})

	t.Run("string", func(t *testing.T) {
		tc := newCase(t, `expected inequality via !=`)
		t.Cleanup(tc.assert)
		NotEqOp(tc, "foo", "foo")
	})

	t.Run("duration", func(t *testing.T) {
		tc := newCase(t, `expected inequality via !=`)
		t.Cleanup(tc.assert)
		NotEqOp(tc, 3*time.Second, 3*time.Second)
	})
}

func TestNotEqOp_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	NotEqOp(tc, 1, 1, tc.TestPostScript("not eq op"))
}

func TestNotEqFunc(t *testing.T) {
	tc := newCase(t, `expected inequality via 'eq' function`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 100, Name: "Alice"}

	NotEqFunc(tc, a, b, func(a, b *Person) bool {
		return a.ID == b.ID && a.Name == b.Name
	})
}

func TestNotEqFunc_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	NotEqFunc(tc, 1, 1, func(a, b int) bool {
		return a == b
	}, tc.TestPostScript("not eq func"))
}

func TestEqJSON(t *testing.T) {
	tc := newCase(t, `expected equality via json marshalling`)
	t.Cleanup(tc.assert)

	EqJSON(tc, `{"a":1, "b":2}`, `{"b":2, "a":9}`)
}

func TestEqJSON_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	EqJSON(tc, `"one"`, `"two"`, tc.TestPostScript("eq json"))
}

func TestValidJSON(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.assert)

	ValidJSON(tc, `{"a":1, "b":}`)
}

func TestValidJSONBytes(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.assert)

	ValidJSONBytes(tc, []byte(`{"a":1, "b":}`))
}

func TestSliceEqFunc(t *testing.T) {
	t.Run("length", func(t *testing.T) {
		tc := newCase(t, `expected slices of same length`)
		t.Cleanup(tc.assert)

		a := []int{1, 2, 3}
		b := []int{1, 2}
		SliceEqFunc(tc, a, b, func(a, b int) bool {
			return false
		})
	})

	t.Run("elements", func(t *testing.T) {
		tc := newCase(t, `expected slice equality via 'eq' function`)
		t.Cleanup(tc.assert)

		a := []*Person{
			{ID: 100, Name: "Alice"},
			{ID: 101, Name: "Bob"},
			{ID: 102, Name: "Carl"},
		}
		b := []*Person{
			{ID: 100, Name: "Alice"},
			{ID: 101, Name: "Bob"},
			{ID: 103, Name: "Dian"},
		}

		SliceEqFunc(tc, a, b, func(a, b *Person) bool {
			return a.ID == b.ID
		})
	})

	t.Run("translate", func(t *testing.T) {
		tc := newCase(t, `expected slice equality via 'eq' function`)
		t.Cleanup(tc.assert)

		values := []*Person{
			{ID: 100, Name: "Alice"},
			{ID: 101, Name: "Bob"},
		}
		exp := []string{"Alice", "Carl"}
		SliceEqFunc(tc, exp, values, func(a *Person, name string) bool {
			return a.Name == name
		})
	})
}

func TestSliceEqFunc_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	a := []int{1, 2, 3}
	b := []int{1, 2}
	SliceEqFunc(tc, a, b, func(a, b int) bool {
		return false
	}, tc.TestPostScript("eq slice func"))
}

// Person implements the Equal and Less functions.
type Person struct {
	ID   int
	Name string
}

func (p *Person) Equal(o *Person) bool {
	return p.ID == o.ID
}

func (p *Person) Less(o *Person) bool {
	return p.ID < o.ID
}

func TestEqual(t *testing.T) {
	tc := newCase(t, `expected equality via .Equal method`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 150, Name: "Alice"}

	Equal(tc, a, b)
}

func TestEqual_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 150, Name: "Alice"}

	Equal(tc, a, b, tc.TestPostScript("equal"))
}

func TestNotEqual(t *testing.T) {
	tc := newCase(t, `expected inequality via .Equal method`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 100, Name: "Alice"}

	NotEqual(tc, a, b)
}

func TestNotEqual_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 100, Name: "Alice"}

	NotEqual(tc, a, b, tc.TestPostScript("not equal"))
}

func TestSliceEqual(t *testing.T) {
	t.Run("length", func(t *testing.T) {
		tc := newCase(t, `expected slices of same length`)
		t.Cleanup(tc.assert)

		a := []*Person{
			{ID: 100, Name: "Alice"},
			{ID: 101, Name: "Bob"},
			{ID: 102, Name: "Carl"},
		}
		b := []*Person{
			{ID: 100, Name: "Alice"},
			{ID: 101, Name: "Bob"},
		}
		SliceEqual(tc, a, b)
	})

	t.Run("elements", func(t *testing.T) {
		tc := newCase(t, `expected slice equality via .Equal method`)
		t.Cleanup(tc.assert)

		a := []*Person{
			{ID: 100, Name: "Alice"},
			{ID: 101, Name: "Bob"},
			{ID: 102, Name: "Carl"},
		}
		b := []*Person{
			{ID: 100, Name: "Alice"},
			{ID: 101, Name: "Bob"},
			{ID: 103, Name: "Dian"},
		}

		SliceEqual(tc, a, b)
	})
}

func TestLesser(t *testing.T) {
	tc := newCase(t, `expected val to be less via .Less method`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 200, Name: "Alice"}
	b := &Person{ID: 100, Name: "Bob"}

	Lesser(tc, b, a)
}

func TestLesser_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)

	a := &Person{ID: 200, Name: "Alice"}
	b := &Person{ID: 100, Name: "Bob"}

	Lesser(tc, b, a, tc.TestPostScript("lesser"))
}

func TestSliceEmpty(t *testing.T) {
	tc := newCase(t, `expected slice to be empty`)
	t.Cleanup(tc.assert)
	SliceEmpty(tc, []int{1, 2})
}

func TestSliceEmpty_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)
	SliceEmpty(tc, []int{1, 2}, tc.TestPostScript("empty slice"))
}

func TestSliceNotEmpty(t *testing.T) {
	tc := newCase(t, `expected slice to not be empty`)
	t.Cleanup(tc.assert)
	SliceNotEmpty(tc, []int{})
}

func TestSliceLen(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected slice to be different length`)
		t.Cleanup(tc.assert)
		SliceLen(tc, 2, []string{"a", "b", "c"})
	})

	t.Run("numbers", func(t *testing.T) {
		tc := newCase(t, `expected slice to be different length`)
		t.Cleanup(tc.assert)
		SliceLen(tc, 3, []int{8, 9})
	})
}

func TestSliceLen_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)
	SliceLen(tc, 3, []int{1, 2}, tc.TestPostScript("len slice"))
}

func TestLen(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected slice to be different length`)
		t.Cleanup(tc.assert)
		Len(tc, 2, []string{"a", "b", "c"})
	})

	t.Run("numbers", func(t *testing.T) {
		tc := newCase(t, `expected slice to be different length`)
		t.Cleanup(tc.assert)
		Len(tc, 3, []int{8, 9})
	})
}

func TestLen_PS(t *testing.T) {
	tc := newCapture(t)
	t.Cleanup(tc.post)
	Len(tc, 3, []int{1, 2}, tc.TestPostScript("len"))
}

func TestSliceContainsOp(t *testing.T) {
	t.Run("numbers", func(t *testing.T) {
		tc := newCase(t, `expected slice to contain missing item via == operator`)
		t.Cleanup(tc.assert)
		SliceContainsOp(tc, []int{3, 4, 5}, 7)
	})

	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected slice to contain missing item via == operator`)
		t.Cleanup(tc.assert)
		SliceContainsOp(tc, []string{"alice", "carl"}, "bob")
	})
}

func TestSliceContainsFunc(t *testing.T) {
	tc := newCase(t, `expected slice to contain missing item via 'eq' function`)
	t.Cleanup(tc.assert)

	s := []*Person{
		{ID: 100, Name: "Alice"},
		{ID: 101, Name: "Bob"},
	}

	SliceContainsFunc(tc, s, "Carl", func(a *Person, name string) bool {
		return a.Name == name
	})
}

func TestSliceContainsEqual(t *testing.T) {
	tc := newCase(t, `expected slice to contain missing item via .Equal method`)
	t.Cleanup(tc.assert)

	s := []*Person{
		{ID: 100, Name: "Alice"},
		{ID: 101, Name: "Bob"},
	}

	SliceContainsEqual(tc, s, &Person{ID: 102, Name: "Carl"})
}

func TestSliceContains(t *testing.T) {
	tc := newCase(t, `expected slice to contain missing item via cmp.Equal method`)
	t.Cleanup(tc.assert)

	s := []*Person{
		{ID: 100, Name: "Alice"},
		{ID: 101, Name: "Bob"},
	}

	SliceContains(tc, s, &Person{ID: 102, Name: "Carl"})
}

func TestSliceNotContains(t *testing.T) {
	tc := newCase(t, `expected slice to not contain item but it does`)
	t.Cleanup(tc.assert)

	s := []*Person{
		{ID: 100, Name: "Alice"},
		{ID: 101, Name: "Bob"},
		{ID: 102, Name: "Carla"},
	}

	SliceNotContains(tc, s, &Person{ID: 101, Name: "Bob"})
}

func TestSliceNotContainsFunc(t *testing.T) {
	tc := newCase(t, `expected slice to not contain item but it does`)
	t.Cleanup(tc.assert)

	s := []*Person{
		{ID: 100, Name: "Alice"},
		{ID: 101, Name: "Bob"},
		{ID: 102, Name: "Carla"},
	}

	f := func(a, b *Person) bool {
		return a.Name == b.Name && a.ID == b.ID
	}

	SliceNotContainsFunc(tc, s, &Person{ID: 101, Name: "Bob"}, f)
}

func TestSliceContainsAll(t *testing.T) {
	t.Run("wrong element", func(t *testing.T) {
		tc := newCase(t, `expected slice to contain missing item`)
		s := []*Person{
			{ID: 100, Name: "Alice"},
			{ID: 101, Name: "Bob"},
		}
		SliceContainsAll(tc, s, []*Person{{ID: 101, Name: "Bob"}, {ID: 105, Name: "Eve"}})
		t.Cleanup(tc.assert)
	})

	t.Run("too large", func(t *testing.T) {
		tc := newCase(t, `expected slice and items to contain same number of elements`)
		s := []*Person{
			{ID: 100, Name: "Alice"},
			{ID: 101, Name: "Bob"},
			{ID: 102, Name: "Carl"},
		}
		SliceContainsAll(tc, s, []*Person{{ID: 101, Name: "Bob"}, {ID: 100, Name: "Alice"}})
		t.Cleanup(tc.assert)
	})

	t.Run("too small", func(t *testing.T) {
		tc := newCase(t, `expected slice and items to contain same number of elements`)
		s := []*Person{
			{ID: 101, Name: "Bob"},
		}
		SliceContainsAll(tc, s, []*Person{{ID: 101, Name: "Bob"}, {ID: 100, Name: "Alice"}})
		t.Cleanup(tc.assert)
	})
}

func TestSliceContainsSubset(t *testing.T) {
	tc := newCase(t, `expected slice to contain missing item`)
	t.Cleanup(tc.assert)

	s := []*Person{
		{ID: 100, Name: "Alice"},
		{ID: 101, Name: "Bob"},
		{ID: 103, Name: "Carl"},
		{ID: 104, Name: "Dora"},
	}

	SliceContainsSubset(tc, s, []*Person{{ID: 101, Name: "Bob"}, {ID: 105, Name: "Eve"}})
}

func TestPositive(t *testing.T) {
	tc := newCase(t, `expected positive value`)
	t.Cleanup(tc.assert)

	Positive(tc, -1)
}

func TestNonPositive(t *testing.T) {
	tc := newCase(t, `expected non-positive value`)
	t.Cleanup(tc.assert)

	NonPositive(tc, 1)
}

func TestNegative(t *testing.T) {
	tc := newCase(t, `expected negative value`)
	t.Cleanup(tc.assert)

	Negative(tc, 1)
}

func TestNonNegative(t *testing.T) {
	tc := newCase(t, `expected non-negative value`)
	t.Cleanup(tc.assert)

	NonNegative(tc, -1)
}

func TestZero(t *testing.T) {
	tc := newCase(t, `expected value of 0`)
	t.Cleanup(tc.assert)

	Zero(tc, 1)
}

func TestNonZero(t *testing.T) {
	tc := newCase(t, `expected non-zero value`)
	t.Cleanup(tc.assert)

	NonZero(tc, 0)
}

func TestOne(t *testing.T) {
	tc := newCase(t, `expected value of 1`)
	t.Cleanup(tc.assert)

	One(tc, 1.1)
}

func TestLess(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		tc := newCase(t, `expected 7 < 5`)
		t.Cleanup(tc.assert)
		Less(tc, 5, 7)
	})

	t.Run("floats", func(t *testing.T) {
		tc := newCase(t, `expected 7.5 < 5.5`)
		t.Cleanup(tc.assert)
		Less(tc, 5.5, 7.5)
	})

	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected foo < bar`)
		t.Cleanup(tc.assert)
		Less(tc, "bar", "foo")
	})

	t.Run("equal", func(t *testing.T) {
		tc := newCase(t, `expected 7 < 7`)
		t.Cleanup(tc.assert)
		Less(tc, 7, 7)
	})
}

func TestLessEq(t *testing.T) {
	tc := newCase(t, `expected 7 ≤ 5`)
	t.Cleanup(tc.assert)
	LessEq(tc, 5, 7)
}

func TestGreater(t *testing.T) {
	t.Run("integer", func(t *testing.T) {
		tc := newCase(t, `expected 5 > 7`)
		t.Cleanup(tc.assert)
		Greater(tc, 7, 5)
	})

	t.Run("floats", func(t *testing.T) {
		tc := newCase(t, `expected 5.5 > 7.7`)
		t.Cleanup(tc.assert)
		Greater(tc, 7.7, 5.5)
	})

	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected bar > foo`)
		t.Cleanup(tc.assert)
		Greater(tc, "foo", "bar")
	})

	t.Run("equal", func(t *testing.T) {
		tc := newCase(t, `expected bar > bar`)
		t.Cleanup(tc.assert)
		Greater(tc, "bar", "bar")
	})
}

func TestGreaterEq(t *testing.T) {
	tc := newCase(t, `expected 5 ≥ 7`)
	t.Cleanup(tc.assert)
	GreaterEq(tc, 7, 5)
}

func TestBetween(t *testing.T) {
	t.Run("too high", func(t *testing.T) {
		tc := newCase(t, `expected val in range (3 ≤ val ≤ 5)`)
		t.Cleanup(tc.assert)
		Between(tc, 3, 7, 5)
	})

	t.Run("too low", func(t *testing.T) {
		tc := newCase(t, `expected val in range (3 ≤ val ≤ 5)`)
		t.Cleanup(tc.assert)
		Between(tc, 3, 1, 5)
	})
}

func TestBetweenExclusive(t *testing.T) {
	t.Run("too high", func(t *testing.T) {
		tc := newCase(t, `expected val in range (3 < val < 5)`)
		t.Cleanup(tc.assert)
		BetweenExclusive(tc, 3, 7, 5)
	})

	t.Run("too high fence", func(t *testing.T) {
		tc := newCase(t, `expected val in range (3 < val < 5)`)
		t.Cleanup(tc.assert)
		BetweenExclusive(tc, 3, 5, 5)
	})

	t.Run("too low", func(t *testing.T) {
		tc := newCase(t, `expected val in range (3 < val < 5)`)
		t.Cleanup(tc.assert)
		BetweenExclusive(tc, 3, 1, 5)
	})

	t.Run("too low fence", func(t *testing.T) {
		tc := newCase(t, `expected val in range (3 < val < 5)`)
		t.Cleanup(tc.assert)
		BetweenExclusive(tc, 3, 3, 5)
	})
}

type number int

func (n number) Min() number {
	return n
}

func (n number) Max() number {
	return n
}

func TestMin(t *testing.T) {
	tc := newCase(t, `expected a different value for min`)
	t.Cleanup(tc.assert)

	Min[number](tc, 42, number(100))
}

func TestMax(t *testing.T) {
	tc := newCase(t, `expected a different value for max`)
	t.Cleanup(tc.assert)

	Max[number](tc, 71, number(100))
}

func TestAscending(t *testing.T) {
	t.Run("numbers", func(t *testing.T) {
		tc := newCase(t, `expected slice`)
		t.Cleanup(tc.assert)

		l := []int{1, 2, 3, 5, 4}
		Ascending(tc, l)
	})

	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected slice`)
		t.Cleanup(tc.assert)

		l := []string{"alpha", "beta", "gamma", "delta"}
		Ascending(tc, l)
	})
}

func TestAscendingFunc(t *testing.T) {
	tc := newCase(t, `expected less`)
	t.Cleanup(tc.assert)

	l := []*Person{
		{ID: 100, Name: "Alice"},
		{ID: 200, Name: "Bob"},
		{ID: 300, Name: "Dale"},
		{ID: 400, Name: "Carl"},
	}

	AscendingFunc(tc, l, func(p1, p2 *Person) bool {
		return p1.Name < p2.Name
	})
}

func TestAscendingCmp(t *testing.T) {
	tc := newCase(t, `expected compare`)
	t.Cleanup(tc.assert)

	l := []*Person{
		{ID: 100, Name: "Alice"},
		{ID: 200, Name: "Bob"},
		{ID: 300, Name: "Dale"},
		{ID: 400, Name: "Carl"},
	}

	cmp := func(p1, p2 *Person) int {
		switch {
		case p1.Name < p2.Name:
			return -1
		case p1.Name > p2.Name:
			return 1
		default:
			return 0
		}
	}

	AscendingCmp(tc, l, cmp)
}

func TestAscendingLess(t *testing.T) {
	tc := newCase(t, `expected slice`)
	t.Cleanup(tc.assert)

	l := []*Person{
		{ID: 100, Name: "Alice"},
		{ID: 200, Name: "Bob"},
		{ID: 150, Name: "Carl"},
	}
	AscendingLess(tc, l)
}

func TestDescending(t *testing.T) {
	t.Run("numbers", func(t *testing.T) {
		tc := newCase(t, `expected slice`)
		t.Cleanup(tc.assert)

		l := []int{6, 5, 3, 4, 1}
		Descending(tc, l)
	})

	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected slice`)
		t.Cleanup(tc.assert)

		l := []string{"zoo", "yard", "boat", "xray"}
		Descending(tc, l)
	})
}

func TestDescendingFunc(t *testing.T) {
	tc := newCase(t, `expected less`)
	t.Cleanup(tc.assert)
	l := []*Person{
		{ID: 400, Name: "Dale"},
		{ID: 300, Name: "Bob"},
		{ID: 200, Name: "Carl"},
		{ID: 100, Name: "Alice"},
	}

	DescendingFunc(tc, l, func(p1, p2 *Person) bool {
		return p1.Name < p2.Name
	})
}

func TestDescendingCmp(t *testing.T) {
	tc := newCase(t, `expected compare`)
	t.Cleanup(tc.assert)
	l := []*Person{
		{ID: 400, Name: "Dale"},
		{ID: 300, Name: "Bob"},
		{ID: 200, Name: "Carl"},
		{ID: 100, Name: "Alice"},
	}

	cmp := func(p1, p2 *Person) int {
		switch {
		case p1.Name < p2.Name:
			return -1
		case p1.Name > p2.Name:
			return 1
		default:
			return 0
		}
	}

	DescendingCmp(tc, l, cmp)
}

func TestDescendingLess(t *testing.T) {
	tc := newCase(t, `expected slice`)
	t.Cleanup(tc.assert)
	l := []*Person{
		{ID: 400, Name: "Dale"},
		{ID: 300, Name: "Carl"},
		{ID: 100, Name: "Bob"},
		{ID: 200, Name: "Alice"},
	}
	DescendingLess(tc, l)
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
		tc := newCase(t, `expected slices of same length`)
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
		tc := newCase(t, `expected maps of same length`)
		t.Cleanup(tc.assert)
		a := map[string]int{"a": 1}
		b := map[string]int{"a": 1, "b": 2}
		MapEq(tc, a, b)
	})

	t.Run("different keys", func(t *testing.T) {
		tc := newCase(t, `expected maps of same keys`)
		t.Cleanup(tc.assert)
		a := map[int]string{1: "a", 2: "b"}
		b := map[int]string{1: "a", 3: "c"}
		MapEq(tc, a, b)
	})

	t.Run("different values", func(t *testing.T) {
		tc := newCase(t, `expected maps of same values via cmp.Equal function`)
		t.Cleanup(tc.assert)
		a := map[string]string{"a": "amp", "b": "bar"}
		b := map[string]string{"a": "amp", "b": "foo"}
		MapEq(tc, a, b)
	})

	t.Run("custom types", func(t *testing.T) {
		tc := newCase(t, `expected maps of same values via cmp.Equal function`)
		t.Cleanup(tc.assert)

		type custom1 map[string]int
		a := custom1{"key": 1}
		type custom2 map[string]int
		b := custom2{"key": 2}
		MapEq(tc, a, b)
	})
}

func TestMapEqFunc(t *testing.T) {
	t.Run("different values", func(t *testing.T) {
		tc := newCase(t, `expected maps of same values via 'eq' function`)
		t.Cleanup(tc.assert)

		a := map[int]Person{
			0: {ID: 100, Name: "Alice"},
			1: {ID: 101, Name: "Bob"},
		}

		b := map[int]Person{
			0: {ID: 100, Name: "Alice"},
			1: {ID: 101, Name: "Bob B."},
		}

		MapEqFunc(tc, a, b, func(p1, p2 Person) bool {
			return p1.ID == p2.ID && p1.Name == p2.Name
		})
	})
}

func TestMapEqual(t *testing.T) {
	t.Run("different values", func(t *testing.T) {
		tc := newCase(t, `expected maps of same values via .Equal method`)
		t.Cleanup(tc.assert)

		a := map[int]*Person{
			0: {ID: 100, Name: "Alice"},
			1: {ID: 101, Name: "Bob"},
		}

		b := map[int]*Person{
			0: {ID: 100, Name: "Alice"},
			1: {ID: 200, Name: "Bob"},
		}

		MapEqual(tc, a, b)
	})
}

func TestMapLen(t *testing.T) {
	tc := newCase(t, `expected map to be different length`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	MapLen(tc, 2, m)
}

func TestMapEmpty(t *testing.T) {
	tc := newCase(t, `expected map to be empty`)
	t.Cleanup(tc.assert)
	m := map[string]int{"a": 1, "b": 2}
	MapEmpty(tc, m)
}

func TestMapEmptyCustom(t *testing.T) {
	tc := newCase(t, `expected map to be empty`)
	t.Cleanup(tc.assert)
	type custom map[string]int
	m := make(custom)
	m["a"] = 1
	m["b"] = 2
	MapEmpty(tc, m)
}

func TestMapNotEmpty(t *testing.T) {
	tc := newCase(t, `expected map to not be empty`)
	t.Cleanup(tc.assert)
	m := make(map[string]string)
	MapNotEmpty(tc, m)
}

func TestMapContainsKey(t *testing.T) {
	tc := newCase(t, `expected map to contain key`)
	t.Cleanup(tc.assert)
	m := map[string]int{"a": 1, "b": 2}
	MapContainsKey(tc, m, "c")
}

func TestMapNotContainsKey(t *testing.T) {
	tc := newCase(t, `expected map to not contain key`)
	t.Cleanup(tc.assert)
	m := map[string]int{"a": 1, "b": 2}
	MapNotContainsKey(tc, m, "b")
}

func TestMapContainsKeys(t *testing.T) {
	tc := newCase(t, `expected map to contain keys`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	MapContainsKeys(tc, m, []string{"z", "a", "b", "c", "d"})
}

func TestMapNotContainsKeys(t *testing.T) {
	tc := newCase(t, `expected map to not contain keys`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	MapNotContainsKeys(tc, m, []string{"z", "b", "y", "c"})
}

func TestMapContainsValues(t *testing.T) {
	tc := newCase(t, `expected map to contain values`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	MapContainsValues(tc, m, []int{9, 1, 2, 7})
}

func TestMapNotContainsValues(t *testing.T) {
	tc := newCase(t, `expected map to not contain values`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	MapNotContainsValues(tc, m, []int{9, 8, 2, 7})
}

func TestMapContainsValuesFunc(t *testing.T) {
	tc := newCase(t, `expected map to contain values`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	MapContainsValuesFunc(tc, m, []int{9, 1, 2, 7}, func(a, b int) bool {
		return a == b
	})
}

func TestMapNotContainsValuesFunc(t *testing.T) {
	tc := newCase(t, `expected map to not contain values`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	MapNotContainsValuesFunc(tc, m, []int{2, 4, 6, 8}, func(a, b int) bool {
		return a == b
	})
}

func TestMapContainsValuesEqual(t *testing.T) {
	tc := newCase(t, `expected map to contain values`)
	t.Cleanup(tc.assert)

	m := map[int]*Person{
		1: {ID: 100, Name: "Alice"},
		2: {ID: 200, Name: "Bob"},
		3: {ID: 300, Name: "Carl"},
	}
	MapContainsValuesEqual(tc, m, []*Person{
		{ID: 201, Name: "Bob"},
	})
}

func TestMapNotContainsValuesEqual(t *testing.T) {
	tc := newCase(t, `expected map to not contain values`)
	t.Cleanup(tc.assert)

	m := map[int]*Person{
		1: {ID: 100, Name: "Alice"},
		2: {ID: 200, Name: "Bob"},
		3: {ID: 300, Name: "Carl"},
	}
	MapNotContainsValuesEqual(tc, m, []*Person{
		{ID: 201, Name: "Bob"}, {ID: 200, Name: "Daisy"},
	})
}

func TestFileExistsFS(t *testing.T) {
	tc := newCase(t, `expected file to exist`)
	t.Cleanup(tc.assert)

	FileExistsFS(tc, os.DirFS("/etc"), "hosts2")
}

func TestFileExists(t *testing.T) {
	tc := newCase(t, `expected file to exist`)
	t.Cleanup(tc.assert)

	FileExists(tc, "/etc/hosts2")
}

func TestFileNotExistsFS(t *testing.T) {
	needsOS(t, "linux")

	tc := newCase(t, `expected file to not exist`)
	t.Cleanup(tc.assert)

	FileNotExistsFS(tc, os.DirFS("/etc"), "hosts")
}

func TestFileNotExists(t *testing.T) {
	needsOS(t, "linux")

	tc := newCase(t, `expected file to not exist`)
	t.Cleanup(tc.assert)

	FileNotExists(tc, "/etc/hosts")
}

func TestDirExistsFS(t *testing.T) {
	tc := newCase(t, `expected directory to exist`)
	t.Cleanup(tc.assert)

	DirExistsFS(tc, os.DirFS("/usr/local"), "bin2")
}

func TestDirExists(t *testing.T) {
	tc := newCase(t, `expected directory to exist`)
	t.Cleanup(tc.assert)

	DirExists(tc, "/usr/local/bin2")
}

func TestDirNotExistsFS(t *testing.T) {
	needsOS(t, "linux")

	tc := newCase(t, `expected directory to not exist`)
	t.Cleanup(tc.assert)

	DirNotExistsFS(tc, os.DirFS("/usr"), "local")
}

func TestDirNotExists(t *testing.T) {
	needsOS(t, "linux")

	tc := newCase(t, `expected directory to not exist`)
	t.Cleanup(tc.assert)

	DirNotExists(tc, "/usr/local")
}

func TestFileModeFS(t *testing.T) {
	needsOS(t, "linux")

	tc := newCase(t, `expected different file permissions`)
	t.Cleanup(tc.assert)

	var unexpected os.FileMode = 0673 // (actual 0655)
	FileModeFS(tc, os.DirFS("/bin"), "find", unexpected)
}

func TestFileMode(t *testing.T) {
	needsOS(t, "linux")

	tc := newCase(t, `expected different file permissions`)
	t.Cleanup(tc.assert)

	var unexpected os.FileMode = 0673 // (actual 0655)
	FileMode(tc, "/bin/find", unexpected)
}

func TestFileContainsFS(t *testing.T) {
	needsOS(t, "linux")

	tc := newCase(t, `expected file contents`)
	t.Cleanup(tc.assert)

	FileContainsFS(tc, os.DirFS("/etc"), "hosts", "127.0.0.999")
}

func TestFileContains(t *testing.T) {
	needsOS(t, "linux")

	tc := newCase(t, `expected file contents`)
	t.Cleanup(tc.assert)

	FileContains(tc, "/etc/hosts", "127.0.0.999")
}

func TestFilePathValid(t *testing.T) {
	tc := newCase(t, `expected valid file path`)
	t.Cleanup(tc.assert)

	FilePathValid(tc, "foo/../bar")
}

func TestStrEqFold(t *testing.T) {
	tc := newCase(t, `expected strings to be equal ignoring case`)
	t.Cleanup(tc.assert)

	StrEqFold(tc, "hello", "hi")
}

func TestStrNotEqFold(t *testing.T) {
	tc := newCase(t, `expected strings to not be equal ignoring case; but they are`)
	t.Cleanup(tc.assert)

	StrNotEqFold(tc, "hello", "HeLLo")
}

func TestStrContains(t *testing.T) {
	tc := newCase(t, `expected string to contain substring; it does not`)
	t.Cleanup(tc.assert)

	StrContains(tc, "banana", "band")
}

func TestStrContainsFold(t *testing.T) {
	tc := newCase(t, `expected string to contain substring; it does not`)
	t.Cleanup(tc.assert)

	StrContainsFold(tc, "banana", "band")
}

func TestStrNotContains(t *testing.T) {
	tc := newCase(t, `expected string to not contain substring; but it does`)
	t.Cleanup(tc.assert)

	StrNotContains(tc, "banana", "ana")
}

func TestStrNotContainsFold(t *testing.T) {
	tc := newCase(t, `expected string to not contain substring; but it does`)
	t.Cleanup(tc.assert)

	StrNotContainsFold(tc, "banana", "aNA")
}

func TestStrContainsAny(t *testing.T) {
	tc := newCase(t, `expected string to contain one or more code points`)
	t.Cleanup(tc.assert)

	StrContainsAny(tc, "banana", "xyz")
}

func TestStrNotContainsAny(t *testing.T) {
	tc := newCase(t, `expected string to not contain code points; but it does`)
	t.Cleanup(tc.assert)

	StrNotContainsAny(tc, "banana", "xnz")
}

func TestStrCount(t *testing.T) {
	tc := newCase(t, `expected string to contain`)
	t.Cleanup(tc.assert)

	StrCount(tc, "mississippi", "ssi", 3)
}

func TestStrContainsFields(t *testing.T) {
	tc := newCase(t, `expected fields of string to contain subset of values`)
	t.Cleanup(tc.assert)

	StrContainsFields(tc, "one too three", []string{"one", "two", "three", "nine"})
}

func TestStrHasPrefix(t *testing.T) {
	tc := newCase(t, `expected string to have prefix`)
	t.Cleanup(tc.assert)

	StrHasPrefix(tc, "mrs.", "mr. biggles")
}

func TestStrNotHasPrefix(t *testing.T) {
	tc := newCase(t, `expected string to not have prefix; but it does`)
	t.Cleanup(tc.assert)

	StrNotHasPrefix(tc, "mr.", "mr. biggles")
}

func TestStrHasSuffix(t *testing.T) {
	tc := newCase(t, `expected string to have suffix`)
	t.Cleanup(tc.assert)

	StrHasSuffix(tc, "wiggles", "mr. biggles")
}

func TestStringNotHasSuffix(t *testing.T) {
	tc := newCase(t, `expected string to not have suffix; but it does`)
	t.Cleanup(tc.assert)

	StrNotHasSuffix(tc, "biggles", "mr. biggles")
}

func TestRegexMatch(t *testing.T) {
	tc := newCase(t, `expected regexp match`)
	t.Cleanup(tc.assert)

	re := regexp.MustCompile(`abc\d`)
	RegexMatch(tc, re, "abcX")
}

func TestRegexCompiles(t *testing.T) {
	tc := newCase(t, `expected regular expression to compile`)
	t.Cleanup(tc.assert)

	RegexCompiles(tc, "ab"+`\`+"ef")
}

func TestRegexCompilesPOSIX(t *testing.T) {
	tc := newCase(t, `expected regular expression to compile (posix)`)
	t.Cleanup(tc.assert)

	RegexCompilesPOSIX(tc, "ab"+`\`+"ef")
}

func TestPS_Sprintf(t *testing.T) {
	tc := newCapture(t)
	Eq(tc, "a", "b", Sprintf("hello %s", "world"))
}

func TestPS_Sprint(t *testing.T) {
	tc := newCapture(t)
	Eq(tc, "a", "b", Sprint("hello", 42, "hi"))
}

func TestPS_Values(t *testing.T) {
	tc := newCapture(t)
	Eq(tc, "a", "b", Values("foo", "bar", 1, 2, "now", time.Now()))
}

func TestPS_Func(t *testing.T) {
	tc := newCapture(t)
	Eq(tc, "a", "b", Func(func() string {
		return "hello"
	}))
}

func TestEq_Combo(t *testing.T) {
	tc := newCapture(t)
	Eq(tc, "a", "b", Sprintf("this is a note"), Values("foo", "bar", "baz", 3), Func(func() string {
		return "this is the result of a function"
	}))
}

func Test_UUIDv4(t *testing.T) {
	tc := newCase(t, `expected well-formed v4 UUID`)
	t.Cleanup(tc.assert)

	UUIDv4(tc, "abc123")                              // fail
	UUIDv4(t, "12345678-abcd-1234-abcd-aabbccdd1122") // pass
}

type container[T any] struct {
	contains bool
	empty    bool
	size     int
	length   int
}

func (c *container[T]) Contains(_ T) bool {
	return c.contains
}

func (c *container[T]) Empty() bool {
	return c.empty
}

func (c *container[T]) Size() int {
	return c.size
}

func (c *container[T]) Len() int {
	return c.length
}

func (c *container[T]) Copy() *container[T] {
	return &container[T]{
		contains: c.contains,
		empty:    c.empty,
		size:     c.size,
		length:   c.length,
	}
}

func (c *container[T]) Equal(o *container[T]) bool {
	if c == nil || o == nil {
		return c == o
	}
	switch {
	case c.contains != o.contains:
		return false
	case c.empty != o.empty:
		return false
	case c.size != o.size:
		return false
	case c.length != o.length:
		return false
	}
	return true
}

func TestEmpty(t *testing.T) {
	tc := newCase(t, `expected to be empty, but was not`)
	t.Cleanup(tc.assert)

	c := &container[struct{}]{empty: false}
	Empty(tc, c)
}

func TestNotEmpty(t *testing.T) {
	tc := newCase(t, `expected to not be empty, but is`)
	t.Cleanup(tc.assert)

	c := &container[struct{}]{empty: true}
	NotEmpty(tc, c)
}

func TestContains(t *testing.T) {
	tc := newCase(t, `expected to contain element, but does not`)
	t.Cleanup(tc.assert)

	c := &container[string]{contains: false}
	Contains[string](tc, "apple", c)
}

func TestContainsSubset(t *testing.T) {
	tc := newCase(t, `expected to contain element, but does not`)
	t.Cleanup(tc.assert)

	c := &container[string]{contains: false}
	ContainsSubset[string](tc, []string{"a", "b"}, c)
}

func TestNotContains(t *testing.T) {
	tc := newCase(t, `expected not to contain element, but it does`)
	t.Cleanup(tc.assert)

	c := &container[string]{contains: true}
	NotContains[string](tc, "apple", c)
}

func TestSize(t *testing.T) {
	tc := newCase(t, `expected different size`)
	t.Cleanup(tc.assert)

	c := &container[string]{size: 3}
	Size(tc, 2, c)
}

func TestLength(t *testing.T) {
	tc := newCase(t, `expected different length`)
	t.Cleanup(tc.assert)

	c := &container[string]{length: 3}
	Length(tc, 4, c)
}

func TestWait_BoolFunc(t *testing.T) {
	tc := newCase(t, `expected condition to pass within wait context`)
	t.Cleanup(tc.assert)

	Wait(tc, wait.InitialSuccess(
		wait.BoolFunc(func() bool { return false }),
		wait.Timeout(100*time.Millisecond),
	))
}

func TestWait_ErrorFunc(t *testing.T) {
	tc := newCase(t, `expected condition to pass within wait context`)
	t.Cleanup(tc.assert)

	Wait(tc, wait.InitialSuccess(
		wait.ErrorFunc(func() error { return errors.New("fail") }),
		wait.Timeout(100*time.Millisecond),
	))
}

func TestWait_TestFunc(t *testing.T) {
	tc := newCase(t, `expected condition to pass within wait context`)
	t.Cleanup(tc.assert)

	Wait(tc, wait.InitialSuccess(
		wait.TestFunc(func() (bool, error) { return false, errors.New("fail") }),
		wait.Timeout(100*time.Millisecond),
	))
}

func TestStructEqual(t *testing.T) {
	tc := newCase(t, `expected inequality via .Equal method`)
	t.Cleanup(tc.assert)

	StructEqual(tc, &container[int]{
		contains: true,
		empty:    true,
		size:     1,
		length:   2,
	}, []Tweak[*container[int]]{{
		Field: "contains",
		Apply: func(c *container[int]) { c.contains = false },
	}, {
		Field: "empty",
		Apply: func(c *container[int]) { c.empty = false },
	}, {
		Field: "size",
		Apply: func(c *container[int]) { c.size = 9 },
	}, {
		Field: "length",
		Apply: func(c *container[int]) { c.length = 2 }, // no mod
	}})
}

func TestStructEqual_Tweaks(t *testing.T) {
	tc := newCase(t, `expected inequality via .Equal method`)
	t.Cleanup(tc.assert)

	StructEqual(tc, &container[int]{
		contains: true,
		empty:    true,
		size:     1,
		length:   2,
	}, Tweaks[*container[int]]{{
		Field: "contains",
		Apply: func(c *container[int]) { c.contains = false },
	}, {
		Field: "empty",
		Apply: func(c *container[int]) { c.empty = false },
	}, {
		Field: "size",
		Apply: func(c *container[int]) { c.size = 9 },
	}, {
		Field: "length",
		Apply: func(c *container[int]) { c.length = 2 }, // no mod
	}})
}

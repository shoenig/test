package test

import (
	"errors"
	"math"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestNil(t *testing.T) {
	tc := newCase(t, `expected to be nil; is not nil`)
	t.Cleanup(tc.assert)

	Nil(tc, 42)
	Nil(tc, "hello")
	Nil(tc, time.UTC)
	Nil(tc, []string{"foo"})
	Nil(tc, map[string]int{"foo": 1})
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

func TestTrue(t *testing.T) {
	tc := newCase(t, `expected condition to be true; is false`)
	t.Cleanup(tc.assert)

	True(tc, false)
}

func TestFalse(t *testing.T) {
	tc := newCase(t, `expected condition to be false; is true`)
	t.Cleanup(tc.assert)

	False(tc, true)
}

func TestUnreachable(t *testing.T) {
	tc := newCase(t, `expected not to execute this code path`)
	t.Cleanup(tc.assert)

	Unreachable(tc)
}

func TestError(t *testing.T) {
	tc := newCase(t, `expected non-nil error; is nil`)
	t.Cleanup(tc.assert)

	Error(tc, nil)
}

func TestEqError(t *testing.T) {
	tc := newCase(t, `expected matching error strings`)
	t.Cleanup(tc.assert)

	EqError(tc, errors.New("oops"), "blah")
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

func TestEqOp(t *testing.T) {
	t.Run("number", func(t *testing.T) {
		tc := newCase(t, `expected equality via ==`)
		t.Cleanup(tc.assert)
		EqOp(tc, "foo", "bar")
	})
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

func TestNotEq(t *testing.T) {
	tc := newCase(t, `expected inequality via cmp.Equal function`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 100, Name: "Alice"}

	NotEq(tc, a, b)
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

func TestNotEqFunc(t *testing.T) {
	tc := newCase(t, `expected inequality via 'eq' function`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 100, Name: "Alice"}

	NotEqFunc(tc, a, b, func(a, b *Person) bool {
		return a.ID == b.ID && a.Name == b.Name
	})
}

func TestEqJSON(t *testing.T) {
	tc := newCase(t, `expected equality via json marshalling`)
	t.Cleanup(tc.assert)

	EqJSON(tc, `{"a":1, "b":2}`, `{"b":2, "a":9}`)
}

func TestEqSliceFunc(t *testing.T) {
	t.Run("length", func(t *testing.T) {
		tc := newCase(t, `expected slices of same length`)
		t.Cleanup(tc.assert)

		a := []int{1, 2, 3}
		b := []int{1, 2}
		EqSliceFunc(tc, a, b, func(a, b int) bool {
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

		EqSliceFunc(tc, a, b, func(a, b *Person) bool {
			return a.ID == b.ID
		})
	})
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
	tc := newCase(t, `expected equality via .Equals method`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 150, Name: "Alice"}

	Equals(tc, a, b)
}

func TestNotEquals(t *testing.T) {
	tc := newCase(t, `expected inequality via .Equals method`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 100, Name: "Alice"}
	b := &Person{ID: 100, Name: "Alice"}

	NotEquals(tc, a, b)
}

func TestEqualsSlice(t *testing.T) {
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
		EqualsSlice(tc, a, b)
	})

	t.Run("elements", func(t *testing.T) {
		tc := newCase(t, `expected slice equality via .Equals method`)
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

		EqualsSlice(tc, a, b)
	})
}

func TestLesser(t *testing.T) {
	tc := newCase(t, `expected to be less via .Less method`)
	t.Cleanup(tc.assert)

	a := &Person{ID: 200, Name: "Alice"}
	b := &Person{ID: 100, Name: "Bob"}

	Lesser(tc, a, b)
}

func TestEmptySlice(t *testing.T) {
	tc := newCase(t, `expected slice to be empty`)
	t.Cleanup(tc.assert)

	EmptySlice(tc, []int{1, 2})
}

func TestEmpty(t *testing.T) {
	tc := newCase(t, `expected slice to be empty`)
	t.Cleanup(tc.assert)

	Empty(tc, []int{1, 2})
}

func TestLenSlice(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected slice to be different length`)
		t.Cleanup(tc.assert)
		LenSlice(tc, 2, []string{"a", "b", "c"})
	})

	t.Run("numbers", func(t *testing.T) {
		tc := newCase(t, `expected slice to be different length`)
		t.Cleanup(tc.assert)
		LenSlice(tc, 3, []int{8, 9})
	})
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

func TestContains(t *testing.T) {
	t.Run("people", func(t *testing.T) {
		tc := newCase(t, `expected slice to contain missing item via cmp.Equal function`)
		t.Cleanup(tc.assert)
		a := []*Person{
			{ID: 100, Name: "Alice"},
			{ID: 101, Name: "Bob"},
		}
		target := &Person{ID: 102, Name: "Carl"}
		Contains(tc, a, target)
	})
}

func TestContainsOp(t *testing.T) {
	t.Run("numbers", func(t *testing.T) {
		tc := newCase(t, `expected slice to contain missing item via == operator`)
		t.Cleanup(tc.assert)
		ContainsOp(tc, []int{3, 4, 5}, 7)
	})

	t.Run("strings", func(t *testing.T) {
		tc := newCase(t, `expected slice to contain missing item via == operator`)
		t.Cleanup(tc.assert)
		ContainsOp(tc, []string{"alice", "carl"}, "bob")
	})
}

func TestContainsFunc(t *testing.T) {
	tc := newCase(t, `expected slice to contain missing item via 'eq' function`)
	t.Cleanup(tc.assert)

	s := []*Person{
		{ID: 100, Name: "Alice"},
		{ID: 101, Name: "Bob"},
	}

	ContainsFunc(tc, s, &Person{ID: 102, Name: "Carl"}, func(a, b *Person) bool {
		return a.ID == b.ID && a.Name == b.Name
	})
}

func TestContainsEquals(t *testing.T) {
	tc := newCase(t, `expected slice to contain missing item via .Equals method`)
	t.Cleanup(tc.assert)

	s := []*Person{
		{ID: 100, Name: "Alice"},
		{ID: 101, Name: "Bob"},
	}

	ContainsEquals(tc, s, &Person{ID: 102, Name: "Carl"})
}

func TestContainsString(t *testing.T) {
	tc := newCase(t, `expected to contain substring`)
	t.Cleanup(tc.assert)

	ContainsString(tc, "foobar", "food")
}

func TestPositive(t *testing.T) {
	tc := newCase(t, `expected positive value`)
	t.Cleanup(tc.assert)

	Positive(tc, -1)
}

func TestNegative(t *testing.T) {
	tc := newCase(t, `expected negative value`)
	t.Cleanup(tc.assert)

	Negative(tc, 1)
}

func TestZero(t *testing.T) {
	tc := newCase(t, `expected zero`)
	t.Cleanup(tc.assert)

	Zero(tc, 1)
}

func TestNonZero(t *testing.T) {
	tc := newCase(t, `expected non-zero`)
	t.Cleanup(tc.assert)

	NonZero(tc, 0)
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

func TestLessEq(t *testing.T) {
	tc := newCase(t, `expected 7 <= 5`)
	t.Cleanup(tc.assert)
	LessEq(tc, 7, 5)
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

func TestGreaterEq(t *testing.T) {
	tc := newCase(t, `expected 5 >= 7`)
	t.Cleanup(tc.assert)
	GreaterEq(tc, 5, 7)
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
		tc := newCase(t, `expected maps of same values via cmp.Diff function`)
		t.Cleanup(tc.assert)
		a := map[string]string{"a": "amp", "b": "bar"}
		b := map[string]string{"a": "amp", "b": "foo"}
		MapEq(tc, a, b)
	})

	t.Run("custom types", func(t *testing.T) {
		tc := newCase(t, `expected maps of same values via cmp.Diff function`)
		t.Cleanup(tc.assert)

		type custom1 map[string]int
		a := custom1{"key": 1}
		type custom2 map[string]int
		b := custom2{"key": 2}
		MapEq(tc, a, b)
	})
}

func TestMapEqFunc(t *testing.T) {
	t.Run("different value", func(t *testing.T) {
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

func TestMapEquals(t *testing.T) {
	t.Run("different value", func(t *testing.T) {
		tc := newCase(t, `expected maps of same values via .Equals method`)
		t.Cleanup(tc.assert)

		a := map[int]*Person{
			0: {ID: 100, Name: "Alice"},
			1: {ID: 101, Name: "Bob"},
		}

		b := map[int]*Person{
			0: {ID: 100, Name: "Alice"},
			1: {ID: 200, Name: "Bob"},
		}

		MapEquals(tc, a, b)
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

func TestMapContainsKeys(t *testing.T) {
	tc := newCase(t, `expected map to contain keys`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	MapContainsKeys(tc, m, []string{"z", "a", "b", "c", "d"})
}

func TestMapContainsValues(t *testing.T) {
	tc := newCase(t, `expected map to contain values`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	MapContainsValues(tc, m, []int{9, 1, 2, 7})
}

func TestMapContainsValuesFunc(t *testing.T) {
	tc := newCase(t, `expected map to contain values`)
	t.Cleanup(tc.assert)

	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	MapContainsValuesFunc(tc, m, []int{9, 1, 2, 7}, func(a, b int) bool {
		return a == b
	})
}

func TestMapContainsValuesEquals(t *testing.T) {
	tc := newCase(t, `expected map to contain values`)
	t.Cleanup(tc.assert)

	m := map[int]*Person{
		1: {ID: 100, Name: "Alice"},
		2: {ID: 200, Name: "Bob"},
		3: {ID: 300, Name: "Carl"},
	}
	MapContainsValuesEquals(tc, m, []*Person{
		{ID: 201, Name: "Bob"},
	})
}

func TestFileExists(t *testing.T) {
	tc := newCase(t, `expected file to exist`)
	t.Cleanup(tc.assert)

	FileExists(tc, os.DirFS("/etc"), "hosts2")
}

func TestFileNotExists(t *testing.T) {
	tc := newCase(t, `expected file to not exist`)
	t.Cleanup(tc.assert)

	FileNotExists(tc, os.DirFS("/etc"), "hosts")
}

func TestDirExists(t *testing.T) {
	tc := newCase(t, `expected directory to exist`)
	t.Cleanup(tc.assert)

	DirExists(tc, os.DirFS("/usr/local"), "bin2")
}

func TestDirNotExists(t *testing.T) {
	tc := newCase(t, `expected directory to not exist`)
	t.Cleanup(tc.assert)

	DirNotExists(tc, os.DirFS("/usr"), "local")
}

func TestFileMode(t *testing.T) {
	tc := newCase(t, `expected different file permissions`)
	t.Cleanup(tc.assert)

	FileMode(tc, os.DirFS("/bin"), "find", 0673) // (actual 0655)
}

func TestFileContains(t *testing.T) {
	tc := newCase(t, `expected file contents`)
	t.Cleanup(tc.assert)

	FileContains(tc, os.DirFS("/etc"), "hosts", "127.0.0.999")
}

func TestFilePathValid(t *testing.T) {
	tc := newCase(t, `expected valid file path`)
	t.Cleanup(tc.assert)

	FilePathValid(tc, "foo/../bar")
}

func TestRegexMatch(t *testing.T) {
	tc := newCase(t, `expected regexp match`)
	t.Cleanup(tc.assert)

	re := regexp.MustCompile(`abc\d`)
	RegexMatch(tc, re, "abcX")
}

# test

<img align="right" width="240" height="244" src="https://i.imgur.com/gmn5mIo.png">

[![GoDoc](https://godoc.org/github.com/shoenig/test?status.svg)](https://godoc.org/github.com/shoenig/test)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoenig/test)](https://goreportcard.com/report/github.com/shoenig/test)
[![CI Tests](https://github.com/shoenig/test/actions/workflows/ci.yaml/badge.svg)](https://github.com/shoenig/test/actions/workflows/ci.yaml)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

`test` is a generics based testing assertions library for Go.

There are four key packages,

- `must` - assertions causing test failure and halt the test case immediately
- `test` - assertions causing test failure and allow the test case to continue
- `wait` - utilities for waiting on conditionals in tests
- `portal` - utilities for allocating free ports for network listeners in tests

### Changes

:rocket: v0.5.0 contains new packages `wait` and `portal`

:warning: v0.4.0 contains breaking changes

 - Slice functions are renamed to be more consistent and to make room for interface based variants.
 - Filesystem assertions now use OS by default, FS interface methods are renamed.
 - Comparison assertions now have the expectation parameter first.

### Requirements

Only depends on `github.com/google/go-cmp`.

The minimum Go version is `go1.18`.

### Install

Use `go get` to grab the latest version of `test`.

```shell
go get -u github.com/shoenig/test@latest
```

### Influence

This library was made after a ~decade of using [testify](https://github.com/stretchr/testify),
quite possibly the most used library in the whole Go ecosystem. All credit of
inspiration belongs them.

### Philosophy

Go has always lacked a strong definition of equivalency, and until recently lacked the
language features necessary to make type-safe yet generic assertive statements based on
the contents of values.

This `test` (and companion `must`) package aims to provide a test-case assertion library
where the caller is in control of how types are compared, and to do so in a strongly typed
way - avoiding erroneous comparisons in the first place.

Generally there are 4 ways of asserting equivalence between types.

#### the == operator

Functions like `EqOp` and `ContainsOp` work on types that are `comparable`, i.e. are
compatible with Go's built-in `==` and `!=` operators.

#### a comparator function

Functions like `EqFunc` and `ContainsFunc` work on any type, as the caller passes in a
function that takes two arguments of that type, returning a boolean indicating equivalence.

#### an .Equal method

Functions like `Equal` and `ContainsEqual` work on types implementing the `EqualFunc`
generic interface (i.e. implement an `.Equal` method). The `.Equal` method is called
to determine equivalence.

#### the cmp.Equal or reflect.DeepEqual functions

Functions like `Eq` and `Contains` work on any type, using the `cmp.Equal` or `reflect.DeepEqual`
functions to determine equivalence. Although this is the easiest / most compatible way
to "just compare stuff", it the least deterministic way of comparing instances of a type.
Changes to the underlying types may cause unexpected changes in their equivalence (e.g.
the addition of unexported fields, function field types, etc.).

#### output

When possible, a nice `diff` output is created to show why an equivalence has failed. This
is done via the `cmp.Diff` function. For incompatible types, their `GoString` values are
printed instead.

All output is directed through `t.Log` functions, and is visible only if test verbosity is
turned on (e.g. `go test -v`).

#### fail fast vs. fail later

The `test` and `must` packages are identical, except for how test cases behave when encountering
a failure. Sometimes it is helpful for a test case to continue running even though a failure has
occurred (e.g. contains cleanup logic not captured via a `t.Cleanup` function). Other times it
make sense to fail immediately and stop the test case execution.

### PostScripts

Some tests are large and complex (like e2e testing). It can be helpful to provide more context
on test case failures beyond the actual assertion. Logging could do this, but often we want to
only produce output on failure.

The `test` and `must` packages provide a `PostScript` interface which can be implemented to
add more context in the output of failed tests. There are handy implementations of the `PostScript`
interface provided - `Sprint`, `Sprintf`, `Values`, and `Func`.

By adding one or more `PostScript` to an assertion, on failure the error message will be appended
with the additional context.

```golang
// Add a single Sprintf-string to the output of a failed test assertion.
must.Eq(t, exp, result, must.Sprintf("some more context: %v", value))
```

```golang
// Add a formatted key-value map to the output of a failed test assertion.
must.Eq(t, exp, result, must.Values(
  "one", 1,
  "two", 2,
  "fruit", "banana",
))
```

```golang
// Add the output from a closure to the output of a failed test assertion.
must.Eq(t, exp, result, must.Func(func() string {
  // ... something interesting
  return s
})
```

### Wait

Sometimes a test needs to wait on a condition for an undeterministic amount of time.
For these cases, the `wait` package provides utilities for waiting on conditionals.

The `wait` package provides constructs to create a constraint that can wait on a 
specified timeout or maximum number of attempts. There are three function types that
can be waited on.

- `BoolFunc` - of type `func() bool`, retry until function returns true
- `ErrorFunc` - of type `func() error`, retry until function returns non-nil
- `TestFunc` of type `func() (bool, error)`, retry until function returns true

The options for configuring a `wait.Constraint` include

- `Timeout(time.Duration)` - maximum amount of time to wait for condition to be satisfied
- `Attempts(int)` - maximum number of attempts to retry for condition to be satisfied
- `Gap(time.Duration)` - amount of time to wait between retry attempts

#### Fundamental form

```go
c := wait.On(
    BoolFunc(f),
    Timeout(10 * time.Seconds),
    Gap(1 * time.Second),
)
err := c.Run()
```

#### Assertions form

The `test` and `must` package implement an assertion helper for using the `wait` package.

```go
must.Wait(t, wait.On(wait.ErrorFunc(f)))
```

### Examples (equality)

```go
import "github.com/shoenig/test/must"

// ... 

e1 := Employee{ID: 100, Name: "Alice"}
e2 := Employee{ID: 101, Name: "Bob"}

// using cmp.Equal (like magic!)
must.Eq(t, e1, e2)

// using == operator
must.EqOp(t, e1, e2)

// using a custom comparator
must.EqFunc(t, e1, e2, func(a, b *Employee) bool {
    return a.ID == b.ID
})

// using .Equal method
must.Equal(t, e1, e2)
```

### Output

The `test` and `must` package attempt to create useful, readable output when an assertions goes awry. Some random examples below.

```text
test_test.go:779: expected different file permissions
↪ name: find
↪ exp: -rw-rwx-wx
↪ got: -rwxr-xr-x
```

```text
tests_test.go:569: expected maps of same values via 'eq' function
↪ difference:
  map[int]test.Person{
  	0: {ID: 100, Name: "Alice"},
  	1: {
  		ID:   101,
- 		Name: "Bob",
+ 		Name: "Bob B.",
  	},
  }
```

```text
test_test.go:520: expected slice[1].Less(slice[2])
↪ slice[1]: &{200 Bob}
↪ slice[2]: &{150 Carl}
```

```text
test_test.go:688: expected maps of same values via .Equal method
↪ differential ↷
  map[int]*test.Person{
  	0: &{ID: 100, Name: "Alice"},
  	1: &{
- 		ID:   101,
+ 		ID:   200,
  		Name: "Bob",
  	},
  }
```

```text
test_test.go:801: expected regexp match
↪ s: abcX
↪ re: abc\d
```

### License

Open source under the [MPL](LICENSE)

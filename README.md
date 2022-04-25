# test

[![GoDoc](https://godoc.org/github.com/shoenig/test?status.svg)](https://godoc.org/github.com/shoenig/test)
[![Go Report Card](https://goreportcard.com/badge/github.com/shoenig/test)](https://goreportcard.com/report/github.com/shoenig/test)
[![Run CI Tests](https://github.com/shoenig/test/actions/workflows/ci.yaml/badge.svg)](https://github.com/shoenig/test/actions/workflows/ci.yaml)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL_2.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)

A modern, generics based testing assertions library for Go.

### Overview

Package `test` provides an opinionated, lightweight library for writing
test case assertions in Go.

There are no external dependencies.

The minimum Go version is **go1.18**.

### Influence

This library was made after a ~decade of using [testify](https://github.com/stretchr/testify),
quite possibly the most used library in the whole Go ecosystem. All credit of
inspiration belongs them.

### Install

Use `go get` to grab the latest version of `test`.

```shell
go get -u github.com/shoenig/test@latest
```

### Examples (basic)

```go
// using ==
e1 := Employee{ID: 100, Name: "Alice"}
e2 := Employee{ID: 101, Name: "Bob"}
test.Eq(t, e1, e2)

// using a custom comparator
e1 := &Employee{ID: 100, Name: "Alice"}
e2 := &Employee{ID: 101, Name: "Bob"}
test.EqFunc(t, e1, e2, func(a, b *Employee) bool {
    return a.ID == b.ID
})

// using .Equals
e1 := &Employee{ID: 100, Name: "Alice"}
e2 := &Employee{ID: 101, Name: "Bob"}
test.Equals(t, e1, e2)
```

### Examples (slices)

```go
a := []*Employee{
  {ID: 100, Name: "Alice"},
  {ID: 101, Name: "Bob"},
  {ID: 102, Name: "Carl"},
}
b := []*Employee{
  {ID: 100, Name: "Alice"},
  {ID: 101, Name: "Bob"},
  {ID: 103, Name: "Dian"},
}

EqSliceFunc(tc, a, b, func(a, b *Person) bool {
  return a.ID == b.ID && a.Name == b.Name
})
```

### Examples (maps)

```go
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
```

### License

Open source under the [MPL](LICENSE)

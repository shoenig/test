// Package test provides a modern generic testing assertions library.
package test

// T is the minimal set of functions to be implemented by any testing framework
// compatible with the test package.
type T interface {
	Helper()
	Fatalf(string, ...any)
}

// EqualsFunc represents a type implementing the Equals method.
type EqualsFunc[A any] interface {
	Equals(A) bool
}

// LessFunc represents any type implementing the Less method.
type LessFunc[A any] interface {
	Less(A) bool
}

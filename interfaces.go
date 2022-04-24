// Package test provides a modern generic testing assertions library.
package test

// T is the minimal set of functions to be implemented by any testing framework
// compatible with the test package.
type T interface {
	Helper()
	Fatalf(string, ...any)
}

// EqualsFunc represents a type implementing the Equals method.
type EqualsFunc[T any] interface {
	Equals(T) bool
}

// LessFunc represents any type implementing the Less method.
type LessFunc[T any] interface {
	Less(T) bool
}

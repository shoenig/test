// Package test provides a modern generic testing assertions library.
package test

// T is the minimal set of functions to be implemented by any testing framework
// compatible with the test package.
type T interface {
	Helper()
	Fail()
	Logf(string, ...any)
}

// EqualsFunc represents a type implementing the Equals method.
type EqualsFunc[A any] interface {
	Equals(A) bool
}

// LessFunc represents any type implementing the Less method.
type LessFunc[A any] interface {
	Less(A) bool
}

// Map represents any map type where keys are comparable.
type Map[K comparable, V any] interface {
	~map[K]V
}

// MapEqualsFunc represents any map type where keys are comparable and values implement .Equals method.
type MapEqualsFunc[K comparable, V EqualsFunc[V]] interface {
	~map[K]V
}
